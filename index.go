package store_default

import (
	"github.com/infrago/store"
)

func Driver() store.Driver {
	return &defaultDriver{}
}

func init() {
	store.Register("default", Driver())
}
