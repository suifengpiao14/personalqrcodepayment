package personalqrcodepayment

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strconv"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/paymentrecord"
	paymentrecordrepository "github.com/suifengpiao14/paymentrecord/repository"
	"github.com/suifengpiao14/personalqrcodepayment/model"
	"github.com/suifengpiao14/sqlbuilder"
)

func MakeQRcode(content string) (png []byte, err error) {
	png, err = qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return png, nil
}

type PayOrderService struct {
	config Config
}

func NewPayOrderService(config Config) PayOrderService {
	return PayOrderService{
		config: config,
	}

}

const (
	PayingAgent_Wechat = "weixin"
	PayingAgent_Alipay = "alipay"
)

type PayOrder struct {
	PayId       string        `json:"payId"`
	OrderId     string        `json:"orderId"`
	OrderPrice  int           `json:"orderPrice"`
	PaidPrice   int           `json:"paidPrice"`
	PayUrl      string        `json:"payUrl"`
	AnyAmount   int           `json:"anyAmount"`
	State       PayOrderState `json:"state"`
	Expire      int           `json:"timeOut"`
	CreatedAt   string        `json:"date"`
	PayingAgent string        `json:"payingAgent"`
}

func (m *PayOrder) GetStateFSM() (stateMachine *PayOrderStateMachine) {
	stateMachine = NewPayOrderStateMachine(m.State)
	return stateMachine
}

type PayOrders []PayOrder

// PaidMoney 支付单中已支付金额总和，使用时要确保全部为一个订单的支付单，否则无意义
func (orders PayOrders) PaidMoney() (paidMoney int) {
	if len(orders) == 0 {
		return 0
	}
	firstOrderId := orders[0].OrderId
	for _, order := range orders {
		if order.OrderId != firstOrderId { // 确保只统计同一个支付单的金额
			err := errors.New("PayOrders.PaidMoney 方法只能用于同一个订单的支付单")
			panic(err)
		}
		if order.State == PayOrderModel_state_paid {
			paidMoney += cast.ToInt(order.PaidPrice)
		}
	}
	return paidMoney
}

func (orders PayOrders) IsPayFinished() (payfinished bool) {
	if len(orders) == 0 {
		return true
	}
	orderPrice := orders[0].OrderPrice
	payfinished = orderPrice >= orders.PaidMoney() // 所有支付单支付金额总和大于等于订单金额即为支付完成
	return payfinished
}

func (orders PayOrders) CanClose() (err error) {
	for _, order := range orders {
		stateMachine := NewPayOrderStateMachine(order.State)
		err = stateMachine.CanClose()
		if err != nil {
			return err
		}

	}
	return nil
}

type Config struct {
	Key       string `json:"key"`
	NotifyUrl string `json:"notifyUrl"`
	ReturnUrl string `json:"returnUrl"`
	PayQf     int    `json:"payQf"`
	PayUrl    string `json:"payUrl"`
}

type PayOrderCreateIn struct {
	PayId            string `json:"payId"`
	OrderId          string `json:"orderId"`
	RecipientAccount string `json:"recipientAccount"` // 收款人ID
	PayAgent         string `json:"payAgent"`         // 支付机构 weixin:微信 alipay:支付宝
	OrderAmount      int    `json:"orderPrice"`       // 订单金额，单位分
	UserId           string `json:"userId"`
	Sign             string `json:"sign"`
	Param            string `json:"param"`
	PaymentAccount   string `json:"paymentAccount"`
	PaymentName      string `json:"paymentName"`
}

func getPayRecordSerivce() *paymentrecord.PayRecordService {
	return paymentrecord.NewPayRecordService(model.DBHander)
}

