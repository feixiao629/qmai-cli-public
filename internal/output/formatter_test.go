package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	data := map[string]string{"name": "Alice", "role": "admin"}
	if err := f.WriteJSON(data); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	// Output should be valid JSON
	var parsed map[string]string
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	if parsed["name"] != "Alice" {
		t.Errorf("name = %q, want %q", parsed["name"], "Alice")
	}
	if parsed["role"] != "admin" {
		t.Errorf("role = %q, want %q", parsed["role"], "admin")
	}

	// Should be pretty-printed (indented)
	if !strings.Contains(buf.String(), "  ") {
		t.Error("expected pretty-printed JSON with indentation")
	}
}

func TestWriteJSON_Slice(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	data := []map[string]int{{"x": 1}, {"x": 2}}
	if err := f.WriteJSON(data); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	var parsed []map[string]int
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(parsed) != 2 {
		t.Errorf("expected 2 items, got %d", len(parsed))
	}
}

func TestWriteTable(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatTable)

	headers := []string{"ID", "Name", "Status"}
	rows := [][]string{
		{"1", "Alice", "active"},
		{"2", "Bob", "inactive"},
	}

	if err := f.WriteTable(headers, rows); err != nil {
		t.Fatalf("WriteTable failed: %v", err)
	}

	out := buf.String()

	// Headers should be present
	if !strings.Contains(out, "ID") {
		t.Error("missing header 'ID'")
	}
	if !strings.Contains(out, "Name") {
		t.Error("missing header 'Name'")
	}
	if !strings.Contains(out, "Status") {
		t.Error("missing header 'Status'")
	}

	// Data rows should be present
	if !strings.Contains(out, "Alice") {
		t.Error("missing data 'Alice'")
	}
	if !strings.Contains(out, "Bob") {
		t.Error("missing data 'Bob'")
	}

	// Should have separator line with dashes
	if !strings.Contains(out, "--") {
		t.Error("missing separator dashes")
	}

	// Should have at least 4 lines (header, separator, 2 data rows)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 4 {
		t.Errorf("expected at least 4 lines, got %d", len(lines))
	}
}

func TestWriteCSV(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatCSV)

	headers := []string{"ID", "Name", "Email"}
	rows := [][]string{
		{"1", "Alice", "alice@example.com"},
		{"2", "Bob", "bob@example.com"},
	}

	if err := f.WriteCSV(headers, rows); err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}

	// Header line
	if lines[0] != "ID,Name,Email" {
		t.Errorf("header line = %q, want %q", lines[0], "ID,Name,Email")
	}

	// First data row
	if lines[1] != "1,Alice,alice@example.com" {
		t.Errorf("row 1 = %q, want %q", lines[1], "1,Alice,alice@example.com")
	}

	// Second data row
	if lines[2] != "2,Bob,bob@example.com" {
		t.Errorf("row 2 = %q, want %q", lines[2], "2,Bob,bob@example.com")
	}
}

func TestWriteCSV_WithCommasInData(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatCSV)

	headers := []string{"Name", "Description"}
	rows := [][]string{
		{"Alice", "Developer, Senior"},
	}

	if err := f.WriteCSV(headers, rows); err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	out := buf.String()
	// The field with a comma should be quoted
	if !strings.Contains(out, `"Developer, Senior"`) {
		t.Errorf("expected quoted field with comma, got: %s", out)
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    Format
		wantErr bool
	}{
		{"json", FormatJSON, false},
		{"JSON", FormatJSON, false},
		{"table", FormatTable, false},
		{"TABLE", FormatTable, false},
		{"", FormatTable, false}, // empty defaults to table
		{"csv", FormatCSV, false},
		{"CSV", FormatCSV, false},
		{"xml", "", true},
		{"yaml", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormat(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWrite_DispatchesToCorrectMethod(t *testing.T) {
	headers := []string{"K", "V"}
	rows := [][]string{{"a", "1"}}
	data := map[string]string{"a": "1"}

	t.Run("json", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewFormatter(&buf, FormatJSON)
		if err := f.Write(data, headers, rows); err != nil {
			t.Fatalf("Write failed: %v", err)
		}
		// Should be valid JSON
		var parsed map[string]string
		if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
			t.Errorf("expected JSON output, got: %s", buf.String())
		}
	})

	t.Run("table", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewFormatter(&buf, FormatTable)
		if err := f.Write(data, headers, rows); err != nil {
			t.Fatalf("Write failed: %v", err)
		}
		if !strings.Contains(buf.String(), "K") {
			t.Error("expected table output with headers")
		}
	})

	t.Run("csv", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewFormatter(&buf, FormatCSV)
		if err := f.Write(data, headers, rows); err != nil {
			t.Fatalf("Write failed: %v", err)
		}
		if !strings.Contains(buf.String(), "K,V") {
			t.Errorf("expected CSV header, got: %s", buf.String())
		}
	})
}
