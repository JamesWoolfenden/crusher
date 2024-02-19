package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	crusher "github.com/JamesWoolfenden/crusher/src"
	"github.com/JamesWoolfenden/crusher/src/version"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"moul.io/banner"
)

func main() {
	fmt.Println(banner.Inline("crusher"))
	fmt.Println("version:", version.Version)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var Content crusher.Crusher

	app := &cli.App{
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:      "version",
				Aliases:   []string{"v"},
				Usage:     "Outputs the application version",
				UsageText: "Crusher version",
				Action: func(*cli.Context) error {
					return nil
				},
			},
			{
				Name:      "clip",
				Aliases:   []string{"c"},
				Usage:     "Compacts BigTable",
				UsageText: "crusher clip",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "instance",
						Aliases:     []string{"i"},
						Destination: &Content.InstanceID,
						Value:       "pangpt",
						Category:    "bigtable",
					},
					&cli.StringFlag{
						Name:        "table",
						Aliases:     []string{"t"},
						Destination: &Content.TableID,
						Value:       "pangpt",
						Category:    "bigtable",
					},
					&cli.StringFlag{
						Name:        "project",
						Aliases:     []string{"p"},
						Value:       "pangpt",
						Destination: &Content.ProjectID,
						Category:    "bigtable",
					},
					&cli.StringFlag{
						Name:        "keyfilter",
						Aliases:     []string{"k"},
						Destination: &Content.KeyFilter,
						Category:    "bigtable",
						Value:       ".*chat_histories$",
					},
					&cli.IntFlag{
						Name:        "days",
						Aliases:     []string{"d"},
						Destination: &Content.Days,
						Category:    "bigtable",
						Value:       180,
					},
					&cli.BoolFlag{
						Name:        "dry-run",
						Aliases:     []string{"r"},
						Destination: &Content.DryRun,
						Category:    "bigtable",
						Value:       false,
					},
				},
				Action: func(*cli.Context) error {
					Content.Clip()
					return nil
				},
			},
		},
		Name:     "crusher",
		Usage:    "crsuher clip",
		Compiled: time.Time{},
		Authors:  []*cli.Author{{Name: "James Woolfenden", Email: "jim.wolf@duck.com"}},
		Version:  version.Version,
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("crusher failure")
	}
}
