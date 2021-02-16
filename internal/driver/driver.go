package driver

type Driver interface {
	Create(string) (string, error)
	Get(string) (string, error)
	List() (map[string]string, error)
	Release(string) error
	Version() string
}
