package app

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("no username provided")
	}

	username := cmd.Arguments[0]
	err := s.Configuration.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set")
}