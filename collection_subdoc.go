package gocb

import (
	"context"
	"errors"
	"time"

	"github.com/couchbase/gocbcore/v8"
)

// LookupInSpec is the representation of an operation available when calling LookupIn
type LookupInSpec struct {
	op gocbcore.SubDocOp
}

// LookupInOptions are the set of options available to LookupIn.
type LookupInOptions struct {
	Context        context.Context
	Timeout        time.Duration
	WithExpiration bool
	Serializer     JSONSerializer
}

// GetSpecOptions are the options available to LookupIn subdoc Get operations.
type GetSpecOptions struct {
	IsXattr bool
}

// GetSpec indicates a path to be retrieved from the document.  The value of the path
// can later be retrieved from the LookupResult.
// The path syntax follows N1QL's path syntax (e.g. `foo.bar.baz`).
func GetSpec(path string, opts *GetSpecOptions) LookupInSpec {
	if opts == nil {
		opts = &GetSpecOptions{}
	}
	return getSpecWithFlags(path, opts.IsXattr)
}

// GetFullSpecOptions are the options available to LookupIn subdoc GetFull operations.
// There are currently no options and this is left empty for future extensibility.
type GetFullSpecOptions struct {
}

// GetFullSpec indicates that a full document should be retrieved. This command allows you
// to do things like combine with Get to fetch a document with certain Xattrs
func GetFullSpec(opts *GetFullSpecOptions) LookupInSpec {
	op := gocbcore.SubDocOp{
		Op:    gocbcore.SubDocOpGetDoc,
		Flags: gocbcore.SubdocFlag(SubdocFlagNone),
	}

	return LookupInSpec{op: op}
}

func getSpecWithFlags(path string, isXattr bool) LookupInSpec {
	var flags gocbcore.SubdocFlag
	if isXattr {
		flags |= gocbcore.SubdocFlag(SubdocFlagXattr)
	}

	op := gocbcore.SubDocOp{
		Op:    gocbcore.SubDocOpGet,
		Path:  path,
		Flags: flags,
	}

	return LookupInSpec{op: op}
}

// ExistsSpecOptions are the options available to LookupIn subdoc Exists operations.
type ExistsSpecOptions struct {
	IsXattr bool
}

// ExistsSpec is similar to Path(), but does not actually retrieve the value from the server.
// This may save bandwidth if you only need to check for the existence of a
// path (without caring for its content). You can check the status of this
// operation by using .ContentAt (and ignoring the value) or .Exists() on the LookupResult.
func ExistsSpec(path string, opts *ExistsSpecOptions) LookupInSpec {
	if opts == nil {
		opts = &ExistsSpecOptions{}
	}

	var flags gocbcore.SubdocFlag
	if opts.IsXattr {
		flags |= gocbcore.SubdocFlag(SubdocFlagXattr)
	}

	op := gocbcore.SubDocOp{
		Op:    gocbcore.SubDocOpExists,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
	}

	return LookupInSpec{op: op}
}

// CountSpecOptions are the options available to LookupIn subdoc Count operations.
type CountSpecOptions struct {
	IsXattr bool
}

// CountSpec allows you to retrieve the number of items in an array or keys within an
// dictionary within an element of a document.
func CountSpec(path string, opts *CountSpecOptions) LookupInSpec {
	if opts == nil {
		opts = &CountSpecOptions{}
	}

	var flags gocbcore.SubdocFlag
	if opts.IsXattr {
		flags |= gocbcore.SubdocFlag(SubdocFlagXattr)
	}

	op := gocbcore.SubDocOp{
		Op:    gocbcore.SubDocOpGetCount,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
	}

	return LookupInSpec{op: op}
}

