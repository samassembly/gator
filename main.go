package main

import (
	"fmt"
	"log"
	"os"
	"github.com/samassembly/gator/internal/config"
	"github.com/samassembly/gator/internal/app"
)

type state struct {
	Configuration *Config
}

func main() {
	s := state{}
	s.Configuration, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", s.Configuration)

	cmds := app.Commands{}
	cmds.Cmds = make(map[string]func(*state, app.command) error)

	cmds.Register("login", app.handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Error, no command provided")
		return
	}

	command := app.Command{
		Name: args[0]
		Arguments: args[1:]
	}

	err := app.run(s, command)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	
}