package server

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/trustacks/trustacks/pkg/api"
)

// New creates a new server instance.
func New(host, port string) {
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("starting server on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