// LookupIn performs a set of subdocument lookup operations on the document identified by key.
func (c *Collection) LookupIn(key string, ops []LookupInSpec, opts *LookupInOptions) (docOut *LookupInResult, errOut error) {
	if opts == nil {
		opts = &LookupInOptions{}
	}

	// Only update ctx if necessary, this means that the original ctx.Done() signal will be triggered as expected
	ctx, cancel := c.context(opts.Context, opts.Timeout)
	if cancel != nil {
		defer cancel()
	}

	res, err := c.lookupIn(ctx, key, ops, *opts)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Collection) lookupIn(ctx context.Context, key string, ops []LookupInSpec, opts LookupInOptions) (docOut *LookupInResult, errOut error) {
	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}

	var subdocs []gocbcore.SubDocOp
	for _, op := range ops {
		subdocs = append(subdocs, op.op)
	}

	// Prepend the expiration get if required, xattrs have to be at the front of the ops list.
	if opts.WithExpiration {
		op := gocbcore.SubDocOp{
			Op:    gocbcore.SubDocOpGet,
			Path:  "$document.exptime",
			Flags: gocbcore.SubdocFlag(SubdocFlagXattr),
		}

		subdocs = append([]gocbcore.SubDocOp{op}, subdocs...)
	}

	if len(ops) > 16 {
		return nil, errors.New("too many lookupIn ops specified, maximum 16")
	}

	serializer := opts.Serializer
	if serializer == nil {
		serializer = &DefaultJSONSerializer{}
	}

	ctrl := c.newOpManager(ctx)
	err = ctrl.wait(agent.LookupInEx(gocbcore.LookupInOptions{
		Key:            []byte(key),
		Ops:            subdocs,
		CollectionName: c.name(),
		ScopeName:      c.scopeName(),
	}, func(res *gocbcore.LookupInResult, err error) {
		if err != nil && !gocbcore.IsErrorStatus(err, gocbcore.StatusSubDocBadMulti) {
			errOut = maybeEnhanceKVErr(err, key, false)
			ctrl.resolve()
			return
		}

		if res != nil {
			resSet := &LookupInResult{}
			resSet.serializer = serializer
			resSet.cas = Cas(res.Cas)
			resSet.contents = make([]lookupInPartial, len(subdocs))

			for i, opRes := range res.Ops {
				// resSet.contents[i].path = opts.spec.ops[i].Path
				resSet.contents[i].err = maybeEnhanceKVErr(opRes.Err, key, false)
				if opRes.Value != nil {
					resSet.contents[i].data = append([]byte(nil), opRes.Value...)
				}
			}

			if opts.WithExpiration {
				// if expiration was requested then extract and remove it from the results
				resSet.withExpiration = true
				err = resSet.ContentAt(0, &resSet.expiration)
				if err != nil {
					errOut = err
					ctrl.resolve()
					return
				}
				resSet.contents = resSet.contents[1:]
			}

			docOut = resSet
		}

		ctrl.resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}

type subDocOp struct {
	Op         gocbcore.SubDocOpType
	Flags      gocbcore.SubdocFlag
	Path       string
	Value      interface{}
	MultiValue bool
}

// MutateInSpec is the representation of an operation available when calling MutateIn
type MutateInSpec struct {
	op subDocOp
}

// MutateInOptions are the set of options available to MutateIn.
type MutateInOptions struct {
	Timeout         time.Duration
	Context         context.Context
	Expiration      uint32
	Cas             Cas
	PersistTo       uint
	ReplicateTo     uint
	DurabilityLevel DurabilityLevel
	UpsertDocument  bool
	InsertDocument  bool
	Serializer      JSONSerializer
	// Internal: This should never be used and is not supported.
	AccessDeleted bool
}

func (c *Collection) encodeMultiArray(in interface{}, serializer JSONSerializer) ([]byte, error) {
	out, err := serializer.Serialize(in)
	if err != nil {
		return nil, err
	}

	// Assert first character is a '['
	if len(out) < 2 || out[0] != '[' {
		return nil, errors.New("not a JSON array")
	}

	out = out[1 : len(out)-1]
	return out, nil
}

// InsertSpecOptions are the options available to subdocument Insert operations.
type InsertSpecOptions struct {
	CreatePath bool
	IsXattr    bool
}

// InsertSpec inserts a value at the specified path within the document.
func InsertSpec(path string, val interface{}, opts *InsertSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &InsertSpecOptions{}
	}
	var flags SubdocFlag
	_, ok := val.(MutationMacro)
	if ok {
		flags |= SubdocFlagUseMacros
		opts.IsXattr = true
	}

	if opts.CreatePath {
		flags |= SubdocFlagCreatePath
	}
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpDictAdd,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: val,
	}

	return MutateInSpec{op: op}
}

// InsertFullSpecOptions are the options available to subdocument InsertFull operations.
type InsertFullSpecOptions struct {
}

