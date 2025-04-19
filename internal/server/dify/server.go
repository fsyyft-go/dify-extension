// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package dify

import (
	"encoding/json"
	"net/http"

	kitlog "github.com/fsyyft-go/kit/log"

	appconf "github.com/fsyyft-go/dify-extension/internal/conf"
)

type (
	requestData struct {
		Point string `json:"point"`
	}
	responseData struct {
		Result string `json:"result"`
	}

	server struct {
		logger kitlog.Logger   // 结构化日志记录器。
		conf   *appconf.Config // 服务器配置对象。
	}
)

// New 创建并初始化一个新的 Passkey 服务器实例。
//
// 参数：
//   - logger：结构化日志记录器，用于服务器运行时的日志记录
//   - conf：服务器配置对象，包含所有必要的配置参数
//
// 返回：
//   - http.Handler：配置完成的 HTTP 请求处理器
func New(logger kitlog.Logger, conf *appconf.Config) http.Handler {
	h := &server{
		logger: logger,
		conf:   conf,
	}

	return h
}

// ServeHTTP 实现了 http.Handler 接口，处理所有 WebAuthn 相关的 HTTP 请求。
// 根据请求路径将请求分发到相应的处理函数。
//
// 参数：
//   - w：HTTP 响应写入器
//   - r：HTTP 请求对象
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	req := requestData{}
	resp := responseData{}

	if errDecoder := json.NewDecoder(r.Body).Decode(&req); errDecoder != nil {
		s.logger.Error("failed to decode request body", "error", errDecoder)
		err = errDecoder
		goto END
	}
	switch req.Point {
	case "ping":
		resp.Result = "pong"
	}
END:
	if err != nil {
		resp.Result = err.Error()

	}
	if errEncode := json.NewEncoder(w).Encode(resp); nil != errEncode {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