// Create 创建订单
func (s PayOrderService) Create(in PayOrderCreateIn) (out *PayOrder, err error) {
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	cfg := s.config
	// 验证签名
	if in.Sign != Signature(in.OrderId, in.Param, in.PayAgent, in.OrderAmount, cfg.Key) {
		return nil, errors.New("签名验证失败")
	}
	settingService := model.NewSettingService()
	expire, err := settingService.GetOrderExpire()
	if err != nil {
		return nil, err
	}

	payQRCodeService := model.NewPayQRCodeRepository()
	payQRCodeModel, err := payQRCodeService.LockQRCodeByOrderId(in.OrderId, in.RecipientAccount, in.OrderAmount, in.PayAgent)
	if err != nil {
		return nil, err
	}
	payrecordService := getPayRecordSerivce()
	payRecordInForQRPay := paymentrecord.PayRecordCreateIn{
		OrderId:          in.OrderId,
		PayAgent:         payQRCodeModel.PayAgent,
		OrderAmount:      in.OrderAmount,
		PayAmount:        payQRCodeModel.Amount,
		Expire:           expire,
		PayId:            paymentrecord.PayIdGenerator(),
		PayParam:         "",
		UserId:           in.UserId,
		ClientIp:         "127.0.0.1",
		RecipientAccount: payQRCodeModel.RecipientAccount,
		RecipientName:    payQRCodeModel.RecipientName,
		PaymentAccount:   in.PaymentAccount,
		PaymentName:      in.PaymentName,
		PayUrl:           payQRCodeModel.PayUrl,
		NotifyUrl:        "",
		ReturnUrl:        "",
		Remark:           "固定金额二维码支付",
	}

	payRecordInForCoupon := paymentrecord.PayRecordCreateIn{
		OrderId:          in.OrderId,
		PayAgent:         paymentrecordrepository.PayingAgent_Coupon,
		OrderAmount:      in.OrderAmount,
		PayAmount:        in.OrderAmount - payQRCodeModel.Amount,
		PayId:            paymentrecord.PayIdGenerator(),
		PayParam:         "",
		UserId:           in.UserId,
		ClientIp:         "127.0.0.1",
		RecipientAccount: payQRCodeModel.RecipientAccount,
		RecipientName:    payQRCodeModel.RecipientName,
		PaymentAccount:   "compon",
		PaymentName:      "compon",
		PayUrl:           "http://coupon.pay.com",
		NotifyUrl:        "",
		ReturnUrl:        "",
		Remark:           "优惠券支付",
	}

	err = payrecordService.Create(payRecordInForQRPay, payRecordInForCoupon)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Validate 验证请求参数
func (req *PayOrderCreateIn) Validate() error {
	// 验证payId
	if req.OrderId == "" {
		return errors.New("请传入商户订单号")
	}
	payingAgent := cast.ToInt(req.PayAgent)
	// 验证type
	if payingAgent == 0 {
		return errors.New("请传入支付方式=>1|微信 2|支付宝")
	}

	if payingAgent != 1 && payingAgent != 2 {
		return errors.New("支付方式错误=>1|微信 2|支付宝")
	}

	// 验证price
	if req.OrderAmount <= 0 {
		return errors.New("订单金额必须大于0")
	}

	// 验证sign
	if req.Sign == "" {
		return errors.New("请传入签名")
	}
	return nil
}

// Signature 生成签名
func Signature(payID, param string, payingAgent string, price int, key string) string {
	// 将参数按顺序拼接
	data := payID + param + payingAgent + strconv.Itoa(price) + key

	// 计算MD5哈希
	hash := md5.Sum([]byte(data))

	// 转换为十六进制字符串
	return hex.EncodeToString(hash[:])
}

// GetOrderPayInfo 获取订单支付信息
func (s PayOrderService) GetOrderPayInfo(orderId string) (models paymentrecordrepository.PayRecordModels, err error) {
	payrecordService := getPayRecordSerivce()
	models, err = payrecordService.GetOrderPayInfo(orderId)
	if err != nil {
		return nil, err
	}
	return models, nil
}

// Pay 支付订单
func (s PayOrderService) Pay(payAmount int, payAgent string) (err error) {
	payrecordService := getPayRecordSerivce()
	fs := sqlbuilder.Fields{
		paymentrecordrepository.NewPayAmount(payAmount),
		paymentrecordrepository.NewPayAgent(payAgent),
		paymentrecordrepository.NewState(paymentrecordrepository.PayOrderModel_state_paid.String()).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.SetOrderFn(sqlbuilder.OrderFnDesc)
		}),
	}
	model, err := payrecordService.GetFirstPayRecordByConditon(fs)
	if err != nil {
		return err
	}

	payIn := paymentrecord.PayIn{
		PayId: model.PayId,
	}
	_, err = payrecordService.Pay(payIn)
	if err != nil {
		return err
	}
	//用户支付成功后，优惠券自动支付,保证整个订单支付完成
	couponFs := sqlbuilder.Fields{
		paymentrecordrepository.NewOrderId(model.OrderId),
		paymentrecordrepository.NewPayAgent(paymentrecordrepository.PayingAgent_Coupon),
		paymentrecordrepository.NewState(paymentrecordrepository.PayOrderModel_state_paid.String()).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.SetOrderFn(sqlbuilder.OrderFnDesc)
		}),
	}
	couponModel, err := payrecordService.GetFirstPayRecordByConditon(couponFs)
	if err != nil {
		return err
	}
	componPayIn := paymentrecord.PayIn{
		PayId: couponModel.PayId,
	}

	_, err = payrecordService.Pay(componPayIn)
	if err != nil {
		return err
	}
	return nil
}
