package resources

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"fmt"
	"os"

	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"

	"github.com/intel/sriov-network-device-plugin/pkg/types"
	"github.com/intel/sriov-network-device-plugin/pkg/types/mocks"
	"github.com/intel/sriov-network-device-plugin/pkg/utils"
)

var _ = Describe("Server", func() {
	var (
		rs *resourceServer
		fs = &utils.FakeFilesystem{Dirs: []string{"tmp"}}
	)
	Describe("creating new instance of resource server", func() {
		Context("valid arguments are passed", func() {
			BeforeEach(func() {
				defer fs.Use()()
				sockDir = fs.RootDir
				rp := mocks.ResourcePool{}
				rp.On("GetResourceName").Return("fakename")
				obj := newResourceServer("fakeprefix", "fakesuffix", &rp)
				rs = obj.(*resourceServer)
			})
			It("should have the properties correctly assigned", func() {
				Expect(rs.resourcePool.GetResourceName()).To(Equal("fakename"))
				Expect(rs.resourceNamePrefix).To(Equal("fakeprefix"))
				Expect(rs.endPoint).To(Equal("fakename.fakesuffix"))
			})
		})
	})
	DescribeTable("registering with Kubelet",
		func(shouldRunServer, shouldServerFail, shouldFail bool) {
			defer fs.Use()()
			sockDir = fs.RootDir
			rp := mocks.ResourcePool{}
			rp.On("GetResourceName").Return("fakename")
			obj := newResourceServer("fakeprefix", "fakesuffix", &rp)
			rs = obj.(*resourceServer)
			registrationServer := createFakeRegistrationServer(shouldServerFail)
			if shouldRunServer {
				os.MkdirAll(pluginapi.DevicePluginPath, 0755)
				registrationServer.start()
			}
			err := rs.register()
			if shouldFail {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
			if shouldRunServer {
				registrationServer.stop()
			}
		},
		Entry("when can't connect to Kubelet should fail", false, true, true),
		Entry("when device plugin unable to register with Kubelet should fail", true, true, true),
		Entry("succesfully shouldn't fail", true, false, false),
	)
	Describe("initializating server", func() {
		Context("when device discovery has failed", func() {
			var err error
			BeforeEach(func() {
				defer fs.Use()()
				sockDir = fs.RootDir
				rp := mocks.ResourcePool{}
				rp.
					On("GetResourceName").Return("fakename").
					On("DiscoverDevices").Return(fmt.Errorf("fake error"))

				rs := newResourceServer("fakeprefix", "fakesuffix", &rp)
				err = rs.Init()
			})
			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when device discovery has been succesful", func() {
			var err error
			BeforeEach(func() {
				defer fs.Use()()
				sockDir = fs.RootDir
				rp := mocks.ResourcePool{}
				rp.
					On("GetResourceName").Return("fakename").
					On("DiscoverDevices").Return(nil)

				rs := newResourceServer("fakeprefix", "fakesuffix", &rp)
				err = rs.Init()
			})
			It("should not fail", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
	Describe("starting, restarting, watching and stopping the resource server", func() {
		// integration-like test for the resource server (positive case)
		Context("succesfully", func() {
			It("should register with kubelet", func(done Done) {
				defer fs.Use()()
				sockDir = fs.RootDir
				fakeConf := &types.ResourceConfig{ResourceName: "fake", RootDevices: []string{"fakeid"}}
				rp := mocks.ResourcePool{}
				rp.
					On("GetConfig").Return(fakeConf).
					On("GetResourceName").Return("fake.com").
					On("DiscoverDevices").Return(nil).
					On("GetDevices").Return(map[string]*pluginapi.Device{}).
					On("Probe", fakeConf, map[string]*pluginapi.Device{}).Return(true)

				rs := newResourceServer("fake.com", "fake", &rp).(*resourceServer)

				registrationServer := createFakeRegistrationServer(false)
				os.MkdirAll(pluginapi.DevicePluginPath, 0755)
				registrationServer.start()
				defer registrationServer.stop()

				err := rs.Start()
				Expect(err).NotTo(HaveOccurred())

				err = rs.restart()
				Expect(err).NotTo(HaveOccurred())

				go func() {
					rs.Watch()
				}()

				go func() {
					err = rs.Stop()
					Expect(err).NotTo(HaveOccurred())
				}()

				Eventually(rs.termSignal).Should(Receive())
				Eventually(rs.stopWatcher).Should(Receive())

				close(done)
			})
		})
	})

	DescribeTable("allocating",
		func(req *pluginapi.AllocateRequest, expectedRespLength int, shouldFail bool) {
			rp := mocks.ResourcePool{}
			rp.On("GetResourceName").
				Return("fake.com").
				On("GetDeviceFiles").
				Return(map[string]string{"00:00.01": "/dev/fake"}).
				On("GetDeviceSpecs", map[string]string{"00:00.01": "/dev/fake"}, []string{"00:00.01"}).
				Return([]*pluginapi.DeviceSpec{{ContainerPath: "/dev/fake", HostPath: "/dev/fake", Permissions: "rw"}}).
				On("GetEnvs", []string{"00:00.01"}).
				Return(map[string]string{"PCIDEVICE_FAKE_COM_ADDR": "00:00.01"}).
				On("GetMounts").
				Return([]*pluginapi.Mount{{ContainerPath: "/dev/fake", HostPath: "/dev/fake", ReadOnly: false}})

			rs := newResourceServer("fake.com", "fake", &rp)

			resp, err := rs.Allocate(nil, req)

			Expect(len(resp.GetContainerResponses())).To(Equal(expectedRespLength))

			if shouldFail {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("allocating succesfully 1 deviceID",
			&pluginapi.AllocateRequest{
				ContainerRequests: []*pluginapi.ContainerAllocateRequest{{DevicesIDs: []string{"00:00.01"}}},
			},
			1,
			false,
		),
		PEntry("allocating deviceID that does not exist",
			&pluginapi.AllocateRequest{
				ContainerRequests: []*pluginapi.ContainerAllocateRequest{{DevicesIDs: []string{"00:00.02"}}},
			},
			0,
			true,
		),
		Entry("empty AllocateRequest", &pluginapi.AllocateRequest{}, 0, false),
	)
	Describe("running PreStartContainer", func() {
		It("should not fail", func() {
			rs := &resourceServer{}
			resp, err := rs.PreStartContainer(nil, nil)
			Expect(resp).NotTo(Equal(nil))
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Describe("running GetDevicePluginOptions", func() {
		It("should not fail", func() {
			rs := &resourceServer{}
			resp, err := rs.GetDevicePluginOptions(nil, nil)
			Expect(resp).NotTo(Equal(nil))
			Expect(err).NotTo(HaveOccurred())
		})
	})
	DescribeTable("getting env variables",
		func(in, expected map[string]string) {
			defer fs.Use()()
			sockDir = fs.RootDir
			deviceIDs := []string{"fakeid"}
			rp := mocks.ResourcePool{}
			rp.
				On("GetResourceName").Return("fake").
				On("GetEnvs", deviceIDs).Return(in)

			obj := newResourceServer("fake.com", "fake", &rp)
			rs := *obj.(*resourceServer)
			actual := rs.getEnvs(deviceIDs)
			if len(in) == 0 {
				Expect(actual).To(BeEmpty())
			}
			for k, v := range expected {
				Expect(actual).To(HaveKeyWithValue(k, v))
			}
		},
		Entry("empty map", map[string]string{}, map[string]string{}),
		Entry("some values",
			map[string]string{
				"key1": "value1",
				"key2": "",
			},
			map[string]string{
				"PCIDEVICE_FAKE_COM_KEY1": "value1",
				"PCIDEVICE_FAKE_COM_KEY2": "",
			},
		),
	)

	Describe("ListAndWatch", func() {
		Context("when first Send call in DevicePlugin_ListAndWatch failed", func() {
			It("should fail", func() {
				defer fs.Use()()
				sockDir = fs.RootDir
				rp := mocks.ResourcePool{}
				rp.On("GetResourceName").Return("fake.com").
					On("GetDevices").Return(map[string]*pluginapi.Device{"00:00.01": {ID: "00:00.01", Health: "Healthy"}}).Once()

				rs := newResourceServer("fake.com", "fake", &rp).(*resourceServer)

				lwSrv := &fakeListAndWatchServer{
					resourceServer: rs,
					sendCallToFail: 1,
				}

				err := rs.ListAndWatch(&pluginapi.Empty{}, lwSrv)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when Send call in DevicePlugin_ListAndWatch breaks", func() {
			It("should receive not fail", func(done Done) {
				defer fs.Use()()
				sockDir = fs.RootDir
				rp := mocks.ResourcePool{}
				rp.On("GetResourceName").Return("fake.com").
					On("GetDevices").Return(map[string]*pluginapi.Device{"00:00.01": {ID: "00:00.01", Health: "Healthy"}}).Once().
					On("GetDevices").Return(map[string]*pluginapi.Device{"00:00.02": {ID: "00:00.02", Health: "Healthy"}}).Once()

				rs := newResourceServer("fake.com", "fake", &rp).(*resourceServer)

				lwSrv := &fakeListAndWatchServer{
					resourceServer: rs,
					sendCallToFail: 2,
					updates:        make(chan bool),
				}

				// run ListAndWatch which will send initial update
				var err error
				go func() {
					err = rs.ListAndWatch(&pluginapi.Empty{}, lwSrv)
					// because DevicePlugin_ListAndWatch breaks...
					Expect(err).To(HaveOccurred())
				}()

				// wait for the initial update to reach ListAndWatchServer
				Eventually(lwSrv.updates).Should(Receive())
				// this time it should break
				rs.updateSignal <- true
				Eventually(lwSrv.updates).ShouldNot(Receive())

				close(done)
			})
		})
		Context("when received multiple update requests and then the term signal", func() {
			It("should receive not fail", func(done Done) {
				defer fs.Use()()
				sockDir = fs.RootDir
				rp := mocks.ResourcePool{}
				rp.On("GetResourceName").Return("fake.com").
					On("GetDevices").Return(map[string]*pluginapi.Device{"00:00.01": {ID: "00:00.01", Health: "Healthy"}}).Once().
					On("GetDevices").Return(map[string]*pluginapi.Device{"00:00.02": {ID: "00:00.02", Health: "Healthy"}}).Once()

				rs := newResourceServer("fake.com", "fake", &rp).(*resourceServer)

				lwSrv := &fakeListAndWatchServer{
					resourceServer: rs,
					sendCallToFail: 0, // no failures on purpose
					updates:        make(chan bool),
				}

				// run ListAndWatch which will send initial update
				var err error
				go func() {
					err = rs.ListAndWatch(&pluginapi.Empty{}, lwSrv)
					Expect(err).NotTo(HaveOccurred())
				}()

				// wait for the initial update to reach ListAndWatchServer
				Eventually(lwSrv.updates).Should(Receive())

				// send another set of updates and wait for the ListAndWatchServer
				rs.updateSignal <- true
				Eventually(lwSrv.updates).Should(Receive())

				// finally send term signal
				rs.termSignal <- true

				close(done)
			})
		})
	})
})
