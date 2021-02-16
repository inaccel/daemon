package driver

import (
	"sync"

	"github.com/inaccel/daemon/pkg"
	"github.com/sirupsen/logrus"
)

type InAccel struct {
	d sync.Mutex

	config  pkg.Config
	version string
}

func NewInAccel(config pkg.Config, version string) (*InAccel, error) {
	if _, err := config.Resources.Create(config.Root, ""); err != nil {
		return nil, err
	}

	return &InAccel{
		config:  config,
		version: version,
	}, nil
}

func (inaccel *InAccel) Create(name string) (string, error) {
	inaccel.d.Lock()
	defer inaccel.d.Unlock()

	create, err := inaccel.config.Resources.Create(inaccel.config.Root, name)

	logrus.WithFields(logrus.Fields{
		"error":  err,
		"name":   name,
		"return": create,
	}).Debug("resources.create")

	return create, err
}

func (inaccel *InAccel) Get(name string) (string, error) {
	inaccel.d.Lock()
	defer inaccel.d.Unlock()

	get, err := inaccel.config.Resources.Get(inaccel.config.Root, name)

	logrus.WithFields(logrus.Fields{
		"error":  err,
		"name":   name,
		"return": get,
	}).Debug("resources.get")

	return get, err
}

func (inaccel *InAccel) List() (map[string]string, error) {
	inaccel.d.Lock()
	defer inaccel.d.Unlock()

	list, err := inaccel.config.Resources.List(inaccel.config.Root)

	logrus.WithFields(logrus.Fields{
		"error":  err,
		"return": list,
	}).Debug("resources.list")

	return list, err
}

func (inaccel *InAccel) Release(name string) error {
	inaccel.d.Lock()
	defer inaccel.d.Unlock()

	err := inaccel.config.Resources.Release(inaccel.config.Root, name)

	logrus.WithFields(logrus.Fields{
		"error": err,
		"name":  name,
	}).Debug("resources.release")

	return err
}

func (inaccel *InAccel) Version() string {
	return inaccel.version
}
