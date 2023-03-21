package store_default

import (
	"github.com/infrago/infra"
	"github.com/infrago/store"
)

func Driver() store.Driver {
	return &defaultDriver{}
}

func init() {
	infra.Register("default", Driver())
}
