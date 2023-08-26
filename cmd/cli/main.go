package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main()  {
  app := &cli.App{
    Name: "cli-add",
    Usage: "For easily adding posts to sqlite db",
    Commands: []*cli.Command{
      {
        Name: "new",
        Aliases: []string{"n"},
        Usage: "Spawns a test html doc for putting in DB",
        Action: func(ctx *cli.Context) error {
          fmt.Println("Making an example form in dir...")
          // todo add file spawner
          return nil
        },
      },
      {
        Name: "tag-name",
        Aliases: []string{"tn"},
        Usage: "For updating tag names",
        Action: func(ctx *cli.Context) error {
          fmt.Println("Making an example form in dir...")

          return nil
        },
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
