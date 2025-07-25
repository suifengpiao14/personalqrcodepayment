package personalqrcodepayment

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/personalqrcodepayment/model"
)

type PayOrderService struct {
	config Config
}

func NewPayOrderService(config Config) PayOrderService {
	return PayOrderService{
		config: config,
	}

}

type PayOrderCreateIn struct {
	PayId       string `json:"payId"`
	PayingAgent string `json:"payingAgent"` // 支付机构 weixin:微信 alipay:支付宝
	OrderAmount int    `json:"orderPrice"`  // 订单金额，单位分
	Sign        string `json:"sign"`
	Param       string `json:"param"`
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

// Create 创建订单
func (s PayOrderService) Create(in PayOrderCreateIn) (out *PayOrder, err error) {
	err = in.Validate()
	if err != nil {
		return nil, err
	}
	cfg := s.config
	// 验证签名
	if in.Sign != Signature(in.PayId, in.Param, in.PayingAgent, in.OrderAmount, cfg.Key) {
		return nil, errors.New("签名验证失败")
	}
	tmpPriceSerivice := model.NewTmpPriceSerivce()
	realPrice := in.OrderAmount
	orderId := OrderIDGenerator()
	isFindTmpPrice := false
	settingService := model.NewSettingService()
	for range 10 {
		tmpPrice := fmt.Sprintf("%d-%s", realPrice, in.PayingAgent)
		err = tmpPriceSerivice.InsertIgnore(tmpPrice, orderId)
		if err == nil {
			isFindTmpPrice = true
			break
		}
		if !errors.Is(err, model.Err_TmpPriceAlreadyExist) {
			return nil, err
		}
		err = nil
		if cfg.PayQf == 1 {
			realPrice++
		} else {
			realPrice--
		}
	}
	if !isFindTmpPrice {
		err = errors.New("订单超出负荷，请稍后重试")
		return nil, err
	}

	payUrl := cfg.PayUrl
	isAnyAmount := true
	payQRCodeService := model.NewPayQRCodeRepository()
	payQRCodeModel, exists, err := payQRCodeService.GetByPrice(realPrice, in.PayingAgent)
	if err != nil {
		return nil, err
	}
	if exists {
		isAnyAmount = false
		payUrl = payQRCodeModel.PayUrl
	}

	payOrderService := model.NewPayOrderRepository()
	_, exists, err = payOrderService.GetByPayId(in.PayId)
	if err != nil {
		return nil, err
	}
	if exists {
		err = errors.New("订单已存在")
		return nil, err
	}

	createdAt := time.Now().Format(time.DateTime)
	expire, err := settingService.GetOrderExpire()
	if err != nil {
		return nil, err
	}
	payOrderIn := model.PayOrderAddIn{
		Expire:      expire,
		CreatedAt:   createdAt,
		AnyAmount:   cast.ToInt(isAnyAmount),
		NotifyUrl:   cfg.NotifyUrl,
		OrderId:     orderId,
		Param:       in.Param,
		PayId:       in.PayId,
		PayUrl:      payUrl,
		Price:       in.OrderAmount,
		PaidPrice:   realPrice,
		ReturnUrl:   cfg.ReturnUrl,
		State:       PayOrderModel_state_pending.String(),
		PayingAgent: in.PayingAgent,
	}
	err = payOrderService.Add(payOrderIn)
	if err != nil {
		return nil, err
	}

	out = &PayOrder{
		PayId:       payOrderIn.PayId,
		OrderId:     payOrderIn.OrderId,
		OrderPrice:  payOrderIn.Price,
		PaidPrice:   payOrderIn.PaidPrice,
		PayUrl:      payOrderIn.PayUrl,
		AnyAmount:   payOrderIn.AnyAmount,
		State:       PayOrderState(payOrderIn.State),
		Expire:      payOrderIn.Expire,
		CreatedAt:   payOrderIn.CreatedAt,
		PayingAgent: payOrderIn.PayingAgent,
	}
	return out, nil
}

// Validate 验证请求参数
func (req *PayOrderCreateIn) Validate() error {
	// 验证payId
	if req.PayId == "" {
		return errors.New("请传入商户订单号")
	}
	payingAgent := cast.ToInt(req.PayingAgent)
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

// GetOrderPayInfo 获取订单支付信息
func (s PayOrderService) GetOrderPayInfo(orderId string) (payOrders PayOrders, err error) {
	r := model.NewPayOrderRepository()
	models, err := r.GetByOrderId(orderId)
	if err != nil {
		return nil, err
	}
	for _, v := range models {
		payOrder := PayOrder{
			PayId:      v.PayId,
			OrderId:    v.OrderId,
			OrderPrice: v.Price,
			PaidPrice:  v.PaidPrice,
			PayUrl:     v.PayUrl,
			AnyAmount:  v.AnyAmount,
			State:      PayOrderState(v.State),
			Expire:     v.Expire,
		}
		payOrders = append(payOrders, payOrder)
	}
	return payOrders, nil
}

// Pay 支付订单
func (s PayOrderService) Pay(payId string) (err error) {
	r := model.NewPayOrderRepository()
	model, err := r.GetByPayIdMust(payId)
	if err != nil {
		return err
	}
	stateFSM := NewPayOrderStateMachine(PayOrderState(model.State))
	err = stateFSM.CanPay()
	if err != nil {
		return err
	}
	err = r.Pay(model.PayId, PayOrderModel_state_paid.String(), model.State)
	if err != nil {
		return err
	}
	return nil
}

func (s PayOrderService) IsPaid(orderId string) (ok bool, err error) {
	records, err := s.GetOrderPayInfo(orderId)
	if err != nil {
		return false, err
	}
	payFinished := records.IsPayFinished()
	return payFinished, nil
}

// CloseByOrderId 关闭订单支付，当订单关闭时，关闭订单对应的支付单
func (s PayOrderService) CloseByOrderId(orderId string) (err error) {
	records, err := s.GetOrderPayInfo(orderId)
	if err != nil {
		return err
	}

	r := model.NewPayOrderRepository()
	closeBatchIn := make([]model.CloseIn, 0)
	for _, v := range records {
		closeIn := model.CloseIn{PayId: v.PayId, NewState: string(PayOrderModel_state_closed), OldState: v.State.String()}
		closeBatchIn = append(closeBatchIn, closeIn)
	}

	err = r.CloseBatch(closeBatchIn...)
	if err != nil {
		return err
	}
	return nil
}

func (s PayOrderService) GetByPayId(payId string) (payOrder *PayOrder, err error) {
	r := model.NewPayOrderRepository()
	model, err := r.GetByPayIdMust(payId)
	if err != nil {
		return nil, err
	}
	out := &PayOrder{
		PayId:       model.PayId,
		OrderId:     model.OrderId,
		OrderPrice:  model.Price,
		PaidPrice:   model.PaidPrice,
		PayUrl:      model.PayUrl,
		AnyAmount:   model.AnyAmount,
		State:       PayOrderState(model.State),
		Expire:      model.Expire,
		CreatedAt:   model.CreatedAt,
		PayingAgent: model.PayingAgent,
	}
	return out, nil
}

func (s PayOrderService) CloseByPayId(payId string) (err error) {
	record, err := s.GetByPayId(payId)
	if err != nil {
		return err
	}
	stateFSM := record.GetStateFSM()
	err = stateFSM.CanPay()
	if err != nil {
		return err
	}

	r := model.NewPayOrderRepository()
	err = r.CloseByPayId(record.PayId, PayOrderModel_state_paid.String(), record.State.String())
	if err != nil {
		return err
	}
	return nil
}
func (s PayOrderService) ExpiredByPayId(payId string) (err error) {
	record, err := s.GetByPayId(payId)
	if err != nil {
		return err
	}
	stateFSM := record.GetStateFSM()
	err = stateFSM.CanExpire()
	if err != nil {
		return err
	}

	r := model.NewPayOrderRepository()
	err = r.CloseByPayId(record.PayId, PayOrderModel_state_expired.String(), record.State.String())
	if err != nil {
		return err
	}
	return nil
}

func (s PayOrderService) Failed(payId string) (err error) {
	record, err := s.GetByPayId(payId)
	if err != nil {
		return err
	}
	stateFSM := record.GetStateFSM()
	err = stateFSM.CanPay()
	if err != nil {
		return err
	}

	r := model.NewPayOrderRepository()
	err = r.Failed(record.PayId, PayOrderModel_state_failed.String(), record.State.String())
	if err != nil {
		return err
	}
	return nil
}
