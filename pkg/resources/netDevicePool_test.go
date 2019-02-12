package resources

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"reflect"

	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	"github.com/intel/sriov-network-device-plugin/pkg/types"
	"github.com/intel/sriov-network-device-plugin/pkg/utils"
)

var _ = Describe("NetDevicePool", func() {
	Describe("creating new newNetDevicePool", func() {
		var pool types.ResourcePool
		BeforeEach(func() {
			pool = newNetDevicePool(nil)
		})
		It("should return valid netDevicePool object", func() {
			Expect(pool).NotTo(Equal(nil))
			Expect(reflect.TypeOf(pool)).To(Equal(reflect.TypeOf(&netDevicePool{})))
		})
	})
	DescribeTable("getting envs",
		func(resourceName string, deviceIDs []string, expected map[string]string) {
			conf := &types.ResourceConfig{ResourceName: resourceName}
			pool := newNetDevicePool(conf)
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
		Entry("when SRIOV is enabled and links become healthy from unhealthy status",
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
		Entry("when SRIOV is enabled and links become unhealthy from healthy status",
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
		Entry("when SRIOV is enabled and links remain in healthy status",
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
		Entry("when SRIOV is disabled and links become healthy",
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
	)
	Describe("getting device file", func() {
		It("should always return empty string", func() {
			pool := netDevicePool{}
			Expect(pool.GetDeviceFile("fake")).To(Equal(""))
		})
	})
	Describe("getting mounts", func() {
		It("should always return empty array", func() {
			pool := netDevicePool{}
			Expect(pool.GetMounts()).To(BeEmpty())
		})
	})
	Describe("getting device specs", func() {
		It("should always return empty map", func() {
			pool := netDevicePool{}
			Expect(pool.GetDeviceSpecs(map[string]string{}, []string{})).To(BeEmpty())
		})
	})
})
