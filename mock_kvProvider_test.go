// Code generated by mockery v1.0.0. DO NOT EDIT.

package gocb

import gocbcore "github.com/couchbase/gocbcore/v8"
import mock "github.com/stretchr/testify/mock"

// mockKvProvider is an autogenerated mock type for the kvProvider type
type mockKvProvider struct {
	mock.Mock
}

// AddEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) AddEx(opts gocbcore.AddOptions, cb gocbcore.StoreExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.AddOptions, gocbcore.StoreExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.AddOptions, gocbcore.StoreExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AppendEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) AppendEx(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecrementEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) DecrementEx(opts gocbcore.CounterOptions, cb gocbcore.CounterExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.CounterOptions, gocbcore.CounterExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.CounterOptions, gocbcore.CounterExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) DeleteEx(opts gocbcore.DeleteOptions, cb gocbcore.DeleteExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.DeleteOptions, gocbcore.DeleteExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.DeleteOptions, gocbcore.DeleteExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAndLockEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) GetAndLockEx(opts gocbcore.GetAndLockOptions, cb gocbcore.GetAndLockExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.GetAndLockOptions, gocbcore.GetAndLockExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.GetAndLockOptions, gocbcore.GetAndLockExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAndTouchEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) GetAndTouchEx(opts gocbcore.GetAndTouchOptions, cb gocbcore.GetAndTouchExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.GetAndTouchOptions, gocbcore.GetAndTouchExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.GetAndTouchOptions, gocbcore.GetAndTouchExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) GetEx(opts gocbcore.GetOptions, cb gocbcore.GetExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.GetOptions, gocbcore.GetExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.GetOptions, gocbcore.GetExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMetaEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) GetMetaEx(opts gocbcore.GetMetaOptions, cb gocbcore.GetMetaExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.GetMetaOptions, gocbcore.GetMetaExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.GetMetaOptions, gocbcore.GetMetaExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOneReplicaEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) GetOneReplicaEx(opts gocbcore.GetOneReplicaOptions, cb gocbcore.GetReplicaExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.GetOneReplicaOptions, gocbcore.GetReplicaExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.GetOneReplicaOptions, gocbcore.GetReplicaExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncrementEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) IncrementEx(opts gocbcore.CounterOptions, cb gocbcore.CounterExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.CounterOptions, gocbcore.CounterExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.CounterOptions, gocbcore.CounterExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LookupInEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) LookupInEx(opts gocbcore.LookupInOptions, cb gocbcore.LookupInExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.LookupInOptions, gocbcore.LookupInExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.LookupInOptions, gocbcore.LookupInExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MutateInEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) MutateInEx(opts gocbcore.MutateInOptions, cb gocbcore.MutateInExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.MutateInOptions, gocbcore.MutateInExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.MutateInOptions, gocbcore.MutateInExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NumReplicas provides a mock function with given fields:
func (_m *mockKvProvider) NumReplicas() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// ObserveEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) ObserveEx(opts gocbcore.ObserveOptions, cb gocbcore.ObserveExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.ObserveOptions, gocbcore.ObserveExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.ObserveOptions, gocbcore.ObserveExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ObserveVbEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) ObserveVbEx(opts gocbcore.ObserveVbOptions, cb gocbcore.ObserveVbExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.ObserveVbOptions, gocbcore.ObserveVbExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.ObserveVbOptions, gocbcore.ObserveVbExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PingKvEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) PingKvEx(opts gocbcore.PingKvOptions, cb gocbcore.PingKvExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.PingKvOptions, gocbcore.PingKvExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.PingKvOptions, gocbcore.PingKvExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrependEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) PrependEx(opts gocbcore.AdjoinOptions, cb gocbcore.AdjoinExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.AdjoinOptions, gocbcore.AdjoinExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReplaceEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) ReplaceEx(opts gocbcore.ReplaceOptions, cb gocbcore.StoreExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.ReplaceOptions, gocbcore.StoreExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.ReplaceOptions, gocbcore.StoreExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) SetEx(opts gocbcore.SetOptions, cb gocbcore.StoreExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.SetOptions, gocbcore.StoreExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.SetOptions, gocbcore.StoreExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TouchEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) TouchEx(opts gocbcore.TouchOptions, cb gocbcore.TouchExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.TouchOptions, gocbcore.TouchExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.TouchOptions, gocbcore.TouchExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnlockEx provides a mock function with given fields: opts, cb
func (_m *mockKvProvider) UnlockEx(opts gocbcore.UnlockOptions, cb gocbcore.UnlockExCallback) (gocbcore.PendingOp, error) {
	ret := _m.Called(opts, cb)

	var r0 gocbcore.PendingOp
	if rf, ok := ret.Get(0).(func(gocbcore.UnlockOptions, gocbcore.UnlockExCallback) gocbcore.PendingOp); ok {
		r0 = rf(opts, cb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gocbcore.PendingOp)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(gocbcore.UnlockOptions, gocbcore.UnlockExCallback) error); ok {
		r1 = rf(opts, cb)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
