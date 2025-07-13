package exporter

import (
	"fmt"

	"path/filepath"
	"sync"

	"github.com/xuri/excelize/v2"
)

type ExcelExporter struct {
	headers   []string
	file      *excelize.File
	sheetName string
	fileMutex sync.Mutex
	filename  string
}

func NewExcelExporter(headers []string, sheetName string) *ExcelExporter {
	return &ExcelExporter{
		headers: headers,
		file:    excelize.NewFile(),
	}
}

func (e *ExcelExporter) Open(filename string) error {
	e.fileMutex.Lock()
	defer e.fileMutex.Unlock()

	e.filename = filename
	sheet := e.file.GetSheetName(0)
	if sheet == "" {
		sheet = e.sheetName
		e.file.NewSheet(sheet)
	}
	allHeaders := e.headers
	allHeaders = append([]string{"id"}, e.headers...)
	allHeaders = append(allHeaders, []string{"created_at", "updated_at"}...)

	// Add headers to the first row
	for i, h := range allHeaders {
		col, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s1", col)
		e.file.SetCellValue(sheet, cell, h)
	}
	return nil
}

func (e *ExcelExporter) AddRow(row []any) error {
	e.fileMutex.Lock()
	defer e.fileMutex.Unlock()

	sheet := e.file.GetSheetName(0)

	rows, err := e.file.GetRows(sheet)
	if err != nil {
		return err
	}
	rowIndex := len(rows) + 1

	for i, cellVal := range row {
		col, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s%d", col, rowIndex)
		// Handle sql.NullString or other nullable types
		e.file.SetCellValue(sheet, cell, cellVal)
	}
	return nil
}

func (e *ExcelExporter) Save(localDir string) (string, error) {
	e.fileMutex.Lock()
	defer e.fileMutex.Unlock()

	path := filepath.Join(localDir, e.filename)
	if err := e.file.SaveAs(path); err != nil {
		return "", err
	}
	return path, nil
}
