// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package data 实现了应用程序的数据访问层，提供数据存储、检索和管理的功能。
package data

import (
	"github.com/google/wire"

	kit_log "github.com/fsyyft-go/kit/log"

	app_conf "github.com/fsyyft-go/dify-extension/internal/conf"
)

var (
	// ProviderSet 是数据层依赖注入的提供者集合。
	ProviderSet = wire.NewSet(
		NewData,
		NewGreeterRepo,
	)
)

type (
	// Data 定义了数据层的接口。
	// 该接口提供了数据层的公共依赖和基础操作能力。
	Data interface{}

	// data 实现了 Data 接口。
	data struct {
		// logger 用于记录日志信息。
		logger kit_log.Logger
		// conf 存储应用配置信息。
		conf *app_conf.Config
	}
)

// NewData 创建一个新的 Data 实例。
//
// 参数：
//   - logger：日志记录器，用于记录操作日志。
//   - conf：应用配置信息。
//
// 返回：
//   - Data：数据层接口实现。
//   - func()：清理函数，用于资源释放。
//   - error：可能的错误信息。
func NewData(logger kit_log.Logger, conf *app_conf.Config) (Data, func(), error) {
	cleanup := func() {
		logger.Info("closing the data resources")
	}
	return &data{
		logger: logger,
		conf:   conf,
	}, cleanup, nil
}
