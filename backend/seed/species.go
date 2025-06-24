package seed

import (
	"context"
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/anish-chanda/ferna/db"
)

//go:embed species.csv
var speciesCSV embed.FS

func SeedSpecies(ctx context.Context, database db.Database) error {
	f, err := speciesCSV.Open("species.csv")
	if err != nil {
		return fmt.Errorf("open embedded CSV: %w", err)
	}
	defer f.Close()
	rdr := csv.NewReader(f)
	// skip header
	if _, err := rdr.Read(); err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	for {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("csv parse: %w", err)
		}
		freq, _ := strconv.Atoi(record[2])

		// Use INSERT OR IGNORE to handle duplicates as scientific name is unique
		if _, err := database.ExecContext(ctx,
			`INSERT OR IGNORE INTO species(common_name, scientific_name, default_watering_frequency_days) VALUES (?, ?, ?)`,
			record[0], record[1], freq); err != nil {
			return fmt.Errorf("insert species: %w", err)
		}
	}
	return nil
}
