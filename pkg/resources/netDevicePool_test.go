package resources

import (
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"reflect"

	//pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	"github.com/intel/sriov-network-device-plugin/pkg/types"
	//"github.com/intel/sriov-network-device-plugin/pkg/utils"
)

var _ = Describe("NetDevicePool", func() {
	Describe("creating new newNetDevicePool", func() {
		var pool types.DeviceInfoProvider
		BeforeEach(func() {
			pool = newNetDevicePool()
		})
		It("should return valid netDevicePool object", func() {
			Expect(pool).NotTo(Equal(nil))
			Expect(reflect.TypeOf(pool)).To(Equal(reflect.TypeOf(&netDevicePool{})))
		})
	})
	/*
		DescribeTable("health probe",
			func(fs *utils.FakeFilesystem,
				sriovMode bool,
				devices map[string]*pluginapi.Device,
				expectedResult bool,
				expectedDevices map[string]*pluginapi.Device) {

				defer fs.Use()()

				conf := &types.ResourceConfig{
					ResourceName: "net",
					RootDevices:  []string{"0000:02:00.0"},
					DeviceType:   "netdevice",
					SriovMode:    sriovMode,
				}

				pool := newNetDevicePool(conf)
				result := pool.Probe(conf, devices)

				Expect(result).To(Equal(expectedResult))

				for k, v := range expectedDevices {
					Expect(devices).To(HaveKeyWithValue(k, v))
				}
			},
			PEntry("when SRIOV is enabled and links become healthy from unhealthy status",
				&utils.FakeFilesystem{
					Dirs: []string{
						"sys/bus/pci/devices/0000:02:10.0", "sys/bus/pci/devices/0000:02:00.0/net/eth0",
					},
					Files: map[string][]byte{
						"sys/bus/pci/devices/0000:02:00.0/sriov_numvfs":       []byte("1"),
						"sys/bus/pci/devices/0000:02:00.0/sriov_totalvfs":     []byte("1"),
						"sys/bus/pci/devices/0000:02:00.0/net/eth0/operstate": []byte("up"),
					},
					Symlinks: map[string]string{
						"sys/bus/pci/devices/0000:02:00.0/virtfn0": "../0000:02:10.0",
					},
				},
				true,
				map[string]*pluginapi.Device{"0000:02:10.0": &pluginapi.Device{Health: pluginapi.Unhealthy}},
				true,
				map[string]*pluginapi.Device{"0000:02:10.0": &pluginapi.Device{Health: pluginapi.Healthy}},
			),
			PEntry("when SRIOV is enabled and links become unhealthy from healthy status",
				&utils.FakeFilesystem{
					Dirs: []string{
						"sys/bus/pci/devices/0000:02:10.0", "sys/bus/pci/devices/0000:02:00.0/net/eth0",
					},
					Files: map[string][]byte{
						"sys/bus/pci/devices/0000:02:00.0/sriov_numvfs":       []byte("1"),
						"sys/bus/pci/devices/0000:02:00.0/sriov_totalvfs":     []byte("1"),
						"sys/bus/pci/devices/0000:02:00.0/net/eth0/operstate": []byte("down"),
					},
					Symlinks: map[string]string{
						"sys/bus/pci/devices/0000:02:00.0/virtfn0": "../0000:02:10.0",
					},
				},
				true,
				map[string]*pluginapi.Device{"0000:02:10.0": &pluginapi.Device{Health: pluginapi.Healthy}},
				true,
				map[string]*pluginapi.Device{"0000:02:10.0": &pluginapi.Device{Health: pluginapi.Unhealthy}},
			),
			PEntry("when SRIOV is enabled and links remain in healthy status",
				&utils.FakeFilesystem{
					Dirs: []string{
						"sys/bus/pci/devices/0000:02:10.0", "sys/bus/pci/devices/0000:02:00.0/net/eth0",
					},
					Files: map[string][]byte{
						"sys/bus/pci/devices/0000:02:00.0/sriov_numvfs":       []byte("1"),
						"sys/bus/pci/devices/0000:02:00.0/sriov_totalvfs":     []byte("1"),
						"sys/bus/pci/devices/0000:02:00.0/net/eth0/operstate": []byte("up"),
					},
					Symlinks: map[string]string{
						"sys/bus/pci/devices/0000:02:00.0/virtfn0": "../0000:02:10.0",
					},
				},
				true,
				map[string]*pluginapi.Device{"0000:02:10.0": &pluginapi.Device{Health: pluginapi.Healthy}},
				false,
				map[string]*pluginapi.Device{"0000:02:10.0": &pluginapi.Device{Health: pluginapi.Healthy}},
			),
			PEntry("when SRIOV is disabled and links become healthy",
				&utils.FakeFilesystem{
					Dirs: []string{"sys/bus/pci/devices/0000:02:00.0/net/eth0"},
					Files: map[string][]byte{
						"sys/bus/pci/devices/0000:02:00.0/sriov_numvfs":       []byte("0"),
						"sys/bus/pci/devices/0000:02:00.0/sriov_totalvfs":     []byte("0"),
						"sys/bus/pci/devices/0000:02:00.0/net/eth0/operstate": []byte("up"),
					},
				},
				false,
				map[string]*pluginapi.Device{"0000:02:00.0": &pluginapi.Device{Health: pluginapi.Unhealthy}},
				true,
				map[string]*pluginapi.Device{"0000:02:00.0": &pluginapi.Device{Health: pluginapi.Healthy}},
			),
		)*/
	Describe("getting mounts", func() {
		It("should always return an empty array", func() {
			pool := netDevicePool{}
			Expect(pool.GetMounts("fakePCIAddr")).To(BeEmpty())
		})
	})
	Describe("getting device specs", func() {
		It("should always return an empty map", func() {
			pool := netDevicePool{}
			Expect(pool.GetDeviceSpecs("fakePCIAddr")).To(BeEmpty())
		})
	})
})
