package server

import "github.com/google/wire"

// ProviderSet 把 HTTP 和 RPC Server 打包
var ProviderSet = wire.NewSet(NewRestServer, NewRpcServer)
