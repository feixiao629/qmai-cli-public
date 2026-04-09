package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// Format represents an output format
type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatCSV   Format = "csv"
)

// ParseFormat parses a format string, returns error if invalid
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON, nil
	case "table", "":
		return FormatTable, nil
	case "csv":
		return FormatCSV, nil
	default:
		return "", fmt.Errorf("unsupported format %q (use json, table, or csv)", s)
	}
}

// Formatter writes structured data in different formats
type Formatter struct {
	Writer io.Writer
	Format Format
}

// NewFormatter creates a new Formatter
func NewFormatter(w io.Writer, format Format) *Formatter {
	return &Formatter{Writer: w, Format: format}
}

// WriteJSON outputs data as pretty-printed JSON
func (f *Formatter) WriteJSON(v interface{}) error {
	enc := json.NewEncoder(f.Writer)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// WriteTable outputs data as an aligned table
func (f *Formatter) WriteTable(headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(f.Writer, 0, 4, 2, ' ', 0)

	// Header
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	// Separator
	seps := make([]string, len(headers))
	for i, h := range headers {
		seps[i] = strings.Repeat("-", len(h))
	}
	fmt.Fprintln(tw, strings.Join(seps, "\t"))

	// Rows
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	return tw.Flush()
}

// WriteCSV outputs data as CSV
func (f *Formatter) WriteCSV(headers []string, rows [][]string) error {
	w := csv.NewWriter(f.Writer)
	if err := w.Write(headers); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// Write outputs data in the configured format. For table/csv, provide headers and rows.
// For JSON, pass v as the data to marshal.
func (f *Formatter) Write(v interface{}, headers []string, rows [][]string) error {
	switch f.Format {
	case FormatJSON:
		return f.WriteJSON(v)
	case FormatCSV:
		return f.WriteCSV(headers, rows)
	default:
		return f.WriteTable(headers, rows)
	}
}
