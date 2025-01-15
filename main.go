package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mizbaulhaquemaruf/log_aggregator/internal/config"
	"github.com/mizbaulhaquemaruf/log_aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerList)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{
		Name: cmdName,
		Args: cmdArgs,
	})
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("Read config: %v\n", cfg)

	// err = cfg.SetUser("lane")
	// if err != nil {
	// 	fmt.Errorf("user name cannot be set, try again. Error: %w", err)
	// }

	// cfg, err = config.Read()
	// if err != nil {
	// 	log.Fatalf("error reading config: %v", err)
	// }

	// fmt.Printf("Read config again %v \n", cfg)
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