// InsertFullSpec creates a new document if it does not exist.
// This command allows you to do things like updating xattrs whilst upserting a document.
func InsertFullSpec(path string, val interface{}, opts *InsertFullSpecOptions) MutateInSpec {
	op := subDocOp{
		Op:    gocbcore.SubDocOpAddDoc,
		Flags: gocbcore.SubdocFlag(SubdocDocFlagAddDoc),
		Value: val,
	}

	return MutateInSpec{op: op}
}

// UpsertSpecOptions are the options available to subdocument Upsert operations.
type UpsertSpecOptions struct {
	CreatePath bool
	IsXattr    bool
}

// UpsertSpec creates a new value at the specified path within the document if it does not exist, if it does exist then it
// updates it.
func UpsertSpec(path string, val interface{}, opts *UpsertSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &UpsertSpecOptions{}
	}
	var flags SubdocFlag
	_, ok := val.(MutationMacro)
	if ok {
		flags |= SubdocFlagUseMacros
		opts.IsXattr = true
	}

	if opts.CreatePath {
		flags |= SubdocFlagCreatePath
	}
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpDictSet,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: val,
	}

	return MutateInSpec{op: op}
}

// UpsertFullSpecOptions are the options available to subdocument UpsertFull operations.
// Currently exists for future extensibility and consistency.
type UpsertFullSpecOptions struct {
}

// UpsertFullSpec creates a new document if it does not exist, if it does exist then it
// updates it. This command allows you to do things like updating xattrs whilst upserting
// a document.
func UpsertFullSpec(val interface{}, opts *UpsertFullSpecOptions) MutateInSpec {
	op := subDocOp{
		Op:    gocbcore.SubDocOpSetDoc,
		Flags: gocbcore.SubdocFlag(SubdocDocFlagMkDoc),
		Value: val,
	}

	return MutateInSpec{op: op}
}

// ReplaceSpecOptions are the options available to subdocument Replace operations.
type ReplaceSpecOptions struct {
	IsXattr bool
}

// ReplaceSpec replaces the value of the field at path.
func ReplaceSpec(path string, val interface{}, opts *ReplaceSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ReplaceSpecOptions{}
	}
	var flags SubdocFlag
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpReplace,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: val,
	}

	return MutateInSpec{op: op}
}

// RemoveSpecOptions are the options available to subdocument Remove operations.
type RemoveSpecOptions struct {
	IsXattr bool
}

// RemoveSpec removes the field at path.
func RemoveSpec(path string, opts *RemoveSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &RemoveSpecOptions{}
	}
	var flags SubdocFlag
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpDelete,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
	}

	return MutateInSpec{op: op}
}

// ArrayAppendSpecOptions are the options available to subdocument ArrayAppend operations.
type ArrayAppendSpecOptions struct {
	CreatePath bool
	IsXattr    bool
	// HasMultiple adds multiple values as elements to an array.
	// When used `value` in the spec must be an array type
	// ArrayAppend("path", []int{1,2,3,4}, ArrayAppendSpecOptions{HasMultiple:true}) =>
	//   "path" [..., 1,2,3,4]
	//
	// This is a more efficient version (at both the network and server levels)
	// of doing
	// spec.ArrayAppend("path", 1, nil)
	// spec.ArrayAppend("path", 2, nil)
	// spec.ArrayAppend("path", 3, nil)
	HasMultiple bool
}

// ArrayAppendSpec adds an element(s) to the end (i.e. right) of an array
func ArrayAppendSpec(path string, val interface{}, opts *ArrayAppendSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ArrayAppendSpecOptions{}
	}
	var flags SubdocFlag
	_, ok := val.(MutationMacro)
	if ok {
		flags |= SubdocFlagUseMacros
		opts.IsXattr = true
	}
	if opts.CreatePath {
		flags |= SubdocFlagCreatePath
	}
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpArrayPushLast,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: val,
	}

	if opts.HasMultiple {
		op.MultiValue = true
	}

	return MutateInSpec{op: op}
}

// ArrayPrependSpecOptions are the options available to subdocument ArrayPrepend operations.
type ArrayPrependSpecOptions struct {
	CreatePath bool
	IsXattr    bool
	// HasMultiple adds multiple values as elements to an array.
	// When used `value` in the spec must be an array type
	// ArrayPrepend("path", []int{1,2,3,4}, ArrayPrependSpecOptions{HasMultiple:true}) =>
	//   "path" [1,2,3,4, ....]
	//
	// This is a more efficient version (at both the network and server levels)
	// of doing
	// spec.ArrayPrepend("path", 1, nil)
	// spec.ArrayPrepend("path", 2, nil)
	// spec.ArrayPrepend("path", 3, nil)
	HasMultiple bool
}

