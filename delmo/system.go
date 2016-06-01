package delmo

type System struct {
	config SystemConfig
}

func NewSystem(config SystemConfig) *System {
	return &System{config: config}
}

func (s *System) NewRuntime(name string) (Runtime, error) {
	return NewDockerCompose(s.config.CompleteFilePath, name)
}
