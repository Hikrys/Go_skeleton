package data

import (
	"github.com/google/wire"
)

// 把 NewUserRepo 暴露出去
var ProviderSet = wire.NewSet(NewUserRepo)
