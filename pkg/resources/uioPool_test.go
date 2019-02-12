package resources

import (
	"reflect"

	"github.com/intel/sriov-network-device-plugin/pkg/types"

	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("UioPool", func() {
	Describe("creating new UIO resource pool", func() {
		var uioPool types.ResourcePool
		BeforeEach(func() {
			uioPool = newUioResourcePool(nil)
		})
		It("should return valid uioResourcePool object", func() {
			Expect(uioPool).NotTo(Equal(nil))
			Expect(reflect.TypeOf(uioPool)).To(Equal(reflect.TypeOf(&uioResourcePool{})))
		})
	})
	DescribeTable("getting device specs",
		func(deviceFiles map[string]string, deviceIDs []string, expected []*pluginapi.DeviceSpec) {
			pool := newUioResourcePool(nil)
			specs := pool.GetDeviceSpecs(deviceFiles, deviceIDs)
			Expect(specs).To(ConsistOf(expected))
		},
		Entry("empty",
			map[string]string{},
			[]string{},
			[]*pluginapi.DeviceSpec{},
		),
		Entry("multiple devices IDs passed, returns files for all of them and the default vfio one",
			map[string]string{
				"fake0": "/dev/uio0",
				"fake1": "/dev/uio1",
			},
			[]string{"fake0", "fake1"},
			[]*pluginapi.DeviceSpec{
				{HostPath: "/dev/uio0", ContainerPath: "/dev/uio0", Permissions: "mrw"},
				{HostPath: "/dev/uio1", ContainerPath: "/dev/uio1", Permissions: "mrw"},
			},
		),
	)
	DescribeTable("getting envs",
		func(resourceName string, deviceIDs []string, expected map[string]string) {
			conf := &types.ResourceConfig{ResourceName: resourceName}
			pool := newUioResourcePool(conf)
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
			pool := newUioResourcePool(nil)
			result := pool.GetMounts()
			Expect(result).To(BeEmpty())
		})
	})
	Describe("probing", func() {
		It("should always return false", func() {
			pool := newUioResourcePool(nil)
			result := pool.Probe(nil, nil)
			Expect(result).To(BeFalse())
		})
	})
})
