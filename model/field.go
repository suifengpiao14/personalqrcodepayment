package model

import (
	"github.com/suifengpiao14/commonlanguage"
	"github.com/suifengpiao14/sqlbuilder"
)

/*
CREATE TABLE `pay_order` (
  `id` bigint(20) NOT NULL,
  `close_date` bigint(20) NOT NULL,
  `create_date` bigint(20) NOT NULL,
  `is_auto` int(11) NOT NULL,
  `notify_url` varchar(255) DEFAULT NULL,
  `order_id` varchar(255) DEFAULT NULL,
  `param` varchar(255) DEFAULT NULL,
  `pay_date` bigint(20) NOT NULL,
  `pay_id` varchar(255) DEFAULT NULL,
  `pay_url` varchar(255) DEFAULT NULL,
  `price` double NOT NULL,
  `really_price` double NOT NULL,
  `return_url` varchar(255) DEFAULT NULL,
  `state` int(11) NOT NULL,
  `type` int(11) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*/

func NewId(id int) *sqlbuilder.Field {
	return commonlanguage.NewId(id)
}

func NewClosedAt(closedAt string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(closedAt, "closedAt", "关单时间", 0)
}

func NewExpire(expire int) *sqlbuilder.Field {
	return sqlbuilder.NewIntField(expire, "expire", "超时时间，单位秒", 0)
}

var NewCreatedAt = commonlanguage.NewCreatedAt

func NewAnyAmount(isAuto int) *sqlbuilder.Field {
	return sqlbuilder.NewIntField(isAuto, "isAuto", "是否自定义金额1-是,0-否", 0)
}

func NewNotifyUrl(notifyUrl string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(notifyUrl, "notifyUrl", "支付成功回调地址", 0)
}

func NewOrderId(orderId string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(orderId, "orderId", "订单号", 0)
}

func NewParam(param string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(param, "param", "支付参数", 0)
}

func NewPaidAt(paidAt string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(paidAt, "paidAt", "支付时间", 0)
}

func NewPayId[T string | []string](payId T) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(payId, "payId", "支付流水号", 0)
}

func NewPayUrl(payUrl string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(payUrl, "payUrl", "支付链接地址", 0)
}
func NewPrice(price int) *sqlbuilder.Field {
	return sqlbuilder.NewField(price).SetName("price").SetTitle("金额")
}
func NewPriceStr(price string) *sqlbuilder.Field {
	return sqlbuilder.NewField(price).SetName("price").SetTitle("金额")
}
func NewPaidPrice(paidPrice int) *sqlbuilder.Field {
	return sqlbuilder.NewField(paidPrice).SetName("paidPrice").SetTitle("实际支付金额")
}
func NewReturnUrl(returnUrl string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(returnUrl, "returnUrl", "支付完成跳转地址", 0)
}

var NewDeletedAt = commonlanguage.NewDeletedAt

func NewState(state string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(state, "state", "支付状态", 0)
}

const (
	Pay_order_type_wechat = 1
	Pay_order_type_alipay = 2
)

func NewPayingAgent(type_ string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(type_, "type", "支付类型1-微信，2-支付宝", 0)
}
