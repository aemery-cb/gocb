package gocb

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/couchbase/gocbcore/v10"
	"github.com/couchbase/gocbcore/v10/memd"
)

type kvProviderCore struct {
	agent *gocbcore.Agent
}

var _ kvProvider = &kvProviderCore{}

func (p *kvProviderCore) MutateIn(opm *kvOpManager, action StoreSemantics, ops []MutateInSpec, flags SubdocDocFlag) (*MutateInResult, error) {
	synced := newSyncKvOpManager(opm)

	expiry := synced.Expiry()
	preserveTTL := synced.PreserveExpiry()

	docFlags := memd.SubdocDocFlag(flags)
	switch action {
	case StoreSemanticsReplace:
		// this is default behavior
		if expiry > 0 && preserveTTL {
			return nil, makeInvalidArgumentsError("cannot use preserve ttl with expiry for replace store semantics")
		}
	case StoreSemanticsUpsert:
		docFlags |= memd.SubdocDocFlagMkDoc
	case StoreSemanticsInsert:
		if preserveTTL {
			return nil, makeInvalidArgumentsError("cannot use preserve ttl with insert store semantics")
		}

		docFlags |= memd.SubdocDocFlagAddDoc
	default:
		return nil, makeInvalidArgumentsError("invalid StoreSemantics value provided")
	}

	var subdocs []gocbcore.SubDocOp
	for _, op := range ops {
		if op.path == "" {
			switch op.op {
			case memd.SubDocOpDictAdd:
				return nil, makeInvalidArgumentsError("cannot specify a blank path with InsertSpec")
			case memd.SubDocOpDictSet:
				return nil, makeInvalidArgumentsError("cannot specify a blank path with UpsertSpec")
			case memd.SubDocOpDelete:
				op.op = memd.SubDocOpDeleteDoc
			case memd.SubDocOpReplace:
				op.op = memd.SubDocOpSetDoc
			default:
			}
		}

		etrace := synced.parent.startKvOpTrace("request_encoding", opm.TraceSpanContext(), true)
		bytes, flags, err := jsonMarshalMutateSpec(op)
		etrace.End()
		if err != nil {
			return nil, err
		}

		if op.createPath {
			flags |= memd.SubdocFlagMkDirP
		}

		if op.isXattr {
			flags |= memd.SubdocFlagXattrPath
		}

		subdocs = append(subdocs, gocbcore.SubDocOp{
			Op:    op.op,
			Flags: flags,
			Path:  op.path,
			Value: bytes,
		})
	}

	var errOut error
	var mutOut *MutateInResult
	err := synced.Wait(p.agent.MutateIn(gocbcore.MutateInOptions{
		Key:                    synced.DocumentID(),
		Flags:                  docFlags,
		Cas:                    gocbcore.Cas(synced.Cas()),
		Ops:                    subdocs,
		Expiry:                 durationToExpiry(expiry),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
		PreserveExpiry:         preserveTTL,
	}, func(res *gocbcore.MutateInResult, err error) {
		if err != nil {
			// GOCBC-1019: Due to a previous bug in gocbcore we need to convert cas mismatch back to exists.
			if kvErr, ok := err.(*gocbcore.KeyValueError); ok {
				if errors.Is(kvErr.InnerError, ErrCasMismatch) {
					kvErr.InnerError = ErrDocumentExists
				}
			}
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mutOut = &MutateInResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)
		mutOut.contents = make([]mutateInPartial, len(res.Ops))
		for i, op := range res.Ops {
			mutOut.contents[i] = mutateInPartial{data: op.Value}
		}

		synced.Resolve(mutOut.mt)
	}))

	if err != nil {
		errOut = err
	}
	return mutOut, errOut
}

