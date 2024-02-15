package crusher

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"cloud.google.com/go/bigtable"
)

//func init() {
//	var Crusher Crusher
//	// Register a CloudEvent function with the Functions Framework
//	functions.HTTP("btDelete", Crusher.btDelete)
//}

func (content *Crusher) ReadWithFilter() ([]string, error) {
	//var w bytes.Buffer
	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, content.ProjectID, content.InstanceID)

	if err != nil {
		return nil, fmt.Errorf("bigtable.NewAdminClient: %w", err)
	}

	startTime := time.Unix(0, 0)
	endTime := time.Now().AddDate(0, 0, -content.Days)

	//TimestampRangeFilter returns a filter that matches any cells whose timestamp is within the given time bounds.
	timeFilter := bigtable.TimestampRangeFilter(startTime, endTime)
	content.Filter = bigtable.ChainFilters(bigtable.RowKeyFilter(content.KeyFilter), timeFilter)

	var keys []string

	tbl := client.Open(content.TableID)
	err = tbl.ReadRows(ctx, bigtable.RowRange{},
		func(row bigtable.Row) bool {
			keys = append(keys, row.Key())
			//printRow(&w, row)
			return true
		}, bigtable.RowFilter(content.Filter))

	if err = client.Close(); err != nil {
		return nil, fmt.Errorf("client.Close(): %w", err)
	}

	return keys, nil
}

func printRow(w io.Writer, row bigtable.Row) {
	log.Print(row.Key())
	fmt.Fprintf(w, "Reading data for %s:\n", row.Key())
	for columnFamily, cols := range row {
		fmt.Fprintf(w, "Column Family %s\n", columnFamily)
		for _, col := range cols {
			qualifier := col.Column[strings.IndexByte(col.Column, ':')+1:]
			fmt.Fprintf(w, "\t%s: %s @%d\n", qualifier, col.Value, col.Timestamp)
			log.Print(col.Timestamp.Time())
			//log.Printf("%d\n", temp)
			//log.Print(col.Timestamp)
		}
	}
	fmt.Fprintln(w)
}

func (content *Crusher) insertRows(projectID, instanceID string, rows []string) {
	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, projectID, instanceID)

	if err != nil {
		log.Info().Msgf("bigtable.NewAdminClient: %s", err)
	}
	tbl := client.Open("crusher")
	mut := bigtable.NewMutation()
	// To use numeric values that will later be incremented,
	// they need to be big-endian encoded as 64-bit integers.
	buf := new(bytes.Buffer)
	initialLinkCount := 1 // The initial number of links.
	if err := binary.Write(buf, binary.BigEndian, initialLinkCount); err != nil {
		// TODO: handle err.
	}
	mut.Set("links", "maps.google.com", bigtable.Now(), buf.Bytes())
	mut.Set("links", "golang.org", bigtable.Now(), buf.Bytes())

	for _, row := range rows {
		err = tbl.Apply(ctx, row, mut)
	}

	if err != nil {
		// TODO: handle err.
	}

	if err = client.Close(); err != nil {
		log.Info().Msgf("client.Close(): %s", err)
	}
}

func (content *Crusher) DeleteRows(rows []string) error {

	if !content.DryRun {
		ctx := context.Background()
		client, err := bigtable.NewClient(ctx, content.ProjectID, content.TableID)

		if err != nil {
			log.Info().Msgf("bigtable.NewAdminClient: %s", err)
		}

		tbl := client.Open(content.TableID)
		mut := bigtable.NewMutation()

		// To use numeric values that will later be incremented,
		// they need to be big-endian encoded as 64-bit integers.
		buf := new(bytes.Buffer)
		initialLinkCount := 1 // The initial number of links.
		if err := binary.Write(buf, binary.BigEndian, initialLinkCount); err != nil {
			return err
		}

		mut.DeleteRow()

		for _, row := range rows {
			err = tbl.Apply(ctx, row, mut)
		}

		if err != nil {
			return err
		}
	} else {
		log.Info().Msg("dry-run only")
	}

	return nil
}

func (content *Crusher) Clip() error {
	rows, err := content.ReadWithFilter()

	if err != nil {
		log.Info().Err(err)
		return nil
	}

	if len(rows) != 0 {
		err = content.DeleteRows(rows)
	} else {
		log.Info().Msg("nothing to delete")
		return nil
	}

	if err != nil {
		log.Info().Err(err)
		return nil
	}

	log.Info().Msgf("deleted %d rows", len(rows))

	return nil
}
