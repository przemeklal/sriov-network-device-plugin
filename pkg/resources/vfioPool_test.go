package resources

import (
	"reflect"

	"github.com/intel/sriov-network-device-plugin/pkg/types"
	"github.com/intel/sriov-network-device-plugin/pkg/utils"

	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("VfioPool", func() {
	Describe("creating new VFIO resource pool", func() {
		var vfioPool types.ResourcePool
		BeforeEach(func() {
			vfioPool = newVfioResourcePool(nil)
		})
		It("should return valid vfioResourcePool object", func() {
			Expect(vfioPool).NotTo(Equal(nil))
			Expect(reflect.TypeOf(vfioPool)).To(Equal(reflect.TypeOf(&vfioResourcePool{})))
		})
	})
	Describe("device discovery", func() {
		var (
			pool *vfioResourcePool
			err  error
			fs   *utils.FakeFilesystem
		)
		Context("in SRIOV mode when SRIOV is enabled", func() {
			BeforeEach(func() {
				fs = &utils.FakeFilesystem{
					Dirs: []string{
						"sys/bus/pci/devices/0000:02:00.0",
						"sys/bus/pci/devices/0000:02:00.1",
						"sys/bus/pci/devices/0000:02:10.0",
						"sys/bus/pci/devices/0000:02:11.0",
						"sys/kernel/iommu_groups/0",
						"sys/kernel/iommu_groups/1",
					},
					Files: map[string][]byte{
						"sys/bus/pci/devices/0000:02:00.0/sriov_numvfs":   []byte("1"),
						"sys/bus/pci/devices/0000:02:00.1/sriov_numvfs":   []byte("1"),
						"sys/bus/pci/devices/0000:02:00.0/sriov_totalvfs": []byte("1"),
						"sys/bus/pci/devices/0000:02:00.1/sriov_totalvfs": []byte("1"),
					},
					Symlinks: map[string]string{
						"sys/bus/pci/devices/0000:02:00.0/virtfn0":     "../0000:02:10.0",
						"sys/bus/pci/devices/0000:02:00.1/virtfn0":     "../0000:02:11.0",
						"sys/bus/pci/devices/0000:02:10.0/iommu_group": "../../../../kernel/iommu_groups/0",
						"sys/bus/pci/devices/0000:02:11.0/iommu_group": "../../../../kernel/iommu_groups/1",
					},
				}
				defer fs.Use()()

				conf := &types.ResourceConfig{
					ResourceName: "net",
					RootDevices:  []string{"0000:02:00.0", "0000:02:00.1"},
					DeviceType:   "vfio",
					SriovMode:    true,
				}

				pool = newVfioResourcePool(conf).(*vfioResourcePool)
				err = pool.DiscoverDevices()
			})
			It("should not fail", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should populate devices list", func() {
				expected := map[string]*pluginapi.Device{
					"0000:02:10.0": &pluginapi.Device{"0000:02:10.0", "Healthy"},
					"0000:02:11.0": &pluginapi.Device{"0000:02:11.0", "Healthy"},
				}
				for k, v := range expected {
					Expect(pool.devices).To(HaveKeyWithValue(k, v))
				}
			})
			It("should populate device files list", func() {
				expected := map[string]string{
					"0000:02:10.0": "/dev/vfio/0",
					"0000:02:11.0": "/dev/vfio/1",
				}
				for k, v := range expected {
					Expect(pool.deviceFiles).To(HaveKeyWithValue(k, v))
				}
			})
		})
		Context("in SRIOV mode when SRIOV is disabled", func() {
			BeforeEach(func() {
				fs = &utils.FakeFilesystem{
					Dirs: []string{
						"sys/bus/pci/devices/0000:02:00.0",
						"sys/bus/pci/devices/0000:02:00.1",
					},
					Files: map[string][]byte{
						"sys/bus/pci/devices/0000:02:00.0/sriov_numvfs":   []byte("0"),
						"sys/bus/pci/devices/0000:02:00.1/sriov_numvfs":   []byte("0"),
						"sys/bus/pci/devices/0000:02:00.0/sriov_totalvfs": []byte("1"),
						"sys/bus/pci/devices/0000:02:00.1/sriov_totalvfs": []byte("1"),
					},
				}
				defer fs.Use()()

				conf := &types.ResourceConfig{
					ResourceName: "net",
					RootDevices:  []string{"0000:02:00.0", "0000:02:00.1"},
					DeviceType:   "vfio",
					SriovMode:    true,
				}

				pool = newVfioResourcePool(conf).(*vfioResourcePool)
				err = pool.DiscoverDevices()
			})
			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})
			It("should not populate devices list", func() {
				Expect(pool.devices).To(BeEmpty())
			})
		})
		Context("PFs only", func() {
			BeforeEach(func() {

				fs = &utils.FakeFilesystem{
					Dirs: []string{
						"sys/bus/pci/devices/0000:02:00.0",
						"sys/bus/pci/devices/0000:02:00.1",
						"sys/kernel/iommu_groups/0",
						"sys/kernel/iommu_groups/1",
					},
					Symlinks: map[string]string{
						"sys/bus/pci/devices/0000:02:00.0/iommu_group": "../../../../kernel/iommu_groups/0",
						"sys/bus/pci/devices/0000:02:00.1/iommu_group": "../../../../kernel/iommu_groups/1",
					},
				}
				defer fs.Use()()

				conf := &types.ResourceConfig{
					ResourceName: "net",
					RootDevices:  []string{"0000:02:00.0", "0000:02:00.1"},
					DeviceType:   "vfio",
					SriovMode:    false,
				}
				pool = newVfioResourcePool(conf).(*vfioResourcePool)
				err = pool.DiscoverDevices()
			})
			It("should not fail", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should populate devices list", func() {
				expected := map[string]string{
					"0000:02:00.0": "/dev/vfio/0",
					"0000:02:00.1": "/dev/vfio/1",
				}
				for k, v := range expected {
					Expect(pool.deviceFiles).To(HaveKeyWithValue(k, v))
				}
			})
		})
	})
	DescribeTable("GetDeviceSpecs",
		func(deviceFiles map[string]string, deviceIDs []string, expected []*pluginapi.DeviceSpec) {
			pool := newVfioResourcePool(nil)
			specs := pool.GetDeviceSpecs(deviceFiles, deviceIDs)
			Expect(specs).To(ConsistOf(expected))
		},
		Entry("empty and returning default common vfio device file only",
			map[string]string{},
			[]string{},
			[]*pluginapi.DeviceSpec{
				{HostPath: "/dev/vfio/vfio", ContainerPath: "/dev/vfio/vfio", Permissions: "mrw"},
			},
		),
		Entry("multiple devices IDs passed, returns files for all of them and the default vfio one",
			map[string]string{
				"fake0": "/dev/vfio/0",
				"fake1": "/dev/vfio/1",
			},
			[]string{"fake0", "fake1"},
			[]*pluginapi.DeviceSpec{
				{HostPath: "/dev/vfio/0", ContainerPath: "/dev/vfio/0", Permissions: "mrw"},
				{HostPath: "/dev/vfio/1", ContainerPath: "/dev/vfio/1", Permissions: "mrw"},
				{HostPath: "/dev/vfio/vfio", ContainerPath: "/dev/vfio/vfio", Permissions: "mrw"},
			},
		),
	)
	DescribeTable("getting envs",
		func(resourceName string, deviceIDs []string, expected map[string]string) {
			conf := &types.ResourceConfig{ResourceName: resourceName}
			pool := newVfioResourcePool(conf)
			envs := pool.GetEnvs(deviceIDs)
			for k, v := range expected {
				Expect(envs).To(HaveKeyWithValue(k, v))
			}
		},
		Entry("for empty device IDs should return empty map of envs",
			"fake",
			[]string{},
			map[string]string{},
		),
		Entry("for empty device IDs should return empty map of envs",
			"fake", []string{"02:00.0", "02:00.1"},
			map[string]string{"fake": "02:00.0,02:00.1"},
		),
	)
	Describe("getting mounts", func() {
		It("should always return empty array of mounts", func() {
			pool := newVfioResourcePool(nil)
			result := pool.GetMounts()
			Expect(result).To(BeEmpty())
		})
	})
	Describe("probing", func() {
		It("should always return false", func() {
			pool := newVfioResourcePool(nil)
			result := pool.Probe(nil, nil)
			Expect(result).To(BeFalse())
		})
	})
})
