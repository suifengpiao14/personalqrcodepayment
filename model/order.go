package model

import "github.com/suifengpiao14/sqlbuilder"

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

type PayOrderModel struct {
	Id          int64   `json:"id" gorm:"column:id"`
	CloseDate   int64   `json:"close_date" gorm:"column:close_date"`
	CreateDate  int64   `json:"create_date" gorm:"column:create_date"`
	IsAuto      int     `json:"is_auto" gorm:"column:is_auto"`
	NotifyUrl   string  `json:"notify_url" gorm:"column:notify_url"`
	OrderId     string  `json:"order_id" gorm:"column:order_id"`
	Param       string  `json:"param" gorm:"column:param"`
	PayDate     int64   `json:"payDate" gorm:"column:pay_date"`
	PayId       string  `json:"payId" gorm:"column:pay_id"`
	PayUrl      string  `json:"pay_url" gorm:"column:pay_url"`
	Price       float64 `json:"price" gorm:"column:price"`
	ReallyPrice float64 `json:"really_price" gorm:"column:really_price"`
	ReturnUrl   string  `json:"return_url" gorm:"column:return_url"`
	State       int     `json:"state" gorm:"column:state"`
	Type        int     `json:"type" gorm:"column:type"`
}

var table_pay_order = sqlbuilder.NewTableConfig("pay_order").AddColumns(
	sqlbuilder.NewColumn("id", sqlbuilder.GetField(NewId)),
	sqlbuilder.NewColumn("close_date", sqlbuilder.GetField(NewCloseDate)),
	sqlbuilder.NewColumn("create_date", sqlbuilder.GetField(NewCreateDate)),
	sqlbuilder.NewColumn("is_auto", sqlbuilder.GetField(NewIsAuto)),
	sqlbuilder.NewColumn("notify_url", sqlbuilder.GetField(NewNotifyUrl)),
	sqlbuilder.NewColumn("order_id", sqlbuilder.GetField(NewOrderId)),
	sqlbuilder.NewColumn("param", sqlbuilder.GetField(NewParam)),
	sqlbuilder.NewColumn("pay_date", sqlbuilder.GetField(NewPayDate)),
	sqlbuilder.NewColumn("pay_id", sqlbuilder.GetField(NewPayId)),
	sqlbuilder.NewColumn("pay_url", sqlbuilder.GetField(NewPayUrl)),
	sqlbuilder.NewColumn("price", sqlbuilder.GetField(NewPrice)),
	sqlbuilder.NewColumn("really_price", sqlbuilder.GetField(NewReallyPrice)),
	sqlbuilder.NewColumn("return_url", sqlbuilder.GetField(NewReturnUrl)),
	sqlbuilder.NewColumn("state", sqlbuilder.GetField(NewState)),
	sqlbuilder.NewColumn("type", sqlbuilder.GetField(NewType)),
)

type PayOrderRepository struct {
	sqlbuilder.Repository[PayOrderModel]
}

func NewPayOrderRepository() PayOrderRepository {
	return PayOrderRepository{
		Repository: sqlbuilder.NewRepository[PayOrderModel](table_pay_order),
	}
}

type PayOrderAddIn struct{}

func (in PayOrderAddIn) Fields() sqlbuilder.Fields {
	return sqlbuilder.Fields{}
}

func (po PayOrderRepository) Add(in PayOrderAddIn) (err error) {
	err = po.RepositoryCommand.Insert(in.Fields())
	if err != nil {
		return err
	}
	return nil

}
