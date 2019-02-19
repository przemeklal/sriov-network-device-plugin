// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import types "github.com/intel/sriov-network-device-plugin/pkg/types"

// ResourceFactory is an autogenerated mock type for the ResourceFactory type
type ResourceFactory struct {
	mock.Mock
}

// GetResourcePool provides a mock function with given fields: _a0
func (_m *ResourceFactory) GetResourcePool(_a0 *types.ResourceConfig) types.ResourcePool {
	ret := _m.Called(_a0)

	var r0 types.ResourcePool
	if rf, ok := ret.Get(0).(func(*types.ResourceConfig) types.ResourcePool); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.ResourcePool)
		}
	}

	return r0
}

// GetResourceServer provides a mock function with given fields: _a0
func (_m *ResourceFactory) GetResourceServer(_a0 types.ResourcePool) (types.ResourceServer, error) {
	ret := _m.Called(_a0)

	var r0 types.ResourceServer
	if rf, ok := ret.Get(0).(func(types.ResourcePool) types.ResourceServer); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.ResourceServer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.ResourcePool) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