// ArrayPrependSpec adds an element to the beginning (i.e. left) of an array
func ArrayPrependSpec(path string, val interface{}, opts *ArrayPrependSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ArrayPrependSpecOptions{}
	}
	var flags SubdocFlag
	_, ok := val.(MutationMacro)
	if ok {
		flags |= SubdocFlagUseMacros
		opts.IsXattr = true
	}
	if opts.CreatePath {
		flags |= SubdocFlagCreatePath
	}
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpArrayPushFirst,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: val,
	}

	if opts.HasMultiple {
		op.MultiValue = true
	}

	return MutateInSpec{op: op}
}

// ArrayInsertSpecOptions are the options available to subdocument ArrayInsert operations.
type ArrayInsertSpecOptions struct {
	CreatePath bool
	IsXattr    bool
	// HasMultiple adds multiple values as elements to an array.
	// When used `value` in the spec must be an array type
	// ArrayInsert("path[1]", []int{1,2,3,4}, ArrayInsertSpecOptions{HasMultiple:true}) =>
	//   "path" [..., 1,2,3,4]
	//
	// This is a more efficient version (at both the network and server levels)
	// of doing
	// spec.ArrayInsert("path[2]", 1, nil)
	// spec.ArrayInsert("path[3]", 2, nil)
	// spec.ArrayInsert("path[4]", 3, nil)
	HasMultiple bool
}

// ArrayInsertSpec inserts an element at a given position within an array. The position should be
// specified as part of the path, e.g. path.to.array[3]
func ArrayInsertSpec(path string, val interface{}, opts *ArrayInsertSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ArrayInsertSpecOptions{}
	}
	var flags SubdocFlag
	_, ok := val.(MutationMacro)
	if ok {
		flags |= SubdocFlagUseMacros
		opts.IsXattr = true
	}
	if opts.CreatePath {
		flags |= SubdocFlagCreatePath
	}
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpArrayInsert,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: val,
	}

	if opts.HasMultiple {
		op.MultiValue = true
	}

	return MutateInSpec{op: op}
}

// ArrayAddUniqueSpecOptions are the options available to subdocument ArrayAddUnique operations.
type ArrayAddUniqueSpecOptions struct {
	CreatePath bool
	IsXattr    bool
}

// ArrayAddUniqueSpec adds an dictionary add unique operation to this mutation operation set.
func ArrayAddUniqueSpec(path string, val interface{}, opts *ArrayAddUniqueSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &ArrayAddUniqueSpecOptions{}
	}
	var flags SubdocFlag
	_, ok := val.(MutationMacro)
	if ok {
		flags |= SubdocFlagUseMacros
		opts.IsXattr = true
	}

	if opts.CreatePath {
		flags |= SubdocFlagCreatePath
	}
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpArrayAddUnique,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: val,
	}

	return MutateInSpec{op: op}
}

// CounterSpecOptions are the options available to subdocument Increment and Decrement operations.
type CounterSpecOptions struct {
	CreatePath bool
	IsXattr    bool
}

// IncrementSpec adds an increment operation to this mutation operation set.
func IncrementSpec(path string, delta int64, opts *CounterSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &CounterSpecOptions{}
	}
	var flags SubdocFlag
	if opts.CreatePath {
		flags |= SubdocFlagCreatePath
	}
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpCounter,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: delta,
	}

	return MutateInSpec{op: op}
}

// DecrementSpec adds a decrement operation to this mutation operation set.
func DecrementSpec(path string, delta int64, opts *CounterSpecOptions) MutateInSpec {
	if opts == nil {
		opts = &CounterSpecOptions{}
	}
	var flags SubdocFlag
	if opts.CreatePath {
		flags |= SubdocFlagCreatePath
	}
	if opts.IsXattr {
		flags |= SubdocFlagXattr
	}

	op := subDocOp{
		Op:    gocbcore.SubDocOpCounter,
		Path:  path,
		Flags: gocbcore.SubdocFlag(flags),
		Value: -delta,
	}

	return MutateInSpec{op: op}
}

