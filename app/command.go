package app

import "fmt"

type Command struct {
	Name string
	Arguments []string
}

type Commands struct {
	Cmds map[string]func(*state, Command) error
}

func (c *Commands) run(s *state, cmd Command) error {
	_, ok := c.Cmds[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return c.Cmds[cmd.Name](s, cmd)
}

func (c *Commands) register(name string, f func(*state, Command) error) {
	// register a new signature/command
	if c.Cmds == nil {
		c.Cmds = make(map[string]func(*state, Command) error)
	}
	c.Cmds[name] = f
}