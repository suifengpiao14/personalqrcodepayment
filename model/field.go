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

func NewCloseDate(closeDate int) *sqlbuilder.Field {
	return sqlbuilder.NewIntField(closeDate, "closeDate", "关单时间", 0)
}

func NewCreateDate(createDate int) *sqlbuilder.Field {
	return sqlbuilder.NewIntField(createDate, "createDate", "创建时间", 0)
}
func NewIsAuto(isAuto int) *sqlbuilder.Field {
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

func NewPayDate(payDate int) *sqlbuilder.Field {
	return sqlbuilder.NewIntField(payDate, "payDate", "支付时间", 0)
}

func NewPayId(payId string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(payId, "payId", "支付流水号", 0)
}
func NewPayUrl(payUrl string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(payUrl, "payUrl", "支付链接地址", 0)
}
func NewPrice(price float64) *sqlbuilder.Field {
	return sqlbuilder.NewField(price).SetName("price").SetTitle("金额")
}
func NewReallyPrice(reallyPrice float64) *sqlbuilder.Field {
	return sqlbuilder.NewField(reallyPrice).SetName("reallyPrice").SetTitle("实际金额")
}
func NewReturnUrl(returnUrl string) *sqlbuilder.Field {
	return sqlbuilder.NewStringField(returnUrl, "returnUrl", "支付完成跳转地址", 0)
}

const (
	Pay_order_state_success = 1
	Pay_order_state_fail    = 0
)

func NewState(state int) *sqlbuilder.Field {
	return sqlbuilder.NewIntField(state, "state", "支付状态1-成功，0-失败", 0)
}

const (
	Pay_order_type_wechat = 1
	Pay_order_type_alipay = 2
)

func NewType(type_ int) *sqlbuilder.Field {
	return sqlbuilder.NewIntField(type_, "type", "支付类型1-微信，2-支付宝", 0)
}
