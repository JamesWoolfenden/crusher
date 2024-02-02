package crusher

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	var Crusher Crusher
	// Register a CloudEvent function with the Functions Framework
	functions.HTTP("btDelete", Crusher.btDelete)
}

func (content *Crusher) ReadWithFilter() ([]string, error) {
	var w bytes.Buffer
	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, *content.ProjectID, *content.InstanceID)

	var keys []string

	if err != nil {
		return nil, fmt.Errorf("bigtable.NewAdminClient: %w", err)
	}
	tbl := client.Open(*content.TableID)
	err = tbl.ReadRows(ctx, bigtable.RowRange{},
		func(row bigtable.Row) bool {
			keys = append(keys, row.Key())
			printRow(&w, row)
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
		log.Println("bigtable.NewAdminClient: %w", err)
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
		log.Println("client.Close(): %w", err)
	}
}

func (content *Crusher) DeleteRows(rows []string) error {
	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, *content.ProjectID, *content.TableID)

	if err != nil {
		log.Println("bigtable.NewAdminClient: %w", err)
	}
	tbl := client.Open(*content.TableID)
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
	return nil
}

func (content *Crusher) btDelete(w http.ResponseWriter, r *http.Request) {
	//ctx := context.Background()
	//credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//projectID := credentials.ProjectID // The Google Cloud Platform project ID
	//instanceID := "pangpt"             // The Google Cloud Bigtable instance ID
	//tableID := "pangpt"                // The Google Cloud Bigtable table

	//// [END bigtable_quickstart]
	//// Override with -project, -instance, -table flags
	//flag.StringVar(&projectID, "project", projectID, "The Google Cloud Platform project ID.")
	//flag.StringVar(&instanceID, "instance", instanceID, "The Google Cloud Bigtable instance ID.")
	//flag.StringVar(&tableID, "table", tableID, "The Google Cloud Bigtable table ID.")
	//flag.Parse()

	startTime := time.Unix(0, 0)
	endTime := time.Now().AddDate(0, 0, -90)

	//TimestampRangeFilter returns a filter that matches any cells whose timestamp is within the given time bounds.
	timeFilter := bigtable.TimestampRangeFilter(startTime, endTime)
	chatFilter := bigtable.RowKeyFilter(".*chat_histories$")

	content.Filter = bigtable.ChainFilters(chatFilter, timeFilter)
	rows, err := content.ReadWithFilter()

	if err != nil || rows == nil {
		log.Fatalf("read failure %s", err)
		return
	}

	err = content.DeleteRows(rows)
	if err != nil {
		log.Fatalf("delete failure %s", err)
		return
	}

	fmt.Fprintln(w, "ok")
}
