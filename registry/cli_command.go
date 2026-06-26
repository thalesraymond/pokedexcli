package registry

type CLICommand struct {
	name        string
	description string
	callback    func() error
}

type CLICommandInterface interface {
	GetName() string
	GetDescription() string
	Execute() error
}

func (c *CLICommand) GetName() string {
	return c.name
}

func (c *CLICommand) GetDescription() string {
	return c.description
}

func (c *CLICommand) Execute() error {
	return c.callback()
}
