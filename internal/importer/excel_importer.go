package importer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/xuri/excelize/v2"
)

// RowHandlerFunc defines how to process a row.
type RowHandlerFunc func(data map[string]string, userID int) error

type ExcelImporter struct {
	ExpectedHeaders []string
	HandleRow       RowHandlerFunc
	Mutex           *sync.Mutex
}

func NewExcelImporter(headers []string, handler RowHandlerFunc) *ExcelImporter {
	return &ExcelImporter{
		ExpectedHeaders: headers,
		HandleRow:       handler,
		Mutex:           &sync.Mutex{},
	}
}
func (e *ExcelImporter) Import(filePath string, _ []string,userID int) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return fmt.Errorf("read rows error: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("no data rows found")
	}

	headers := rows[0]
	if err := e.ValidateHeaders(headers); err != nil {
		return fmt.Errorf("header validation failed: %w", err)
	}

	// Stream-process each row
	for _, row := range rows[1:] {
		if len(row) == 0 {
			continue
		}

		record := make(map[string]string)
		for i, cell := range row {
			if i < len(headers) {
				record[headers[i]] = cell
			}
		}

		e.Mutex.Lock()
		err := e.HandleRow(record, userID)
		e.Mutex.Unlock()

		if err != nil {
			return fmt.Errorf("row handler failed: %w", err)
		}
	}

	return nil
}

func (e *ExcelImporter) ValidateHeaders(actual []string) error {
	var missing []string
	for _, expected := range e.ExpectedHeaders {
		found := false
		for _, got := range actual {
			if strings.EqualFold(expected, got) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, expected)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required headers: %v", missing)
	}
	return nil
}
