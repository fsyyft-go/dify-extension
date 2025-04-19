// Copyright 2025 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	kitlog "github.com/fsyyft-go/kit/log"

	appconf "github.com/fsyyft-go/dify-extension/internal/conf"
)

const (
	// 健康检查接口的路径点。
	pointPing = "ping"
	// 食物选择的类型标识。
	choiceFoods = "foods"
	// 饮品选择的类型标识。
	choiceDrinks = "drinks"

	// 健康检查接口的响应结果。
	resultPing = "pong"
	// 食物列表的预设响应结果，包含多种水果名称，以逗号分隔。
	resultFoods = "apple,banana,orange,grape,watermelon,kiwi,peach,pear,pineapple,lemon,mango,blueberry,raspberry,blackberry,grapefruit,apricot,avocado,coconut,fig,guava,lychee,olive,papaya,passion fruit,pomegranate,star fruit,dragon fruit,plum"
	// 饮品列表的预设响应结果，包含多种饮品名称，以逗号分隔。
	resultDrinks = "coffee,tea,juice,water,milk,beer,wine,whiskey,rum,vodka,gin,brandy,tequila,coke,pepsi,sprite,7up,fanta,red bull,monster"
)

var (
	apikey = "f513e140-b10d-4075-98ad-4f75b41d7c3e"
)

type (
	// requestInput 定义了请求输入的数据结构。
	requestInput struct {
		// 选择类型，可以是 foods 或 drinks。
		Choice string `json:"choice"`
	}

	// requestParams 定义了请求参数的数据结构。
	requestParams struct {
		// 应用程序标识符。
		AppID string `json:"app_id"`
		// 工具变量名称。
		ToolVariable string `json:"tool_variable"`
		// 输入参数。
		Inputs requestInput `json:"inputs"`
		// 查询字符串。
		Query string `json:"query"`
	}

	// requestData 定义了完整的请求数据结构。
	requestData struct {
		// 请求的路径点。
		Point string `json:"point"`
		// 请求的参数集合。
		Params requestParams `json:"params"`
	}

	// responseData 定义了响应数据的结构。
	responseData struct {
		// 响应结果字符串。
		Result string `json:"result"`
	}

	// server 定义了服务器的核心结构。
	server struct {
		// 结构化日志记录器。
		logger kitlog.Logger
		// 服务器配置对象。
		conf *appconf.Config
	}
)

// New 创建并初始化一个新的服务器实例。
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

// ServeHTTP 实现了 http.Handler 接口，处理所有 HTTP 请求。
// 根据请求路径将请求分发到相应的处理函数。
//
// 参数：
//   - w：HTTP 响应写入器
//   - r：HTTP 请求对象
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	// 初始化请求和响应结构体。
	req := requestData{}
	resp := responseData{}

	// 获取基础日志记录器。
	l := s.logger
	// 获取 Authorization 头部值，可能包含多个值。
	authHeaders := r.Header["Authorization"]
	authorized := false
	bearer := fmt.Sprintf("Bearer %s", apikey)

	// 读取请求体内容。
	body, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("读取请求体失败", "error", err)
		goto END
	}
	// 重新设置请求体，供后续使用。
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// 遍历所有 Authorization 值，检查是否有匹配的 apikey。
	for _, auth := range authHeaders {
		if auth == bearer {
			authorized = true
			break
		}
	}
	if !authorized {
		err = fmt.Errorf("unauthorized")
		l.Error("未授权")
		goto END
	}

	// 解析请求体 JSON 数据。
	if errDecoder := json.NewDecoder(r.Body).Decode(&req); errDecoder != nil {
		l.Error("解析请求体 JSON 失败", "error", errDecoder)
		err = errDecoder
		goto END
	}

	// 添加请求相关的日志字段。
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

	// 根据请求点和选择类型处理请求。
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

	// 记录响应结果。
	l.Info(resp.Result)

END:
	// 处理错误情况并返回响应。
	if err != nil {
		resp.Result = err.Error()
	}
	// 将响应编码为 JSON 并写入响应流。
	if errEncode := json.NewEncoder(w).Encode(resp); nil != errEncode {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
