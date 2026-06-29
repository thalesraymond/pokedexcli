package registry

type CLICommand struct {
	name        string
	description string
	callback    func(cfg *PokedexContext, args ...string) error
}

type CLICommandInterface interface {
	GetName() string
	GetDescription() string
	Execute(context *PokedexContext, args ...string) error
}

func (c *CLICommand) GetName() string {
	return c.name
}

func (c *CLICommand) GetDescription() string {
	return c.description
}

func (c *CLICommand) Execute(context *PokedexContext, args ...string) error {
	return c.callback(context, args...)
}
