// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package dify

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	kitlog "github.com/fsyyft-go/kit/log"

	appconf "github.com/fsyyft-go/dify-extension/internal/conf"
)

const (
	pointPing    = "ping"
	choiceFoods  = "foods"
	choiceDrinks = "drinks"

	resultPing   = "pong"
	resultFoods  = "apple,banana,orange,grape,watermelon,kiwi,peach,pear,pineapple,lemon,mango,blueberry,raspberry,blackberry,grapefruit,apricot,avocado,coconut,fig,guava,lychee,olive,papaya,passion fruit,pomegranate,star fruit,dragon fruit,plum"
	resultDrinks = "coffee,tea,juice,water,milk,beer,wine,whiskey,rum,vodka,gin,brandy,tequila,coke,pepsi,sprite,7up,fanta,red bull,monster"
)

type (
	requestInput struct {
		Choice string `json:"choice"`
	}

	requestParams struct {
		AppID        string       `json:"app_id"`
		ToolVariable string       `json:"tool_variable"`
		Inputs       requestInput `json:"inputs"`
		Query        string       `json:"query"`
	}
	requestData struct {
		Point  string        `json:"point"`
		Params requestParams `json:"params"`
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

	l := s.logger

	// 读取请求体内容。
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("读取请求体失败", "error", err)
		goto END
	}
	// 重新设置请求体，供后续使用。
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if errDecoder := json.NewDecoder(r.Body).Decode(&req); errDecoder != nil {
		s.logger.Error("failed to decode request body", "error", errDecoder)
		err = errDecoder
		goto END
	}

	l = l.WithField("point", req.Point)
	if len(req.Params.AppID) > 0 {
		l = l.WithField("app_id", req.Params.AppID)
	}
	if len(req.Params.ToolVariable) > 0 {
		l = l.WithField("tool_variable", req.Params.ToolVariable)
	}
	if len(req.Params.Query) > 0 {
		l = l.WithField("query", req.Params.Query)
	}
	if len(req.Params.Inputs.Choice) > 0 {
		l = l.WithField("choice", req.Params.Inputs.Choice)
	}

	switch req.Point {
	case pointPing:
		resp.Result = resultPing
	default:
		switch req.Params.Inputs.Choice {
		case choiceFoods:
			resp.Result = resultFoods
		case choiceDrinks:
			resp.Result = resultDrinks
		}
	}

	l.Info(resp.Result)

END:
	if err != nil {
		resp.Result = err.Error()

	}
	if errEncode := json.NewEncoder(w).Encode(resp); nil != errEncode {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
