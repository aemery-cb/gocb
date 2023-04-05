package gocb

import (
	"context"
	"errors"
	"time"

	"github.com/couchbase/gocbcore/v10"
	"github.com/couchbase/gocbcore/v10/memd"
)

type kvProviderCore struct {
	agent kvProviderCoreProvider
}

var _ kvProvider = &kvProviderCore{}

func (p *kvProviderCore) MutateIn(opm *kvOpManager, action StoreSemantics, ops []MutateInSpec, flags SubdocDocFlag) (*MutateInResult, error) {
	synced := newCoreKvOpManager(opm, p)

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
			// TODO: this might be a bug.
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
	synced := newCoreKvOpManager(opm, p)

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
				docOut.contents[i].data = opRes.Value
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

func (p *kvProviderCore) Scan(c *Collection, scanType ScanType, opts *ScanOptions) (*ScanResult, error) {

	config, err := p.waitForConfigSnapshot(opts.Context, time.Now().Add(opts.Timeout))
	if err != nil {
		return nil, err
	}

	numVbuckets, err := config.NumVbuckets()
	if err != nil {
		return nil, err
	}

	if numVbuckets == 0 {
		return nil, makeInvalidArgumentsError("can only use RangeScan with couchbase buckets")
	}

	opm, err := p.newRangeScanOpManager(c, scanType, numVbuckets, p.agent, opts.ParentSpan, opts.ConsistentWith,
		opts.IDsOnly, opts.Sort)
	if err != nil {
		return nil, err
	}

	opm.SetTranscoder(opts.Transcoder)
	opm.SetContext(opts.Context)
	opm.SetImpersonate(opts.Internal.User)
	opm.SetTimeout(opts.Timeout)
	opm.SetRetryStrategy(opts.RetryStrategy)
	opm.SetItemLimit(opts.BatchItemLimit)
	opm.SetByteLimit(opts.BatchByteLimit)

	if err := opm.CheckReadyForOp(); err != nil {
		return nil, err
	}

	return opm.Scan()
}

func (p *kvProviderCore) Add(opm *kvOpManager) (*MutationResult, error) {
	synced := newCoreKvOpManager(opm, p)

	var errOut error
	var mutOut *MutationResult
	err := synced.Wait(p.agent.Add(gocbcore.AddOptions{
		Key:                    synced.DocumentID(),
		Value:                  synced.ValueBytes(),
		Flags:                  synced.ValueFlags(),
		Expiry:                 durationToExpiry(opm.Expiry()),
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
	synced := newCoreKvOpManager(opm, p)

	var errOut error
	var mutOut *MutationResult
	err := synced.Wait(p.agent.Set(gocbcore.SetOptions{
		Key:                    synced.DocumentID(),
		Value:                  synced.ValueBytes(),
		Flags:                  synced.ValueFlags(),
		Expiry:                 durationToExpiry(synced.Expiry()),
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
	synced := newCoreKvOpManager(opm, p)

	var errOut error
	var mutOut *MutationResult

	err := synced.Wait(p.agent.Replace(gocbcore.ReplaceOptions{
		Key:                    synced.DocumentID(),
		Value:                  synced.ValueBytes(),
		Flags:                  synced.ValueFlags(),
		Expiry:                 durationToExpiry(opm.Expiry()),
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
	synced := newCoreKvOpManager(opm, p)

	var docOut *GetResult
	var errOut error
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

		docOut = doc

		synced.Resolve(nil)
	}))
	if err != nil {
		errOut = err
	}

	return docOut, errOut
}

func (p *kvProviderCore) GetAndTouch(opm *kvOpManager) (*GetResult, error) {
	synced := newCoreKvOpManager(opm, p)

	var getOut *GetResult
	var errOut error

	err := synced.Wait(p.agent.GetAndTouch(gocbcore.GetAndTouchOptions{
		Key:            synced.DocumentID(),
		Expiry:         durationToExpiry(opm.Expiry()),
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
	synced := newCoreKvOpManager(opm, p)

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
	synced := newCoreKvOpManager(opm, p)

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
	synced := newCoreKvOpManager(opm, p)

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
	synced := newCoreKvOpManager(opm, p)

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
	synced := newCoreKvOpManager(opm, p)

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
	synced := newCoreKvOpManager(opm, p)

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

func (p *kvProviderCore) GetAllReplicas(c *Collection, id string, opts *GetAllReplicaOptions) (*GetAllReplicasResult, error) {
	if opts == nil {
		opts = &GetAllReplicaOptions{}
	}

	var tracectx RequestSpanContext
	if opts.ParentSpan != nil {
		tracectx = opts.ParentSpan.Context()
	}

	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	span := c.startKvOpTrace("get_all_replicas", tracectx, false)

	// Timeout needs to be adjusted here, since we use it at the bottom of this
	// function, but the remaining options are all passed downwards and get handled
	// by those functions rather than us.
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = c.timeoutsConfig.KVTimeout
	}

	deadline := time.Now().Add(timeout)
	transcoder := opts.Transcoder
	retryStrategy := opts.RetryStrategy

	snapshot, err := p.waitForConfigSnapshot(ctx, deadline)
	if err != nil {
		return nil, err
	}

	numReplicas, err := snapshot.NumReplicas()
	if err != nil {
		return nil, err
	}

	numServers := numReplicas + 1
	outCh := make(chan *GetReplicaResult, numServers)
	cancelCh := make(chan struct{})

	var recorder ValueRecorder
	if !opts.noMetrics {
		recorder, err = c.meter.ValueRecorder(meterValueServiceKV, "get_all_replicas")
		if err != nil {
			logDebugf("Failed to create value recorder: %v", err)
		}
	}

	repRes := &GetAllReplicasResult{
		totalRequests:       uint32(numServers),
		resCh:               outCh,
		cancelCh:            cancelCh,
		span:                span,
		childReqsCompleteCh: make(chan struct{}),
		valueRecorder:       recorder,
		startedTime:         time.Now(),
	}

	// Loop all the servers and populate the result object
	for replicaIdx := 0; replicaIdx < numServers; replicaIdx++ {
		go func(replicaIdx int) {
			// This timeout value will cause the getOneReplica operation to timeout after our deadline has expired,
			// as the deadline has already begun. getOneReplica timing out before our deadline would cause inconsistent
			// behaviour.
			res, err := p.getOneReplica(context.Background(), span, id, replicaIdx, transcoder, retryStrategy, cancelCh,
				timeout, opts.Internal.User, c)
			if err != nil {
				repRes.addFailed()
				logDebugf("Failed to fetch replica from replica %d: %s", replicaIdx, err)
			} else {
				repRes.addResult(res)
			}
		}(replicaIdx)
	}

	// Start a timer to close it after the deadline
	go func() {
		select {
		case <-time.After(time.Until(deadline)):
			// If we timeout, we should close the result
			err := repRes.Close()
			if err != nil {
				logDebugf("failed to close GetAllReplicas response: %s", err)
			}
		case <-cancelCh:
		// If the cancel channel closes, we are done
		case <-ctx.Done():
			err := repRes.Close()
			if err != nil {
				logDebugf("failed to close GetAllReplicas response: %s", err)
			}
		}
	}()

	return repRes, nil
}

func (p *kvProviderCore) getOneReplica(
	ctx context.Context,
	span RequestSpan,
	id string,
	replicaIdx int,
	transcoder Transcoder,
	retryStrategy RetryStrategy,
	cancelCh chan struct{},
	timeout time.Duration,
	user string,
	c *Collection,
) (*GetReplicaResult, error) {
	opm := newKvOpManager(c, "get_replica", span)
	defer opm.Finish(false)

	synced := newCoreKvOpManager(opm, p)

	synced.SetDocumentID(id)
	synced.SetTranscoder(transcoder)
	synced.SetRetryStrategy(retryStrategy)
	synced.SetTimeout(timeout)
	synced.SetCancelCh(cancelCh)
	synced.SetImpersonate(user)
	synced.SetContext(ctx)

	if replicaIdx == 0 {
		var docOut *GetReplicaResult
		var errOut error
		err := synced.Wait(p.agent.Get(gocbcore.GetOptions{
			Key:            synced.DocumentID(),
			CollectionName: synced.CollectionName(),
			ScopeName:      synced.ScopeName(),
			RetryStrategy:  synced.RetryStrategy(),
			TraceContext:   synced.TraceSpanContext(),
			Deadline:       synced.Deadline(),
			User:           synced.Impersonate(),
		}, func(res *gocbcore.GetResult, err error) {
			if err != nil {
				errOut = opm.EnhanceErr(err)
				synced.Reject()
				return
			}

			docOut = &GetReplicaResult{}
			docOut.cas = Cas(res.Cas)
			docOut.transcoder = synced.Transcoder()
			docOut.contents = res.Value
			docOut.flags = res.Flags
			docOut.isReplica = false

			synced.Resolve(nil)
		}))
		if err != nil {
			errOut = err
		}
		return docOut, errOut
	}

	var docOut *GetReplicaResult
	var errOut error
	err := synced.Wait(p.agent.GetOneReplica(gocbcore.GetOneReplicaOptions{
		Key:            synced.DocumentID(),
		ReplicaIdx:     replicaIdx,
		CollectionName: synced.CollectionName(),
		ScopeName:      synced.ScopeName(),
		RetryStrategy:  synced.RetryStrategy(),
		TraceContext:   synced.TraceSpanContext(),
		Deadline:       synced.Deadline(),
		User:           synced.Impersonate(),
	}, func(res *gocbcore.GetReplicaResult, err error) {
		if err != nil {
			errOut = opm.EnhanceErr(err)
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
	synced := newCoreKvOpManager(opm, p)

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
	synced := newCoreKvOpManager(opm, p)

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
	synced := newCoreKvOpManager(opm, p)

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
	synced := newCoreKvOpManager(opm, p)

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

func (p *kvProviderCore) BulkGet(opts gocbcore.GetOptions, cb gocbcore.GetCallback) (gocbcore.PendingOp, error) {
	return p.agent.Get(opts, cb)
}
func (p *kvProviderCore) BulkGetAndTouch(opts gocbcore.GetAndTouchOptions, cb gocbcore.GetAndTouchCallback) (gocbcore.PendingOp, error) {
	return p.agent.GetAndTouch(opts, cb)
}
func (p *kvProviderCore) BulkTouch(opts gocbcore.TouchOptions, cb gocbcore.TouchCallback) (gocbcore.PendingOp, error) {
	return p.agent.Touch(opts, cb)
}
func (p *kvProviderCore) BulkDelete(opts gocbcore.DeleteOptions, cb gocbcore.DeleteCallback) (gocbcore.PendingOp, error) {
	return p.agent.Delete(opts, cb)
}
func (p *kvProviderCore) BulkSet(opts gocbcore.SetOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error) {
	return p.agent.Set(opts, cb)
}
func (p *kvProviderCore) BulkAdd(opts gocbcore.AddOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error) {
	return p.agent.Add(opts, cb)
}
func (p *kvProviderCore) BulkReplace(opts gocbcore.ReplaceOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error) {
	return p.agent.Replace(opts, cb)
}
func (p *kvProviderCore) BulkAppend(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinCallback) (gocbcore.PendingOp, error) {
	return p.agent.Append(opts, cb)
}
func (p *kvProviderCore) BulkPrepend(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinCallback) (gocbcore.PendingOp, error) {
	return p.agent.Prepend(opts, cb)
}
func (p *kvProviderCore) BulkIncrement(opts gocbcore.CounterOptions, cb gocbcore.CounterCallback) (gocbcore.PendingOp, error) {
	return p.agent.Increment(opts, cb)
}
func (p *kvProviderCore) BulkDecrement(opts gocbcore.CounterOptions, cb gocbcore.CounterCallback) (gocbcore.PendingOp, error) {
	return p.agent.Decrement(opts, cb)
}
