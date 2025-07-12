package importer

import "mime/multipart"

type Importer interface {
	Import(file multipart.File) error
}