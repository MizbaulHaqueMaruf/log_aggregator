package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mizbaulhaquemaruf/log_aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully!")
	printUser(user)
	return nil
}

func printUser(user database.User) {
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Created At: %s\n", user.CreatedAt)
	fmt.Printf("Updated At: %s\n", user.UpdatedAt)
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting all users: %w", err)
	}

	fmt.Println("All users deleted successfully!")
	return nil
}

func handlerList(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error listing users: %w", err)
	}

	currentUser := s.cfg.CurrentUserName

	for _, user := range users {
		curr := user.Name
		if currentUser == user.Name {
			curr += " (current)"
		}
		fmt.Println("* " + curr)
	}

	return nil
}
