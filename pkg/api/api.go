package api

import (
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
)

type APIV1Handler struct{}

func init() {
	rpcServer := jsonrpc.NewServer()
	rpcServer.Register("apiv1", &APIV1Handler{})
	http.Handle("/rpc/v1", rpcServer)
}
