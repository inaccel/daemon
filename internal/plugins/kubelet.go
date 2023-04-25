package plugins

import (
	"context"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/docker/docker/daemon/graphdriver"
	"github.com/inaccel/daemon/internal/driver"
	"github.com/inaccel/daemon/pkg/plugin"
	"github.com/moby/sys/mount"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pluginregistrationv1 "k8s.io/kubelet/pkg/apis/pluginregistration/v1"
)

type Kubelet struct {
	path string

	driver driver.Driver
	plugin.Plugin
}

func NewKubelet(ctx context.Context, driver driver.Driver) plugin.Plugin {
	ctx, cancel := context.WithCancel(ctx)

	kubelet := &Kubelet{
		path: "/var/lib/kubelet/plugins_registry/inaccel.sock",
	}

	kubelet.driver = driver

	kubelet.Plugin = plugin.Base(func() {
		if listener, err := listen(kubelet.path); err == nil {
			go func() {
				<-ctx.Done()

				listener.Close()
			}()

			server := grpc.NewServer()
			csi.RegisterIdentityServer(server, kubelet)
			csi.RegisterNodeServer(server, kubelet)
			pluginregistrationv1.RegisterRegistrationServer(server, kubelet)

			server.Serve(listener)
		} else {
			logrus.Error(err)
		}
	}, cancel)

	return kubelet
}

func (plugin Kubelet) GetInfo(ctx context.Context, request *pluginregistrationv1.InfoRequest) (*pluginregistrationv1.PluginInfo, error) {
	logrus.Info("Kubelet/GetInfo")

	response := &pluginregistrationv1.PluginInfo{
		Type: pluginregistrationv1.CSIPlugin,
		Name: "inaccel",
		SupportedVersions: []string{
			"v1.3.0",
		},
	}

	return response, nil
}

func (plugin Kubelet) GetPluginCapabilities(ctx context.Context, request *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	logrus.Info("Kubelet/GetPluginCapabilities")

	response := &csi.GetPluginCapabilitiesResponse{}

	return response, nil
}

func (plugin Kubelet) GetPluginInfo(ctx context.Context, request *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	logrus.Info("Kubelet/GetPluginInfo")

	response := &csi.GetPluginInfoResponse{
		Name:          "inaccel",
		VendorVersion: plugin.driver.Version(),
	}

	return response, nil
}

func (plugin Kubelet) NodePublishVolume(ctx context.Context, request *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	logrus.Info("Kubelet/NodePublishVolume")

	response := &csi.NodePublishVolumeResponse{}

	mountpoint, err := plugin.driver.Create(request.VolumeId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(request.TargetPath, os.ModePerm); err != nil {
		return nil, err
	}

	if !graphdriver.NewDefaultChecker().IsMounted(request.TargetPath) {
		if err := mount.Mount(mountpoint, request.TargetPath, "none", "rbind"); err != nil {
			return nil, err
		}
	}

	return response, nil
}

func (plugin Kubelet) NodeUnpublishVolume(ctx context.Context, request *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	logrus.Info("Kubelet/NodeUnpublishVolume")

	response := &csi.NodeUnpublishVolumeResponse{}

	if graphdriver.NewDefaultChecker().IsMounted(request.TargetPath) {
		if err := mount.Unmount(request.TargetPath); err != nil {
			return nil, err
		}
	}

	if err := os.RemoveAll(request.TargetPath); err != nil {
		return nil, err
	}

	if err := plugin.driver.Release(request.VolumeId); err != nil {
		return nil, err
	}

	return response, nil
}

func (plugin Kubelet) NodeExpandVolume(ctx context.Context, request *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	logrus.Info("Kubelet/NodeExpandVolume")

	return nil, status.Error(codes.Unimplemented, "")
}

func (plugin Kubelet) NodeGetCapabilities(ctx context.Context, request *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	logrus.Info("Kubelet/NodeGetCapabilities")

	response := &csi.NodeGetCapabilitiesResponse{}

	return response, nil
}

func (plugin Kubelet) NodeGetInfo(ctx context.Context, request *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	logrus.Info("Kubelet/NodeGetInfo")

	response := &csi.NodeGetInfoResponse{}

	if nodeName, ok := os.LookupEnv("NODE_NAME"); ok {
		response.NodeId = nodeName
	} else {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, err
		}

		response.NodeId = hostname
	}

	return response, nil
}

func (plugin Kubelet) NodeGetVolumeStats(ctx context.Context, request *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	logrus.Info("Kubelet/NodeGetVolumeStats")

	return nil, status.Error(codes.Unimplemented, "")
}

func (plugin Kubelet) NodeStageVolume(ctx context.Context, request *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	logrus.Info("Kubelet/NodeStageVolume")

	return nil, status.Error(codes.Unimplemented, "")
}

func (plugin Kubelet) NodeUnstageVolume(ctx context.Context, request *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	logrus.Info("Kubelet/NodeUnstageVolume")

	return nil, status.Error(codes.Unimplemented, "")
}

func (plugin Kubelet) NotifyRegistrationStatus(ctx context.Context, request *pluginregistrationv1.RegistrationStatus) (*pluginregistrationv1.RegistrationStatusResponse, error) {
	logrus.Info("Kubelet/NotifyRegistrationStatus")

	response := &pluginregistrationv1.RegistrationStatusResponse{}

	if !request.PluginRegistered {
		logrus.Error(request.Error)
	}

	return response, nil
}

func (plugin Kubelet) Probe(ctx context.Context, request *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	logrus.Info("Kubelet/Probe")

	response := &csi.ProbeResponse{}

	return response, nil
}
