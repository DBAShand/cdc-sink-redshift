package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const resolvedTableSchema = `
CREATE TABLE IF NOT EXISTS %s (
	endpoint STRING PRIMARY KEY,
	nanos INT NOT NULL,
	logical INT NOT NULL
)
`

// Make this an option?
const resolvedTableName = `_release`

const resolvedTableQuery = `SELECT endpoint, nanos, logical FROM %s WHERE endpoint = $1`

const resolvedTableWrite = `UPSERT INTO %s (endpoint, nanos, logical) VALUES ($1, $2, $3)`

func resolvedFullTableName() string {
	return fmt.Sprintf("%s.%s", *sinkDB, resolvedTableName)
}

// CreateResolvedTable creates a release table if none exists.
func CreateResolvedTable(db *sql.DB) error {
	// Needs retry.
	_, err := db.Exec(fmt.Sprintf(resolvedTableSchema, resolvedFullTableName()))
	return err
}

// ResolvedLine is used to parse a json line in the request body of a resolved
// message.
type ResolvedLine struct {
	// These are use for parsing the resolved line.
	Resolved string `json:"resolved"`

	// There are used for storing back into the resolved table.
	nanos    int64
	logical  int
	endpoint string
}

func parseResolvedLine(rawBytes []byte, endpoint string) (ResolvedLine, error) {
	resolvedLine := ResolvedLine{
		endpoint: endpoint,
	}
	json.Unmarshal(rawBytes, &resolvedLine)
	log.Printf("resolved line: %s", string(rawBytes))

	// Prase the timestamp into nanos and logical.
	var err error
	resolvedLine.nanos, resolvedLine.logical, err = parseSplitTimestamp(resolvedLine.Resolved)
	if err != nil {
		return ResolvedLine{}, err
	}
	if resolvedLine.nanos == 0 {
		return ResolvedLine{}, fmt.Errorf("no nano component to the 'updated' timestamp field")
	}

	log.Printf("resolved: %+v", resolvedLine)

	return resolvedLine, nil
}

// getPreviousResolvedTimestamp returns the last recorded resolved for a
// specific endpoint.
func getPreviousResolved(db *sql.DB, endpoint string) (ResolvedLine, error) {
	// Needs retry.
	row := db.QueryRow(fmt.Sprintf(resolvedTableQuery, resolvedFullTableName()), endpoint)
	var resolvedLine ResolvedLine
	err := row.Scan(&(resolvedLine.endpoint), &(resolvedLine.nanos), &(resolvedLine.logical))
	switch err {
	case sql.ErrNoRows:
		// No line exists yet, go back to the start of time.
		return ResolvedLine{endpoint: endpoint}, nil
	case nil:
		// Found the line.
		return resolvedLine, nil
	default:
		return ResolvedLine{}, err
	}
}

// Writes the updated timestamp to the resolved table.
func (rl ResolvedLine) writeUpdated(db *sql.DB) error {
	// Needs retry.
	_, err := db.Exec(fmt.Sprintf(resolvedTableWrite, resolvedFullTableName()),
		rl.endpoint, rl.nanos, rl.logical,
	)
	return err
}