func (p *kvProviderCore) LookupIn(opm *kvOpManager, ops []LookupInSpec, flags SubdocDocFlag) (*LookupInResult, error) {

	var subdocs []gocbcore.SubDocOp
	for _, op := range ops {
		if op.op == memd.SubDocOpGet && op.path == "" {
			if op.isXattr {
				return nil, errors.New("invalid xattr fetch with no path")
			}

			subdocs = append(subdocs, gocbcore.SubDocOp{
				Op:    memd.SubDocOpGetDoc,
				Flags: memd.SubdocFlag(SubdocFlagNone),
			})
			continue
		} else if op.op == memd.SubDocOpDictSet && op.path == "" {
			//TODO: this might be a bug.
			if op.isXattr {
				return nil, errors.New("invalid xattr set with no path")
			}

			subdocs = append(subdocs, gocbcore.SubDocOp{
				Op:    memd.SubDocOpSetDoc,
				Flags: memd.SubdocFlag(SubdocFlagNone),
			})
			continue
		}

		flags := memd.SubdocFlagNone
		if op.isXattr {
			flags |= memd.SubdocFlagXattrPath
		}

		subdocs = append(subdocs, gocbcore.SubDocOp{
			Op:    op.op,
			Path:  op.path,
			Flags: flags,
		})
	}
	synced := newSyncKvOpManager(opm)

	var errOut error
	var docOut *LookupInResult
	err := synced.Wait(p.agent.LookupIn(gocbcore.LookupInOptions{
		Key:            synced.DocumentID(),
		Ops:            subdocs,
		CollectionName: synced.CollectionName(),
		ScopeName:      synced.ScopeName(),
		RetryStrategy:  synced.RetryStrategy(),
		TraceContext:   synced.TraceSpanContext(),
		Deadline:       synced.Deadline(),
		Flags:          memd.SubdocDocFlag(flags),
		User:           synced.Impersonate(),
	}, func(res *gocbcore.LookupInResult, err error) {
		if err != nil && res == nil {
			errOut = synced.EnhanceErr(err)
		}

		if res != nil {
			docOut = &LookupInResult{}
			docOut.cas = Cas(res.Cas)
			docOut.contents = make([]lookupInPartial, len(subdocs))
			for i, opRes := range res.Ops {
				docOut.contents[i].err = synced.EnhanceErr(opRes.Err)
				docOut.contents[i].data = json.RawMessage(opRes.Value)
			}
		}

		if err == nil {
			synced.Resolve(nil)
		} else {
			synced.Reject()
		}
	}))
	if err != nil {
		errOut = err
	}

	return docOut, errOut
}

func (p *kvProviderCore) Scan(ScanType, *kvOpManager) (*ScanResult, error) {

	return nil, nil
}

func (p *kvProviderCore) Add(opm *kvOpManager) (*MutationResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var mutOut *MutationResult
	err := synced.Wait(p.agent.Add(gocbcore.AddOptions{
		Key:   synced.DocumentID(),
		Value: synced.ValueBytes(),
		Flags: synced.ValueFlags(),
		// Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
	}, func(res *gocbcore.StoreResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = synced.EnhanceMt(res.MutationToken)

		synced.Resolve(mutOut.mt)
	}))

	if err != nil {
		errOut = err
	}

	return mutOut, errOut
}

func (p *kvProviderCore) Set(opm *kvOpManager) (*MutationResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var mutOut *MutationResult

	err := synced.Wait(p.agent.Set(gocbcore.SetOptions{
		Key:   synced.DocumentID(),
		Value: synced.ValueBytes(),
		Flags: synced.ValueFlags(),
		//Expiry:                 durationToExpiry(opts.Expiry),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
	}, func(sr *gocbcore.StoreResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(sr.Cas)
		mutOut.mt = synced.EnhanceMt(sr.MutationToken)

		synced.Resolve(mutOut.mt)

	}))

	if err != nil {
		errOut = err
	}
	return mutOut, errOut
}

func (p *kvProviderCore) Replace(opm *kvOpManager) (*MutationResult, error) {

	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var mutOut *MutationResult

	err := synced.Wait(p.agent.Replace(gocbcore.ReplaceOptions{
		Key:   synced.DocumentID(),
		Value: synced.ValueBytes(),
		Flags: synced.ValueFlags(),
		//Expiry:                 durationToExpiry(opts.Expiry),
		Cas:                    gocbcore.Cas(synced.Cas()),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
		PreserveExpiry:         synced.PreserveExpiry(),
	}, func(sr *gocbcore.StoreResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(sr.Cas)
		mutOut.mt = synced.EnhanceMt(sr.MutationToken)

		synced.Resolve(mutOut.mt)

	}))

	if err != nil {
		errOut = err
	}

	return mutOut, errOut
}

func (p *kvProviderCore) Get(opm *kvOpManager) (*GetResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var getOut *GetResult

	err := synced.Wait(p.agent.Get(gocbcore.GetOptions{
		Key:            opm.DocumentID(),
		CollectionName: opm.CollectionName(),
		ScopeName:      opm.ScopeName(),
		RetryStrategy:  opm.RetryStrategy(),
		TraceContext:   opm.TraceSpanContext(),
		Deadline:       opm.Deadline(),
		User:           opm.Impersonate(),
	}, func(res *gocbcore.GetResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
			synced.Reject()
			return
		}

		doc := &GetResult{
			Result: Result{
				cas: Cas(res.Cas),
			},
			transcoder: opm.Transcoder(),
			contents:   res.Value,
			flags:      res.Flags,
		}

		getOut = doc

		synced.Resolve(nil)
	}))

	if err != nil {
		errOut = err
	}

	return getOut, errOut

}

