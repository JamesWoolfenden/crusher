package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	crusher "github.com/JamesWoolfenden/crusher/src"
	"github.com/JamesWoolfenden/crusher/src/version"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"moul.io/banner"
)

func validateInput(projectID, instanceID, tableID string) error {
	if strings.TrimSpace(projectID) == "" {
		return fmt.Errorf("project ID cannot be empty")
	}
	if strings.TrimSpace(instanceID) == "" {
		return fmt.Errorf("instance ID cannot be empty")
	}
	if strings.TrimSpace(tableID) == "" {
		return fmt.Errorf("table ID cannot be empty")
	}
	return nil
}

func confirmDeletion(content *crusher.Crusher, rowCount int) bool {
	if content.DryRun {
		return true
	}

	fmt.Printf("\nWARNING: You are about to delete %d rows from:\n", rowCount)
	fmt.Printf("  Project:  %s\n", content.ProjectID)
	fmt.Printf("  Instance: %s\n", content.InstanceID)
	fmt.Printf("  Table:    %s\n", content.TableID)
	fmt.Printf("  Filter:   %s\n", content.KeyFilter)
	fmt.Printf("  Days:     >%d days old\n\n", content.Days)
	fmt.Print("Are you sure you want to continue? (yes/no): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Error().Err(err).Msg("failed to read user input")
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes" || response == "y"
}

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
						Value:       "",
						Category:    "bigtable",
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "table",
						Aliases:     []string{"t"},
						Destination: &Content.TableID,
						Value:       "",
						Category:    "bigtable",
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "project",
						Aliases:     []string{"p"},
						Value:       "",
						Destination: &Content.ProjectID,
						Category:    "bigtable",
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "keyfilter",
						Aliases:     []string{"k"},
						Destination: &Content.KeyFilter,
						Category:    "bigtable",
						Value:       ".*",
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
					&cli.BoolFlag{
						Name:     "yes",
						Aliases:  []string{"y"},
						Category: "bigtable",
						Value:    false,
						Usage:    "Skip confirmation prompt",
					},
				},
				Action: func(c *cli.Context) error {
					// Validate inputs
					if err := validateInput(Content.ProjectID, Content.InstanceID, Content.TableID); err != nil {
						return fmt.Errorf("validation failed: %w", err)
					}

					// Create context with timeout
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
					defer cancel()

					// First, do a dry-run to see what would be deleted
					originalDryRun := Content.DryRun
					Content.DryRun = true
					rowCount, err := Content.Clip(ctx)
					if err != nil {
						return fmt.Errorf("failed to preview deletions: %w", err)
					}

					if rowCount == 0 {
						log.Info().Msg("No rows to delete")
						return nil
					}

					// If not originally in dry-run mode, ask for confirmation
					if !originalDryRun {
						Content.DryRun = false
						skipConfirm := c.Bool("yes")

						if !skipConfirm && !confirmDeletion(&Content, rowCount) {
							log.Info().Msg("Operation cancelled by user")
							return nil
						}

						// Actually delete the rows
						_, err = Content.Clip(ctx)
						if err != nil {
							return fmt.Errorf("failed to delete rows: %w", err)
						}
					}

					return nil
				},
			},
		},
		Name:     "crusher",
		Usage:    "crusher clip",
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
