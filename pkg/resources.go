package pkg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/inaccel/daemon/pkg/tmpfs"
)

type Resource struct {
	Mode  string `json:"mode,omitempty" validate:"required,number,max=4,startswith=0,excludesall=89"`
	Tmpfs bool   `json:"tmpfs,omitempty" validate:"required_with=Huge Size"`
	Huge  string `json:"huge,omitempty" validate:"omitempty,oneof=always never within_size"`
	Size  string `json:"size,omitempty" validate:"omitempty,min=2,max=4,endswith=%"`
}

type Resources map[string]*Resource

func (resources Resources) magic() string {
	return "_" + strings.ToLower(reflect.TypeOf(resources).Name())
}

func (resources Resources) Create(root, name string) (string, error) {
	namespace := filepath.Join(root, "."+name)
	mountpoint := filepath.Join(namespace, resources.magic())

	if err := os.MkdirAll(namespace, 0700); err != nil {
		return "", err
	}

	if err := os.Chmod(root, 0755); err != nil {
		return "", err
	}

	if err := os.MkdirAll(mountpoint, 0755); err != nil {
		return "", err
	}

	for id, resource := range resources {
		path := filepath.Join(mountpoint, id)

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return "", err
		}

		if resource.Tmpfs {
			if err := tmpfs.Mount(path, tmpfs.Options{
				Huge: resource.Huge,
				Size: resource.Size,
			}); err != nil {
				return "", err
			}
		}

		if mode, err := strconv.ParseUint(resource.Mode, 8, 32); err != nil {
			return "", err
		} else if err = os.Chmod(path, os.FileMode(mode)); err != nil {
			return "", err
		}

		if err := os.Symlink(filepath.Join(resources.magic(), id), filepath.Join(namespace, id)); err != nil && !os.IsExist(err) {
			return "", err
		}
	}

	return mountpoint, nil
}

func (resources Resources) Get(root string, name string) (string, error) {
	if _, err := os.Stat(filepath.Join(root, "."+name, resources.magic())); err != nil {
		return "", err
	}

	return resources.Create(root, name)
}

func (resources Resources) List(root string) (map[string]string, error) {
	mountpoints, err := filepath.Glob(filepath.Join(root, ".*", resources.magic()))
	if err != nil {
		return nil, err
	}

	names := map[string]string{}

	for _, mountpoint := range mountpoints {
		name := strings.TrimPrefix(filepath.Base(filepath.Dir(mountpoint)), ".")

		if _, err := resources.Create(root, name); err == nil {
			names[name] = mountpoint
		}
	}

	return names, nil
}

func (resources Resources) Release(root, name string) error {
	namespace := filepath.Join(root, "."+name)
	mountpoint := filepath.Join(namespace, resources.magic())

	for id, resource := range resources {
		path := filepath.Join(mountpoint, id)

		if resource.Tmpfs {
			if err := tmpfs.Unmount(path); err != nil {
				return err
			}
		}
	}

	return os.RemoveAll(namespace)
}

func (resources Resources) UnmarshalJSON(data []byte) error {
	var v map[string]json.RawMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for id, rawMessage := range v {
		if string(rawMessage) == "null" {
			delete(resources, id)
		} else if resource, ok := resources[id]; ok {
			if err := json.Unmarshal(rawMessage, resource); err != nil {
				return err
			}
		} else {
			if err := json.Unmarshal(rawMessage, &resource); err != nil {
				return err
			}

			resources[id] = resource
		}
	}

	return nil
}
