package main

import (
	crusher "crusher/src"
	"crusher/src/version"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"moul.io/banner"
)

func main() {
	fmt.Println(banner.Inline("crusher"))
	fmt.Println("version:", version.Version)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var content crusher.Crusher

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
					fmt.Println(version.Version)

					return nil
				},
			},
			{
				Name:      "clip",
				Aliases:   []string{"c"},
				Usage:     "Compacts BigTable ",
				UsageText: "crusher clip",
				Action: func(*cli.Context) error {
					content.ReadWithFilter()
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "instance",
						Aliases:     []string{"i"},
						Usage:       "instance",
						Destination: content.InstanceID,
						Category:    "bigtable",
					},
					&cli.StringFlag{
						Name:        "table",
						Aliases:     []string{"t"},
						Usage:       "table",
						Destination: content.TableID,
						Category:    "bigtable",
					},
					&cli.StringFlag{
						Name:        "project",
						Aliases:     []string{"p"},
						Usage:       "GCloudProject",
						DefaultText: "Pangpt",
						Destination: content.ProjectID,
						Category:    "bigtable",
					},
				},
			},
		},
		Name:     "crusher",
		Usage:    "AISB utility",
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

//func mainold() {
//	projectID := "pangpt"  // The Google Cloud Platform project ID
//	instanceID := "pangpt" // The Google Cloud Bigtable instance ID
//	tableID := "pangpt"    // The Google Cloud Bigtable table
//
//	// [END bigtable_quickstart]
//	// Override with -project, -instance, -table flags
//	flag.StringVar(&projectID, "project", projectID, "The Google Cloud Platform project ID.")
//	flag.StringVar(&instanceID, "instance", instanceID, "The Google Cloud Bigtable instance ID.")
//	flag.StringVar(&tableID, "table", tableID, "The Google Cloud Bigtable table ID.")
//	flag.Parse()
//
//	startTime := time.Unix(0, 0)
//	endTime := time.Now().AddDate(0, 0, -90)
//
//	//TimestampRangeFilter returns a filter that matches any cells whose timestamp is within the given time bounds.
//	timeFilter := bigtable.TimestampRangeFilter(startTime, endTime)
//	chatFilter := bigtable.RowKeyFilter(".*chat_histories$")
//
//	filter := bigtable.ChainFilters(chatFilter, timeFilter)
//	rows, err := content.ReadWithFilter(projectID, instanceID, tableID, filter)
//
//	if err != nil || rows == nil {
//		log.Fatalf("read failure %s", err)
//	}
//
//	err = btdelete.DeleteRows(projectID, instanceID, tableID, rows)
//	if err != nil {
//		log.Fatalf("delete failure %s", err)
//	}
//}
