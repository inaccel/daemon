package plugins

import (
	"context"
	"strings"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/inaccel/daemon/internal/driver"
	"github.com/inaccel/daemon/pkg/plugin"
	"github.com/sirupsen/logrus"
)

type Docker struct {
	path string

	driver driver.Driver
	plugin.Plugin
}

func NewDocker(ctx context.Context, driver driver.Driver) plugin.Plugin {
	ctx, cancel := context.WithCancel(ctx)

	docker := &Docker{
		path: "/run/docker/plugins/inaccel.sock",
	}

	docker.driver = driver

	docker.Plugin = plugin.Base(func() {
		if listener, err := listen(docker.path); err == nil {
			go func() {
				<-ctx.Done()

				listener.Close()
			}()

			server := volume.NewHandler(docker)

			server.Serve(listener)
		} else {
			logrus.Error(err)
		}
	}, cancel)

	return docker
}

func (plugin Docker) Capabilities() *volume.CapabilitiesResponse {
	logrus.Info("Docker/Capabilities")

	response := &volume.CapabilitiesResponse{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}

	return response
}

func (plugin Docker) Create(request *volume.CreateRequest) error {
	logrus.Info("Docker/Create")

	var name string
	if !strings.EqualFold(request.Name, "host") {
		name = request.Name
	}

	_, err := plugin.driver.Create(name)

	return err
}

func (plugin Docker) Get(request *volume.GetRequest) (*volume.GetResponse, error) {
	logrus.Info("Docker/Get")

	response := &volume.GetResponse{}

	if strings.EqualFold(request.Name, "inaccel") {
		response.Volume = &volume.Volume{
			Name: request.Name,
			Status: map[string]interface{}{
				"version": plugin.driver.Version(),
			},
		}
	} else {
		var name string
		if !strings.EqualFold(request.Name, "host") {
			name = request.Name
		}

		mountpoint, err := plugin.driver.Get(name)
		if err != nil {
			return nil, err
		}

		response.Volume = &volume.Volume{
			Name:       request.Name,
			Mountpoint: mountpoint,
			Status: map[string]interface{}{
				"version": plugin.driver.Version(),
			},
		}
	}

	return response, nil
}

func (plugin Docker) List() (*volume.ListResponse, error) {
	logrus.Info("Docker/List")

	response := &volume.ListResponse{}

	volumes, err := plugin.driver.List()
	if err != nil {
		return nil, err
	}

	for name, mountpoint := range volumes {
		response.Volumes = append(response.Volumes, &volume.Volume{
			Name:       name,
			Mountpoint: mountpoint,
			Status: map[string]interface{}{
				"version": plugin.driver.Version(),
			},
		})
	}

	return response, nil
}

func (plugin Docker) Mount(request *volume.MountRequest) (*volume.MountResponse, error) {
	logrus.Info("Docker/Mount")

	response := &volume.MountResponse{}

	if strings.EqualFold(request.Name, "inaccel") {
		mountpoint, err := plugin.driver.Create(request.ID)
		if err != nil {
			return nil, err
		}

		response.Mountpoint = mountpoint
	} else {
		var name string
		if !strings.EqualFold(request.Name, "host") {
			name = request.Name
		}

		mountpoint, err := plugin.driver.Get(name)
		if err != nil {
			return nil, err
		}

		response.Mountpoint = mountpoint
	}

	return response, nil
}

func (plugin Docker) Path(request *volume.PathRequest) (*volume.PathResponse, error) {
	logrus.Info("Docker/Path")

	response := &volume.PathResponse{}

	if !strings.EqualFold(request.Name, "inaccel") {
		var name string
		if !strings.EqualFold(request.Name, "host") {
			name = request.Name
		}

		mountpoint, err := plugin.driver.Get(name)
		if err != nil {
			return nil, err
		}

		response.Mountpoint = mountpoint
	}

	return response, nil
}

func (plugin Docker) Remove(request *volume.RemoveRequest) error {
	logrus.Info("Docker/Remove")

	var name string
	if !strings.EqualFold(request.Name, "host") {
		name = request.Name
	}

	err := plugin.driver.Release(name)

	return err
}

func (plugin Docker) Unmount(request *volume.UnmountRequest) error {
	logrus.Info("Docker/Unmount")

	if strings.EqualFold(request.Name, "inaccel") {
		err := plugin.driver.Release(request.ID)

		return err
	}

	var name string
	if !strings.EqualFold(request.Name, "host") {
		name = request.Name
	}

	_, err := plugin.driver.Get(name)

	return err
}
