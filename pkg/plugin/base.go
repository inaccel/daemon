package plugin

type base struct {
	cancel  func()
	channel chan struct{}
}

func Base(run func(), cancel func()) Plugin {
	plugin := &base{
		cancel: cancel,
	}

	plugin.channel = make(chan struct{})
	go func() {
		defer close(plugin.channel)

		run()
	}()

	return plugin
}

func (plugin base) Close() {
	plugin.cancel()

	<-plugin.channel
}

func (plugin base) IsClosed() bool {
	select {
	case <-plugin.channel:
		return true
	default:
		return false
	}
}
