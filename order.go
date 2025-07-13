package personalqrcodepayment

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type CreateOrderIn struct {
	PayID  string  `json:"payId" form:"payId"`
	Type   int     `json:"type" form:"type"`
	Price  float64 `json:"price" form:"price"`
	Sign   string  `json:"sign" form:"sign"`
	IsHtml int     `json:"isHtml" form:"isHtml"`
	Param  string  `json:"param" form:"param"`
}

type CreateOrderOut struct {
}

type Config struct {
	Key string `json:"key" form:"key"`
}

// CreateOrder 创建订单
func CreateOrder(in CreateOrderIn, cfg Config) (out *CreateOrderOut, err error) {
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	// 验证签名
	if in.Sign != Signature(in.PayID, in.Param, in.Type, in.Price, cfg.Key) {
		return nil, errors.New("签名验证失败")
	}
	return
}

// Validate 验证请求参数
func (req *CreateOrderIn) Validate() error {
	// 验证payId
	if req.PayID == "" {
		return errors.New("请传入商户订单号")
	}

	// 验证type
	if req.Type == 0 {
		return errors.New("请传入支付方式=>1|微信 2|支付宝")
	}

	if req.Type != 1 && req.Type != 2 {
		return errors.New("支付方式错误=>1|微信 2|支付宝")
	}

	// 验证price
	if req.Price <= 0 {
		return errors.New("订单金额必须大于0")
	}

	// 验证sign
	if req.Sign == "" {
		return errors.New("请传入签名")
	}

	// 设置默认值
	if req.IsHtml == 0 {
		req.IsHtml = 0
	}

	if req.Param == "" {
		req.Param = ""
	}

	return nil
}

// Signature 生成签名
func Signature(payID, param string, typ int, price float64, key string) string {
	// 将参数按顺序拼接
	data := payID + param + strconv.Itoa(typ) + strconv.FormatFloat(price, 'f', 2, 64) + key

	// 计算MD5哈希
	hash := md5.Sum([]byte(data))

	// 转换为十六进制字符串
	return hex.EncodeToString(hash[:])
}

// OrderIDGenerator 生成订单ID（格式：YYYYMMDDHHMMSS + 4位随机数）
func OrderIDGenerator() string {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成时间部分（格式：YYYYMMDDHHMMSS）
	timePart := time.Now().Format("20060102150405")

	// 生成4位随机数（1-9之间的数字）
	randPart := fmt.Sprintf("%d%d%d%d",
		rand.Intn(9)+1,
		rand.Intn(9)+1,
		rand.Intn(9)+1,
		rand.Intn(9)+1)

	// 组合订单ID
	return timePart + randPart
}
