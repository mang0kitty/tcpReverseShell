package rsh

type Module interface {
	Type() string
	Request() error
	Execute() error
}

type runModule struct {
	Program string
}

func (m *runModule) Type() string {
	return "run"
}

func (m *runModule) Request() error {
	return nil
}

func (m *runModule) Execute() error {
	return nil
}

type downloadModule struct {
	SrcFile  string
	DestFile string
}

func (m *downloadModule) Type() string {
	return "download"
}

func (m *downloadModule) Request() error {
	return nil
}

func (m *downloadModule) Execute() error {
	return nil
}
