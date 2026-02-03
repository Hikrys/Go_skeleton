package biz

import "github.com/google/wire"

// 假设我们有个 UserUsecase
var ProviderSet = wire.NewSet(NewUserUsecase)
