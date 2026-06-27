package registry

type CLICommand struct {
	name        string
	description string
	callback    func(cfg *PokedexContext) error
}

type CLICommandInterface interface {
	GetName() string
	GetDescription() string
	Execute(pokedexContext *PokedexContext) error
}

func (c *CLICommand) GetName() string {
	return c.name
}

func (c *CLICommand) GetDescription() string {
	return c.description
}

func (c *CLICommand) Execute(pokedexContext *PokedexContext) error {
	return c.callback(pokedexContext)
}
