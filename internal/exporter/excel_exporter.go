package exporter

import (
	"fmt"

	"path/filepath"
	"sync"

	"github.com/xuri/excelize/v2"
)

// ExcelExporter implements Exporter for exporting data to Excel files using excelize.
// It manages file creation, row addition, and saving to disk.
type ExcelExporter struct {
	headers   []string       // Column headers for the Excel sheet
	file      *excelize.File // Excel file object
	sheetName string         // Name of the sheet to write to
	fileMutex sync.Mutex     // Mutex for concurrent access
	filename  string         // Output filename
}

// NewExcelExporter creates a new ExcelExporter with the given headers and sheet name.
// @Summary Create new ExcelExporter
// @Description Initializes a new ExcelExporter instance for exporting data to Excel files
// @Tags Exporter
// @Param headers body []string true "Column headers for Excel sheet"
// @Param sheetName body string true "Sheet name for Excel file"
// @Success 200 {object} ExcelExporter "ExcelExporter instance"
func NewExcelExporter(headers []string, sheetName string) *ExcelExporter {
	return &ExcelExporter{
		headers:   headers,
		file:      excelize.NewFile(),
		sheetName: sheetName,
	}
}

// Open initializes the Excel file and writes the headers to the first row.
// @Summary Open Excel file for export
// @Description Initializes the Excel file and writes headers to the first row
// @Tags Exporter
// @Param filename path string true "Output filename"
// @Success 200 {string} string "Excel file opened successfully"
// @Failure 500 {string} string "Failed to open Excel file"
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

	// Write headers to the first row of the sheet
	for i, h := range allHeaders {
		col, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s1", col)
		e.file.SetCellValue(sheet, cell, h)
	}
	return nil
}

// AddRow appends a row of data to the Excel sheet.
// Handles nullable types and writes each cell value.
// @Summary Add row to Excel sheet
// @Description Appends a row of data to the Excel sheet
// @Tags Exporter
// @Param row body []any true "Row data to append"
// @Success 200 {string} string "Row added successfully"
// @Failure 500 {string} string "Failed to add row"
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
		e.file.SetCellValue(sheet, cell, cellVal)
	}
	return nil
}

// Save writes the Excel file to disk in the specified local directory.
// Returns the full path to the saved file.
// @Summary Save Excel file
// @Description Writes the Excel file to disk in the specified local directory
// @Tags Exporter
// @Param localDir path string true "Local directory to save file"
// @Success 200 {string} string "Path to saved file"
// @Failure 500 {string} string "Failed to save file"
func (e *ExcelExporter) Save(localDir string) (string, error) {
	e.fileMutex.Lock()
	defer e.fileMutex.Unlock()

	path := filepath.Join(localDir, e.filename)
	if err := e.file.SaveAs(path); err != nil {
		return "", err
	}
	return path, nil
}
