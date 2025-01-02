package main

// TODO: include connection with redis
// TODO: include connection with mongodb
// TODO: include connection with Oracle
// TODO: include connection with Mysql
// TODO: include authentication using JWT
import (
	"fmt"
	"log"
	"os"

	cmd "github.com/marinellirubens/dbwrapper/cmd/app"
	cf "github.com/marinellirubens/dbwrapper/internal/config"
	cli "github.com/urfave/cli/v2"
)

const VERSION = "1.0.0"

func main() {
	var cfgPath string

	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config_file",
				Usage:       "Path for the configuration file",
				Aliases:     []string{"f"},
				Value:       cf.DefaultCfgFilePath,
				Destination: &cfgPath,
			},
			&cli.BoolFlag{
				Name:    "version",
				Usage:   "Path for the configuration file",
				Aliases: []string{"v"},
				Value:   false,
				Action: func(ctx *cli.Context, b bool) error {
					fmt.Printf("Version: %v\n", VERSION)
					os.Exit(0)
					return nil
				},
			},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Println("path: ", cfgPath)
			cmd.RunServer(cfgPath)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