func (p *kvProviderCore) GetAndTouch(opm *kvOpManager) (*GetResult, error) {
	synced := newSyncKvOpManager(opm)

	var getOut *GetResult
	var errOut error

	err := synced.Wait(p.agent.GetAndTouch(gocbcore.GetAndTouchOptions{
		Key: synced.DocumentID(),
		//Expiry:         durationToExpiry(expiry),
		CollectionName: synced.CollectionName(),
		ScopeName:      synced.ScopeName(),
		RetryStrategy:  synced.RetryStrategy(),
		TraceContext:   synced.TraceSpanContext(),
		Deadline:       synced.Deadline(),
		User:           synced.Impersonate(),
	}, func(res *gocbcore.GetAndTouchResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		if res != nil {
			getOut = &GetResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				transcoder: opm.Transcoder(),
				contents:   res.Value,
				flags:      res.Flags,
			}
		}

		synced.Resolve(nil)
	}))

	if err != nil {
		errOut = err
	}

	return getOut, errOut

}

func (p *kvProviderCore) GetAndLock(opm *kvOpManager) (*GetResult, error) {
	synced := newSyncKvOpManager(opm)

	var errOut error
	var getResult *GetResult

	err := synced.Wait(p.agent.GetAndLock(gocbcore.GetAndLockOptions{
		Key:            synced.DocumentID(),
		LockTime:       uint32(synced.LockTime() / time.Second),
		CollectionName: synced.CollectionName(),
		ScopeName:      synced.ScopeName(),
		RetryStrategy:  synced.RetryStrategy(),
		TraceContext:   synced.TraceSpanContext(),
		Deadline:       synced.Deadline(),
		User:           synced.Impersonate(),
	}, func(res *gocbcore.GetAndLockResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		if res != nil {
			doc := &GetResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				transcoder: opm.Transcoder(),
				contents:   res.Value,
				flags:      res.Flags,
			}

			getResult = doc
		}

		synced.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}

	return getResult, errOut
}

func (p *kvProviderCore) Exists(opm *kvOpManager) (*ExistsResult, error) {
	synced := newSyncKvOpManager(opm)
	defer synced.Finish(false)

	var docExists *ExistsResult
	var errOut error
	err := synced.Wait(p.agent.GetMeta(gocbcore.GetMetaOptions{
		Key:            synced.DocumentID(),
		CollectionName: synced.CollectionName(),
		ScopeName:      synced.ScopeName(),
		RetryStrategy:  synced.RetryStrategy(),
		TraceContext:   synced.TraceSpanContext(),
		Deadline:       synced.Deadline(),
		User:           synced.Impersonate(),
	}, func(res *gocbcore.GetMetaResult, err error) {
		if errors.Is(err, ErrDocumentNotFound) {
			docExists = &ExistsResult{
				Result: Result{
					cas: Cas(0),
				},
				docExists: false,
			}
			synced.Resolve(nil)
			return
		}

		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		if res != nil {
			docExists = &ExistsResult{
				Result: Result{
					cas: Cas(res.Cas),
				},
				docExists: res.Deleted == 0,
			}
		}

		synced.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}

	return docExists, errOut
}

func (p *kvProviderCore) Delete(opm *kvOpManager) (*MutationResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var mutOut *MutationResult
	err := synced.Wait(p.agent.Delete(gocbcore.DeleteOptions{
		Key:                    synced.DocumentID(),
		Cas:                    gocbcore.Cas(opm.Cas()),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
	}, func(res *gocbcore.DeleteResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = synced.EnhanceMt(res.MutationToken)

		synced.Resolve(mutOut.mt)
	}))
	if err != nil {
		errOut = err
	}

	return mutOut, errOut
}

func (p *kvProviderCore) Unlock(opm *kvOpManager) error {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)
	var errOut error
	err := synced.Wait(p.agent.Unlock(gocbcore.UnlockOptions{
		Key:            synced.DocumentID(),
		Cas:            gocbcore.Cas(opm.Cas()),
		CollectionName: synced.CollectionName(),
		ScopeName:      synced.ScopeName(),
		RetryStrategy:  synced.RetryStrategy(),
		TraceContext:   synced.TraceSpanContext(),
		Deadline:       synced.Deadline(),
		User:           synced.Impersonate(),
	}, func(res *gocbcore.UnlockResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mt := synced.EnhanceMt(res.MutationToken)
		synced.Resolve(mt)
	}))

	if err != nil {
		errOut = err
	}
	return errOut
}

func (p *kvProviderCore) Touch(opm *kvOpManager) (*MutationResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var mutOut *MutationResult
	err := synced.Wait(p.agent.Touch(gocbcore.TouchOptions{
		Key:            synced.DocumentID(),
		Expiry:         durationToExpiry(synced.Expiry()),
		CollectionName: synced.CollectionName(),
		ScopeName:      synced.ScopeName(),
		RetryStrategy:  synced.RetryStrategy(),
		TraceContext:   synced.TraceSpanContext(),
		Deadline:       synced.Deadline(),
		User:           synced.Impersonate(),
	}, func(res *gocbcore.TouchResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = synced.EnhanceMt(res.MutationToken)

		synced.Resolve(mutOut.mt)
	}))

	if err != nil {
		errOut = err
	}

	return mutOut, errOut

}

func (p *kvProviderCore) GetReplica(opm *kvOpManager) (*GetReplicaResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(true)
	var errOut error
	var docOut *GetReplicaResult

	err := synced.Wait(p.agent.GetOneReplica(gocbcore.GetOneReplicaOptions{
		Key:            synced.DocumentID(),
		ReplicaIdx:     synced.ReplicaIndex(),
		CollectionName: synced.CollectionName(),
		ScopeName:      synced.ScopeName(),
		RetryStrategy:  synced.RetryStrategy(),
		TraceContext:   synced.TraceSpanContext(),
		Deadline:       synced.Deadline(),
		User:           synced.Impersonate(),
	}, func(res *gocbcore.GetReplicaResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		docOut = &GetReplicaResult{}
		docOut.cas = Cas(res.Cas)
		docOut.transcoder = synced.Transcoder()
		docOut.contents = res.Value
		docOut.flags = res.Flags
		docOut.isReplica = true

		synced.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}

	return docOut, errOut

}

func (p *kvProviderCore) Prepend(opm *kvOpManager) (*MutationResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var mutOut *MutationResult

	err := synced.Wait(p.agent.Prepend(gocbcore.AdjoinOptions{
		Key:                    synced.DocumentID(),
		Value:                  synced.AdjoinBytes(),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		Cas:                    gocbcore.Cas(synced.Cas()),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
	}, func(res *gocbcore.AdjoinResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		synced.Resolve(mutOut.mt)
	}))

	if err != nil {
		errOut = err
	}

	return mutOut, errOut
}

func (p *kvProviderCore) Append(opm *kvOpManager) (*MutationResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var mutOut *MutationResult

	err := synced.Wait(p.agent.Append(gocbcore.AdjoinOptions{
		Key:                    synced.DocumentID(),
		Value:                  synced.AdjoinBytes(),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		Cas:                    gocbcore.Cas(synced.Cas()),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
	}, func(res *gocbcore.AdjoinResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		mutOut = &MutationResult{}
		mutOut.cas = Cas(res.Cas)
		mutOut.mt = opm.EnhanceMt(res.MutationToken)

		synced.Resolve(mutOut.mt)
	}))

	if err != nil {
		errOut = err
	}

	return mutOut, errOut
}

func (p *kvProviderCore) Increment(opm *kvOpManager) (*CounterResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)

	var errOut error
	var countOut *CounterResult

	err := synced.Wait(p.agent.Increment(gocbcore.CounterOptions{
		Key:                    synced.DocumentID(),
		Delta:                  synced.Delta(),
		Initial:                synced.Initial(),
		Expiry:                 durationToExpiry(synced.Expiry()),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
	}, func(res *gocbcore.CounterResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		countOut = &CounterResult{}
		countOut.cas = Cas(res.Cas)
		countOut.mt = synced.EnhanceMt(res.MutationToken)
		countOut.content = res.Value

		synced.Resolve(countOut.mt)
	}))

	if err != nil {
		errOut = err
	}

	return countOut, errOut

}

func (p *kvProviderCore) Decrement(opm *kvOpManager) (*CounterResult, error) {
	synced := newSyncKvOpManager(opm)

	defer synced.Finish(false)
	var errOut error
	var countOut *CounterResult

	err := synced.Wait(p.agent.Decrement(gocbcore.CounterOptions{
		Key:                    synced.DocumentID(),
		Delta:                  synced.Delta(),
		Initial:                synced.Initial(),
		Expiry:                 durationToExpiry(synced.Expiry()),
		CollectionName:         synced.CollectionName(),
		ScopeName:              synced.ScopeName(),
		DurabilityLevel:        synced.DurabilityLevel(),
		DurabilityLevelTimeout: synced.DurabilityTimeout(),
		RetryStrategy:          synced.RetryStrategy(),
		TraceContext:           synced.TraceSpanContext(),
		Deadline:               synced.Deadline(),
		User:                   synced.Impersonate(),
	}, func(res *gocbcore.CounterResult, err error) {
		if err != nil {
			errOut = synced.EnhanceErr(err)
			synced.Reject()
			return
		}

		countOut = &CounterResult{}
		countOut.cas = Cas(res.Cas)
		countOut.mt = synced.EnhanceMt(res.MutationToken)
		countOut.content = res.Value

		synced.Resolve(countOut.mt)
	}))
	if err != nil {
		errOut = err
	}
	return countOut, errOut

}
