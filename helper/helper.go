package helper

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
	"os"

	"github.com/valyala/fasthttp"
	"github.com/ferluci/fast-realip"
)

var (
	/* 日志输出使用 */
	output  = log.New(os.Stdout, "", 0)

	/* 返回值的 content-type */
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)


/* 处理返回值，返回json */
func RespJson(ctx *fasthttp.RequestCtx, data *map[string]interface{}) {
	(*data)["msg"] = "success"
	respJson := map[string]interface{}{
		"code": 0,
		"appId": "",
		"signType": "plain",
		"encType": "plain",
		"success": true,
		"data": *data,
		"timestamp": time.Now().Unix(),
	}
	doJSONWrite(ctx, fasthttp.StatusOK, respJson)
}

func RespError(ctx *fasthttp.RequestCtx, code int, msg string) {
	log.Println("Error: ", code, msg)
	respJson := map[string]interface{}{
		"code": code,
		"appId": "",
		"signType": "plain",
		"encType": "plain",
		"success": false,
		"data": map[string]interface{}{"msg": msg},
		"timestamp": time.Now().Unix(),
	}
	doJSONWrite(ctx, fasthttp.StatusOK, respJson)
}

func doJSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	start := time.Now()
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		elapsed := time.Since(start)
		log.Printf("", elapsed, err.Error(), obj)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

/*
	接口验签，返回data数据
*/
func CheckSign(content []byte) (*map[string]interface{}, error) {
	fields := make(map[string]interface{})
	if err := json.Unmarshal(content, &fields); err != nil {
		return nil, err
	}

	var appId, version, signType, signData string
	var timestamp int64
	var data map[string]interface{}
	var ok bool

	// 检查参数
	if appId, ok = fields["appId"].(string); !ok {
		return nil, fmt.Errorf("need appid")
	}
	if version, ok = fields["version"].(string); !ok {
		return nil, fmt.Errorf("need version")
	}
	if signType, ok = fields["signType"].(string); !ok {
		return nil, fmt.Errorf("need sign_type")
	}
	if signData, ok = fields["signData"].(string); !ok {
		return nil, fmt.Errorf("need sign_data")
	}
	if _, ok = fields["timestamp"].(float64); !ok {
		return nil, fmt.Errorf("need timestamp")
	} else {
		timestamp = int64(fields["timestamp"].(float64)) // 返回整数
	}
	if data, ok = fields["data"].(map[string]interface{}); !ok {
		return nil, fmt.Errorf("need data")
	}

	// 获取 secret，用户密钥的签名串
	secret, ok := Settings.Api.SECRET_KEY[appId]
	if !ok {
		return nil, fmt.Errorf("wrong appId")
	}

	// 检查版本
	if version != "1" {
		return nil, fmt.Errorf("wrong version")
	}

	// 检查签名类型
	if signType != "SHA256" {
		return nil, fmt.Errorf("unknown signType")
	}

	// 生成参数的key，并排序
	keys := getMapKeys(fields)
	sort.Strings(*keys)
	//fmt.Println(*keys)

	// data 串，用于验签， map已按key排序
	dataStr, _ := json.Marshal(data)

	// 拼接验签串
	var signString = string("")
	for _, k := range *keys {
		if k == "signData" {
			continue
		}
		if k == "data" {
			signString += k + "=" + string(dataStr) + "&"
		} else if k == "timestamp" {
			signString += k + "=" + strconv.FormatInt(timestamp, 10) + "&"
		} else {
			signString += k + "=" + fields[k].(string) + "&"
		}
	}
	signString += "key=" + secret
	//fmt.Println(signString)

	h := sha256.New()
	h.Write([]byte(signString))
	sum := h.Sum(nil)
	sha256Str := fmt.Sprintf("%x", sum)
	signStr := base64.StdEncoding.EncodeToString([]byte(sha256Str))
	//fmt.Println(sha256Str)

	if signStr != signData {
		fmt.Println(signStr)
		fmt.Println(signData)
		return nil, fmt.Errorf("wrong signature")
	}

	return &data, nil
}


// 日志格式处理

// "github.com/AubSs/fasthttplogger"
func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

// Combined format:
// [<time>] <remote-addr> | <HTTP/http-version> | <method> <url> - <status> - <response-time us> | <user-agent>
// [2017/05/31 - 13:27:28] 127.0.0.1:54082 | HTTP/1.1 | GET /hello - 200 - 48.279µs | Paw/3.1.1 (Macintosh; OS X/10.12.5) GCDHTTPRequest
func Combined(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("[%v] %v (%v) | %s | %s %s - %v - %v | %s",
			end.Format("2006/01/02 - 15:04:05"),
			ctx.RemoteAddr(),
			realip.FromRequest(ctx),
			getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	})
}


// 返回 map 所有 key
func getMapKeys(m map[string]interface{}) *[]string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return &keys
}
