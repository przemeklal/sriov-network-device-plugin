package resources

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/intel/sriov-network-device-plugin/pkg/types"
)

var _ = Describe("GenericPool", func() {
	DescribeTable("getting envs",
		func(resourceName string, deviceIDs []string, expected map[string]string) {
			conf := &types.ResourceConfig{ResourceName: resourceName}
			pool := newGenericResourcePool(conf)
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
	Describe("getting device file", func() {
		It("should always return empty string", func() {
			pool := genericResourcePool{}
			Expect(pool.GetDeviceFile("fake")).To(Equal(""))
		})
	})
	Describe("getting mounts", func() {
		It("should always return empty array", func() {
			pool := genericResourcePool{}
			Expect(pool.GetMounts()).To(BeEmpty())
		})
	})
	Describe("getting device specs", func() {
		It("should always return empty array", func() {
			pool := genericResourcePool{}
			result := pool.GetDeviceSpecs(nil, nil)
			Expect(result).To(BeEmpty())
		})
	})
	Describe("getting device specs", func() {
		It("should always return false", func() {
			pool := genericResourcePool{}
			result := pool.Probe(nil, nil)
			Expect(result).To(BeFalse())
		})
	})
})
