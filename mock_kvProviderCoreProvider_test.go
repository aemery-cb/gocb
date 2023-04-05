// Code generated by mockery v2.26.1. DO NOT EDIT.

package gocb

import (
	gocbcore "github.com/couchbase/gocbcore/v10"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// mockKvProviderCoreProvider is an autogenerated mock type for the kvProviderCoreProvider type
type mockKvProviderCoreProvider struct {
	mock.Mock
}

// Add provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Add(opts gocbcore.AddOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.AddOptions, gocbcore.StoreCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.AddOptions, gocbcore.StoreCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.AddOptions, gocbcore.StoreCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Append provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Append(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Decrement provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Decrement(opts gocbcore.CounterOptions, cb gocbcore.CounterCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.CounterOptions, gocbcore.CounterCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.CounterOptions, gocbcore.CounterCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.CounterOptions, gocbcore.CounterCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Delete(opts gocbcore.DeleteOptions, cb gocbcore.DeleteCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.DeleteOptions, gocbcore.DeleteCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.DeleteOptions, gocbcore.DeleteCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.DeleteOptions, gocbcore.DeleteCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Get(opts gocbcore.GetOptions, cb gocbcore.GetCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.GetOptions, gocbcore.GetCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.GetOptions, gocbcore.GetCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.GetOptions, gocbcore.GetCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAndLock provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) GetAndLock(opts gocbcore.GetAndLockOptions, cb gocbcore.GetAndLockCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.GetAndLockOptions, gocbcore.GetAndLockCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.GetAndLockOptions, gocbcore.GetAndLockCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.GetAndLockOptions, gocbcore.GetAndLockCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAndTouch provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) GetAndTouch(opts gocbcore.GetAndTouchOptions, cb gocbcore.GetAndTouchCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.GetAndTouchOptions, gocbcore.GetAndTouchCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.GetAndTouchOptions, gocbcore.GetAndTouchCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.GetAndTouchOptions, gocbcore.GetAndTouchCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMeta provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) GetMeta(opts gocbcore.GetMetaOptions, cb gocbcore.GetMetaCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.GetMetaOptions, gocbcore.GetMetaCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.GetMetaOptions, gocbcore.GetMetaCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.GetMetaOptions, gocbcore.GetMetaCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOneReplica provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) GetOneReplica(opts gocbcore.GetOneReplicaOptions, cb gocbcore.GetReplicaCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.GetOneReplicaOptions, gocbcore.GetReplicaCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.GetOneReplicaOptions, gocbcore.GetReplicaCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.GetOneReplicaOptions, gocbcore.GetReplicaCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Increment provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Increment(opts gocbcore.CounterOptions, cb gocbcore.CounterCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.CounterOptions, gocbcore.CounterCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.CounterOptions, gocbcore.CounterCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.CounterOptions, gocbcore.CounterCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LookupIn provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) LookupIn(opts gocbcore.LookupInOptions, cb gocbcore.LookupInCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.LookupInOptions, gocbcore.LookupInCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.LookupInOptions, gocbcore.LookupInCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.LookupInOptions, gocbcore.LookupInCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MutateIn provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) MutateIn(opts gocbcore.MutateInOptions, cb gocbcore.MutateInCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.MutateInOptions, gocbcore.MutateInCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.MutateInOptions, gocbcore.MutateInCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.MutateInOptions, gocbcore.MutateInCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Observe provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Observe(opts gocbcore.ObserveOptions, cb gocbcore.ObserveCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.ObserveOptions, gocbcore.ObserveCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.ObserveOptions, gocbcore.ObserveCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.ObserveOptions, gocbcore.ObserveCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ObserveVb provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) ObserveVb(opts gocbcore.ObserveVbOptions, cb gocbcore.ObserveVbCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.ObserveVbOptions, gocbcore.ObserveVbCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.ObserveVbOptions, gocbcore.ObserveVbCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.ObserveVbOptions, gocbcore.ObserveVbCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Prepend provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Prepend(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RangeScanCancel provides a mock function with given fields: scanUUID, vbID, opts, cb
func (_m *mockKvProviderCoreProvider) RangeScanCancel(scanUUID []byte, vbID uint16, opts gocbcore.RangeScanCancelOptions, cb gocbcore.RangeScanCancelCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(scanUUID, vbID, opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, uint16, gocbcore.RangeScanCancelOptions, gocbcore.RangeScanCancelCallback) (gocbcore.PendingOp, error)); ok {
		return rf(scanUUID, vbID, opts, cb)
	}
	if rf, ok := ret.Get(0).(func([]byte, uint16, gocbcore.RangeScanCancelOptions, gocbcore.RangeScanCancelCallback) gocbcore.PendingOp); ok {
		r0 = rf(scanUUID, vbID, opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, uint16, gocbcore.RangeScanCancelOptions, gocbcore.RangeScanCancelCallback) error); ok {
		r1 = rf(scanUUID, vbID, opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RangeScanContinue provides a mock function with given fields: scanUUID, vbID, opts, dataCb, actionCb
func (_m *mockKvProviderCoreProvider) RangeScanContinue(scanUUID []byte, vbID uint16, opts gocbcore.RangeScanContinueOptions, dataCb gocbcore.RangeScanContinueDataCallback, actionCb gocbcore.RangeScanContinueActionCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(scanUUID, vbID, opts, dataCb, actionCb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, uint16, gocbcore.RangeScanContinueOptions, gocbcore.RangeScanContinueDataCallback, gocbcore.RangeScanContinueActionCallback) (gocbcore.PendingOp, error)); ok {
		return rf(scanUUID, vbID, opts, dataCb, actionCb)
	}
	if rf, ok := ret.Get(0).(func([]byte, uint16, gocbcore.RangeScanContinueOptions, gocbcore.RangeScanContinueDataCallback, gocbcore.RangeScanContinueActionCallback) gocbcore.PendingOp); ok {
		r0 = rf(scanUUID, vbID, opts, dataCb, actionCb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, uint16, gocbcore.RangeScanContinueOptions, gocbcore.RangeScanContinueDataCallback, gocbcore.RangeScanContinueActionCallback) error); ok {
		r1 = rf(scanUUID, vbID, opts, dataCb, actionCb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RangeScanCreate provides a mock function with given fields: vbID, opts, cb
func (_m *mockKvProviderCoreProvider) RangeScanCreate(vbID uint16, opts gocbcore.RangeScanCreateOptions, cb gocbcore.RangeScanCreateCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(vbID, opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(uint16, gocbcore.RangeScanCreateOptions, gocbcore.RangeScanCreateCallback) (gocbcore.PendingOp, error)); ok {
		return rf(vbID, opts, cb)
	}
	if rf, ok := ret.Get(0).(func(uint16, gocbcore.RangeScanCreateOptions, gocbcore.RangeScanCreateCallback) gocbcore.PendingOp); ok {
		r0 = rf(vbID, opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(uint16, gocbcore.RangeScanCreateOptions, gocbcore.RangeScanCreateCallback) error); ok {
		r1 = rf(vbID, opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Replace provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Replace(opts gocbcore.ReplaceOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.ReplaceOptions, gocbcore.StoreCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.ReplaceOptions, gocbcore.StoreCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.ReplaceOptions, gocbcore.StoreCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Set(opts gocbcore.SetOptions, cb gocbcore.StoreCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.SetOptions, gocbcore.StoreCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.SetOptions, gocbcore.StoreCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.SetOptions, gocbcore.StoreCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Touch provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Touch(opts gocbcore.TouchOptions, cb gocbcore.TouchCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.TouchOptions, gocbcore.TouchCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.TouchOptions, gocbcore.TouchCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.TouchOptions, gocbcore.TouchCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Unlock provides a mock function with given fields: opts, cb
func (_m *mockKvProviderCoreProvider) Unlock(opts gocbcore.UnlockOptions, cb gocbcore.UnlockCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(gocbcore.UnlockOptions, gocbcore.UnlockCallback) (gocbcore.PendingOp, error)); ok {
		return rf(opts, cb)
	}
	if rf, ok := ret.Get(0).(func(gocbcore.UnlockOptions, gocbcore.UnlockCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(gocbcore.UnlockOptions, gocbcore.UnlockCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WaitForConfigSnapshot provides a mock function with given fields: deadline, opts, cb
func (_m *mockKvProviderCoreProvider) WaitForConfigSnapshot(deadline time.Time, opts gocbcore.WaitForConfigSnapshotOptions, cb gocbcore.WaitForConfigSnapshotCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(deadline, opts, cb)

	var r0 gocbcore.PendingOp
	var r1 error
	if rf, ok := ret.Get(0).(func(time.Time, gocbcore.WaitForConfigSnapshotOptions, gocbcore.WaitForConfigSnapshotCallback) (gocbcore.PendingOp, error)); ok {
		return rf(deadline, opts, cb)
	}
	if rf, ok := ret.Get(0).(func(time.Time, gocbcore.WaitForConfigSnapshotOptions, gocbcore.WaitForConfigSnapshotCallback) gocbcore.PendingOp); ok {
		r0 = rf(deadline, opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	if rf, ok := ret.Get(1).(func(time.Time, gocbcore.WaitForConfigSnapshotOptions, gocbcore.WaitForConfigSnapshotCallback) error); ok {
		r1 = rf(deadline, opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockKvProviderCoreProvider interface {
	mock.TestingT
	Cleanup(func())
}

// newMockKvProviderCoreProvider creates a new instance of mockKvProviderCoreProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockKvProviderCoreProvider(t mockConstructorTestingTnewMockKvProviderCoreProvider) *mockKvProviderCoreProvider {
	mock := &mockKvProviderCoreProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
