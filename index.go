package storage_default

import (
	"github.com/infrago/infra"
	"github.com/infrago/storage"
)

func Driver() storage.Driver {
	return &defaultDriver{}
}

func init() {
	infra.Register("default", Driver())
}
