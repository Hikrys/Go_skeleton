//go:build wireinject
// +build wireinject

package main

import (
	"Go_skeleton/internal/biz"
	"Go_skeleton/internal/data"
	"Go_skeleton/internal/service"
	"Go_skeleton/pkg/app"
	"Go_skeleton/pkg/config"
	"Go_skeleton/pkg/server"
	"github.com/google/wire"
)

func initApp(cfg *config.Config) (*AppContainer, func(), error) {
	panic(wire.Build(
		//配置包拆解
		config.ProviderSet,

		// 基础设施
		server.ProviderSet,

		// 业务逻辑
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,

		// 最终组装
		app.NewApp,

		//通知Wire  创建一个 AppContainer 指针，里面的字段（App, UserService）请自动填充 (*)"
		wire.Struct(new(AppContainer), "*"),
	))
}
