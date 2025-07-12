// internal/importer/excel_importer.go
package importer

import (
	"fmt"
	"mime/multipart"
)

type ExcelImporter struct{}

func NewExcelImporter() *ExcelImporter {
	return &ExcelImporter{}
}

func (ei *ExcelImporter) Import(file multipart.File) error {
	// TODO: Use Excel parsing (e.g., excelize)
	fmt.Println("Importing from Excel...")

	// logic here...
	return nil
}