// MutateIn performs a set of subdocument mutations on the document specified by key.
func (c *Collection) MutateIn(key string, ops []MutateInSpec, opts *MutateInOptions) (mutOut *MutateInResult, errOut error) {
	if opts == nil {
		opts = &MutateInOptions{}
	}

	// Only update ctx if necessary, this means that the original ctx.Done() signal will be triggered as expected
	ctx, cancel := c.context(opts.Context, opts.Timeout)
	if cancel != nil {
		defer cancel()
	}

	err := c.verifyObserveOptions(opts.PersistTo, opts.ReplicateTo, opts.DurabilityLevel)
	if err != nil {
		return nil, err
	}

	res, err := c.mutate(ctx, key, ops, *opts)
	if err != nil {
		return nil, err
	}

	if opts.PersistTo == 0 && opts.ReplicateTo == 0 {
		return res, nil
	}
	return res, c.durability(durabilitySettings{
		ctx:            opts.Context,
		key:            key,
		cas:            res.Cas(),
		mt:             res.MutationToken(),
		replicaTo:      opts.ReplicateTo,
		persistTo:      opts.PersistTo,
		forDelete:      false,
		scopeName:      c.scopeName(),
		collectionName: c.name(),
	})
}

func (c *Collection) mutate(ctx context.Context, key string, ops []MutateInSpec, opts MutateInOptions) (mutOut *MutateInResult, errOut error) {
	agent, err := c.getKvProvider()
	if err != nil {
		return nil, err
	}

	if (opts.PersistTo != 0 || opts.ReplicateTo != 0) && !c.sb.UseMutationTokens {
		return nil, configurationError{"cannot use observe based durability without mutation tokens"}
	}

	var isInsertDocument bool
	var flags SubdocDocFlag
	if opts.UpsertDocument {
		flags |= SubdocDocFlagMkDoc
	}
	if opts.InsertDocument {
		flags |= SubdocDocFlagAddDoc
	}
	if opts.AccessDeleted {
		flags |= SubdocDocFlagAccessDeleted
	}

	serializer := opts.Serializer
	if serializer == nil {
		serializer = &DefaultJSONSerializer{}
	}

	var subdocs []gocbcore.SubDocOp
	for _, op := range ops {
		if op.op.Value == nil {
			subdocs = append(subdocs, gocbcore.SubDocOp{
				Op:    op.op.Op,
				Flags: op.op.Flags,
				Path:  op.op.Path,
			})

			continue
		}

		var marshaled []byte
		var err error
		if op.op.MultiValue {
			marshaled, err = c.encodeMultiArray(op.op.Value, serializer)
		} else {
			marshaled, err = serializer.Serialize(op.op.Value)
		}
		if err != nil {
			return nil, err
		}

		subdocs = append(subdocs, gocbcore.SubDocOp{
			Op:    op.op.Op,
			Flags: op.op.Flags,
			Path:  op.op.Path,
			Value: marshaled,
		})
	}

	coerced, durabilityTimeout := c.durabilityTimeout(ctx, opts.DurabilityLevel)
	if coerced {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(durabilityTimeout)*time.Millisecond)
		defer cancel()
	}

	ctrl := c.newOpManager(ctx)
	err = ctrl.wait(agent.MutateInEx(gocbcore.MutateInOptions{
		Key:                    []byte(key),
		Flags:                  gocbcore.SubdocDocFlag(flags),
		Cas:                    gocbcore.Cas(opts.Cas),
		Ops:                    subdocs,
		Expiry:                 opts.Expiration,
		CollectionName:         c.name(),
		ScopeName:              c.scopeName(),
		DurabilityLevel:        gocbcore.DurabilityLevel(opts.DurabilityLevel),
		DurabilityLevelTimeout: durabilityTimeout,
	}, func(res *gocbcore.MutateInResult, err error) {
		if err != nil {
			errOut = maybeEnhanceKVErr(err, key, isInsertDocument)
			ctrl.resolve()
			return
		}

		mutTok := MutationToken{
			token:      res.MutationToken,
			bucketName: c.sb.BucketName,
		}
		mutRes := &MutateInResult{
			MutationResult: MutationResult{
				mt: mutTok,
				Result: Result{
					cas: Cas(res.Cas),
				},
			},
			contents: make([]mutateInPartial, len(res.Ops)),
		}

		for i, op := range res.Ops {
			mutRes.contents[i] = mutateInPartial{data: op.Value}
		}

		mutOut = mutRes

		ctrl.resolve()
	}))
	if err != nil {
		errOut = err
	}

	return
}
