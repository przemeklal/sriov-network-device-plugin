// Copyright 2018 Intel Corp. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"github.com/golang/glog"
	"github.com/intel/sriov-network-device-plugin/pkg/types"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

type resourcePool struct {
	config     *types.ResourceConfig
	devices    map[string]*pluginapi.Device
	devicePool map[string]types.PciNetDevice
}

var _ types.ResourcePool = &resourcePool{}

// NewResourcePool returns an instance of resourcePool
func NewResourcePool(rc *types.ResourceConfig, deviceList []types.PciNetDevice, rFactory types.ResourceFactory) (types.ResourcePool, error) {
	filteredDevice := deviceList

	// filter by vendor list
	if rc.Selectors.Vendors != nil && len(rc.Selectors.Vendors) > 0 {
		if selector, err := rFactory.GetSelector("vendors", rc.Selectors.Vendors); err == nil {
			filteredDevice = selector.Filter(filteredDevice)
		}
	}

	// filter by device list
	if rc.Selectors.Devices != nil && len(rc.Selectors.Devices) > 0 {
		if selector, err := rFactory.GetSelector("devices", rc.Selectors.Devices); err == nil {
			filteredDevice = selector.Filter(filteredDevice)
		}
	}

	// filter by driver list
	if rc.Selectors.Drivers != nil && len(rc.Selectors.Drivers) > 0 {
		if selector, err := rFactory.GetSelector("drivers", rc.Selectors.Drivers); err == nil {
			filteredDevice = selector.Filter(filteredDevice)
		}
	}

	devicePool := make(map[string]types.PciNetDevice, 0)
	apiDevices := make(map[string]*pluginapi.Device)
	for _, dev := range filteredDevice {
		pciAddr := dev.GetPciAddr()
		devicePool[pciAddr] = dev
		apiDevices[pciAddr] = dev.GetAPIDevice()
		glog.Infof("device added: [pciAddr: %s, vendor: %s, device: %s, driver: %s]",
			dev.GetPciAddr(),
			dev.GetVendor(),
			dev.GetDeviceCode(),
			dev.GetDriver())
	}

	rPool := &resourcePool{
		config:     rc,
		devices:    apiDevices,
		devicePool: devicePool,
	}

	return rPool, nil
}

func (rp *resourcePool) GetConfig() *types.ResourceConfig {
	return rp.config
}

func (rp *resourcePool) InitDevice() error {
	// Not implemented
	return nil
}

func (rp *resourcePool) GetResourceName() string {
	return rp.config.ResourceName
}

func (rp *resourcePool) GetDevices() map[string]*pluginapi.Device {
	// returns all devices from devices[]
	return rp.devices
}

func (rp *resourcePool) Probe() bool {
	// TO-DO: Implement this
	return false
}

func (rp *resourcePool) GetDeviceSpecs(deviceIDs []string) []*pluginapi.DeviceSpec {
	glog.Infof("GetDeviceSpecs(): for devices: %v", deviceIDs)
	devSpecs := make([]*pluginapi.DeviceSpec, 0)

	// Add vfio group specific devices
	for _, id := range deviceIDs {
		if dev, ok := rp.devicePool[id]; ok {
			ds := dev.GetDeviceSpecs()
			devSpecs = append(devSpecs, ds...)
		}
	}
	return devSpecs
}

func (rp *resourcePool) GetEnvs(deviceIDs []string) []string {
	glog.Infof("GetEnvs(): for devices: %v", deviceIDs)
	devEnvs := make([]string, 0)

	// Consolidates all Envs
	for _, id := range deviceIDs {
		if dev, ok := rp.devicePool[id]; ok {
			env := dev.GetEnvVal()
			devEnvs = append(devEnvs, env)
		}
	}

	return devEnvs
}

func (rp *resourcePool) GetMounts(deviceIDs []string) []*pluginapi.Mount {
	glog.Infof("GetMounts(): for devices: %v", deviceIDs)
	devMounts := make([]*pluginapi.Mount, 0)

	for _, id := range deviceIDs {
		if dev, ok := rp.devicePool[id]; ok {
			mnt := dev.GetMounts()
			devMounts = append(devMounts, mnt...)
		}
	}
	return devMounts
}
