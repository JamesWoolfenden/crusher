package crusher

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"cloud.google.com/go/bigtable"
)

func (content *Crusher) ReadWithFilter(ctx context.Context) ([]string, error) {
	client, err := bigtable.NewClient(ctx, content.ProjectID, content.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("bigtable.NewClient: %w", err)
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("failed to close bigtable client")
		}
	}()

	startTime := time.Unix(0, 0)
	endTime := time.Now().AddDate(0, 0, -content.Days)

	// TimestampRangeFilter returns a filter that matches any cells whose timestamp is within the given time bounds.
	timeFilter := bigtable.TimestampRangeFilter(startTime, endTime)
	content.Filter = bigtable.ChainFilters(bigtable.RowKeyFilter(content.KeyFilter), timeFilter)

	var keys []string

	tbl := client.Open(content.TableID)
	err = tbl.ReadRows(ctx, bigtable.RowRange{},
		func(row bigtable.Row) bool {
			keys = append(keys, row.Key())
			return true
		}, bigtable.RowFilter(content.Filter))

	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	return keys, nil
}

func (content *Crusher) DeleteRows(ctx context.Context, rows []string) error {
	if content.DryRun {
		log.Info().Msgf("dry-run: would delete %d rows", len(rows))
		return nil
	}

	client, err := bigtable.NewClient(ctx, content.ProjectID, content.InstanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewClient: %w", err)
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("failed to close bigtable client")
		}
	}()

	tbl := client.Open(content.TableID)
	mut := bigtable.NewMutation()
	mut.DeleteRow()

	for _, row := range rows {
		err = tbl.Apply(ctx, row, mut)
		if err != nil {
			return fmt.Errorf("failed to delete row %s: %w", row, err)
		}
	}

	return nil
}

func (content *Crusher) Clip(ctx context.Context) (int, error) {
	rows, err := content.ReadWithFilter(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to read rows")
		return 0, fmt.Errorf("failed to read rows: %w", err)
	}

	if len(rows) == 0 {
		log.Info().Msg("nothing to delete")
		return 0, nil
	}

	err = content.DeleteRows(ctx, rows)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete rows")
		return 0, fmt.Errorf("failed to delete rows: %w", err)
	}

	if content.DryRun {
		log.Info().Msgf("dry-run: found %d rows to delete", len(rows))
	} else {
		log.Info().Msgf("deleted %d rows", len(rows))
	}

	return len(rows), nil
}
