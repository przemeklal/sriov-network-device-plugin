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
		var vfioPool types.DeviceInfoProvider
		BeforeEach(func() {
			vfioPool = newVfioResourcePool()
		})
		It("should return valid vfioResourcePool object", func() {
			Expect(vfioPool).NotTo(Equal(nil))
			Expect(reflect.TypeOf(vfioPool)).To(Equal(reflect.TypeOf(&vfioResourcePool{})))
		})
	})
	/*
		FDescribe("device discovery", func() {
			Context("in SRIOV mode when SRIOV is enabled", func() {
				var (
					pool *vfioResourcePool
					err  error
				)
				BeforeEach(func() {
					fs := &utils.FakeFilesystem{
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
				var (
					pool *vfioResourcePool
					err  error
					fs   *utils.FakeFilesystem
				)
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
				var (
					pool *vfioResourcePool
					err  error
					fs   *utils.FakeFilesystem
				)
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
	*/
	DescribeTable("GetDeviceSpecs",
		func(fs *utils.FakeFilesystem, pciAddr string, expected []*pluginapi.DeviceSpec) {
			defer fs.Use()()

			pool := newVfioResourcePool()
			specs := pool.GetDeviceSpecs(pciAddr)
			Expect(specs).To(ConsistOf(expected))
		},
		Entry("empty and returning default common vfio device file only",
			&utils.FakeFilesystem{},
			"",
			[]*pluginapi.DeviceSpec{
				{HostPath: "/dev/vfio/vfio", ContainerPath: "/dev/vfio/vfio", Permissions: "mrw"},
			},
		),
		Entry("PCI address passed, returns DeviceSpec with paths to its VFIO devices and additional default VFIO path",
			&utils.FakeFilesystem{
				Dirs: []string{
					"sys/bus/pci/devices/0000:02:00.0", "sys/kernel/iommu_groups/0",
				},
				Symlinks: map[string]string{
					"sys/bus/pci/devices/0000:02:00.0/iommu_group": "../../../../kernel/iommu_groups/0",
				},
			},
			"0000:02:00.0",
			[]*pluginapi.DeviceSpec{
				{HostPath: "/dev/vfio/0", ContainerPath: "/dev/vfio/0", Permissions: "mrw"},
				{HostPath: "/dev/vfio/vfio", ContainerPath: "/dev/vfio/vfio", Permissions: "mrw"},
			},
		),
	)
	Describe("getting mounts", func() {
		It("should always return empty array of mounts", func() {
			pool := newVfioResourcePool()
			result := pool.GetMounts("fakeAddr")
			Expect(result).To(BeEmpty())
		})
	})
	Describe("getting env val", func() {
		It("should always return passed PCI address", func() {
			in := "00:02.0"
			pool := newVfioResourcePool()
			out := pool.GetEnvVal(in)
			Expect(out).To(Equal(in))
		})
	})
})
