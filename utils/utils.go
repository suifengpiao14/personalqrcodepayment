package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GenerateOrderID 生成唯一订单号
func GenerateOrderID() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	randomStr := strconv.FormatInt(rand.Int63n(1000000), 10)
	return "P" + timestamp + randomStr
}

// GenerateSign 生成签名
func GenerateSign(params map[string]interface{}, key string) string {
	// 按字典序排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接参数
	var signStr string
	for _, k := range keys {
		if k == "sign" || params[k] == "" {
			continue
		}
		signStr += fmt.Sprintf("%s=%v&", k, params[k])
	}
	signStr += fmt.Sprintf("key=%s", key)

	// 计算MD5
	h := md5.New()
	io.WriteString(h, signStr)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// VerifySign 验证签名
func VerifySign(params map[string]interface{}, sign string, key string) bool {
	// 生成签名
	generatedSign := GenerateSign(params, key)

	// 比较签名
	return generatedSign == sign
}

// SendHTTPRequest 发送HTTP请求
func SendHTTPRequest(method string, url string, params map[string]interface{}) (string, error) {
	// 创建请求
	client := &http.Client{Timeout: 30 * time.Second}

	var req *http.Request
	var err error

	if method == "GET" {
		// 构建GET请求参数
		paramStr := ""
		for k, v := range params {
			paramStr += fmt.Sprintf("%s=%v&", k, v)
		}
		if paramStr != "" {
			paramStr = strings.TrimRight(paramStr, "&")
			url += "?" + paramStr
		}

		req, err = http.NewRequest("GET", url, nil)
	} else {
		// 构建POST请求参数
		jsonData, err := json.Marshal(params)
		if err != nil {
			return "", err
		}

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return "", err
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// WriteJSONResponse 写入JSON响应
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
