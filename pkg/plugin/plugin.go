package plugin

type New func() Plugin

type Plugin interface {
	Close()
	IsClosed() bool
}
