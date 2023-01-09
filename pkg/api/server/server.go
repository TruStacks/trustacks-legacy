package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/trustacks/trustacks/pkg/api"
	_ "github.com/trustacks/trustacks/pkg/api"
)

var rpcServer *jsonrpc.RPCServer

// New creates a new server instance.
func New(host, port string) {
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("starting server on %s\n", addr)
	http.Handle("/rpc", rpcServer)
	http.ListenAndServe(addr, nil)
}

func init() {
	rpcServer = jsonrpc.NewServer()
	rpcServer.Register("v1", &api.APIV1Handler{})
	if os.Getenv("MODE") == "dev" {
		rpcServer.Register("test", &api.TestHandler{})
	}
}
