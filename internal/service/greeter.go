// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package service

import (
	"context"
	"fmt"

	kit_log "github.com/fsyyft-go/kit/log"

	app_helloworld_v1 "github.com/fsyyft-go/dify-extension/api/helloworld/v1"
	app_biz "github.com/fsyyft-go/dify-extension/internal/biz"
	app_conf "github.com/fsyyft-go/dify-extension/internal/conf"
)

type (
	// greeterService 实现了 GreeterHTTPServer 接口，提供问候服务。
	greeterService struct {
		// logger 用于服务日志记录。
		logger kit_log.Logger
		// conf 存储服务配置信息。
		conf *app_conf.Config
		// uc 用于处理问候相关的业务逻辑。
		uc app_biz.GreeterUsecase
	}
)

// NewGreeterService 创建一个新的 GreeterHTTPServer 服务实例。
//
// 参数：
//   - logger：日志记录器，用于服务日志记录。
//   - conf：服务配置信息。
//   - uc：问候用例的业务逻辑实现。
//
// 返回：
//   - app_helloworld_v1.GreeterHTTPServer：问候服务的实现实例。
func NewGreeterService(logger kit_log.Logger, conf *app_conf.Config, uc app_biz.GreeterUsecase) app_helloworld_v1.GreeterHTTPServer {
	return &greeterService{
		logger: logger,
		conf:   conf,
		uc:     uc,
	}
}

// SayHello 发送问候消息。
//
// 参数：
//   - ctx：上下文信息。
//   - in：包含问候请求的参数。
//
// 返回：
//   - *app_helloworld_v1.HelloReply：问候响应。
//   - error：可能发生的错误。
func (s *greeterService) SayHello(ctx context.Context, in *app_helloworld_v1.HelloRequest) (*app_helloworld_v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &app_biz.Greeter{
		Hello: in.Name,
	})
	if nil != err {
		return nil, err
	}
	return &app_helloworld_v1.HelloReply{
		Message: fmt.Sprintf("Hello %s", g.Hello),
	}, nil
}
