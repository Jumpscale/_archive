package rest

import (
	"net/http"
	"net/url"

	"github.com/Jumpscale/aydostorex/fs"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("http-rest")
)

type Service struct {
	client        *http.Client
	store         *fs.Store
	backendStores []*url.URL
}

func NewService(store *fs.Store, client *http.Client, backendStores []*url.URL) *Service {
	if client == nil {
		client = http.DefaultClient
	}

	return &Service{
		client:        client,
		store:         store,
		backendStores: backendStores,
	}
}
