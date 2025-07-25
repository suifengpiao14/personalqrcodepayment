package model

import (
	"time"

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

type PayOrderModel struct {
	Id          int64  `json:"id" gorm:"column:id"`
	ClosedAt    string `json:"closedAt" gorm:"column:closed_at"`
	Expire      int    `json:"expire" gorm:"column:expire"`
	CreatedAt   string `json:"createdAt" gorm:"column:create_at"`
	AnyAmount   int    `json:"anyAmount" gorm:"column:any_amount"`
	NotifyUrl   string `json:"notifyUrl" gorm:"column:notify_url"`
	OrderId     string `json:"orderId" gorm:"column:order_id"`
	Param       string `json:"param" gorm:"column:param"`
	PaidAt      string `json:"paidAt" gorm:"column:paid_at"`
	PayId       string `json:"payId" gorm:"column:pay_id"`
	PayUrl      string `json:"payUrl" gorm:"column:pay_url"`
	Price       int    `json:"price" gorm:"column:price"`
	PaidPrice   int    `json:"paidPrice" gorm:"column:paid_price"`
	ReturnUrl   string `json:"returnUrl" gorm:"column:return_url"`
	State       string `json:"state" gorm:"column:state"`
	PayingAgent string `json:"payingAgent" gorm:"column:paying_agent"`
}

type PayOrderModels []PayOrderModel

func (ms PayOrderModels) TotalPrice() int {
	var total = 0
	for _, m := range ms {
		total += m.Price
	}
	return total
}

var table_pay_order = sqlbuilder.NewTableConfig("pay_order").AddColumns(
	sqlbuilder.NewColumn("id", sqlbuilder.GetField(NewId)),
	sqlbuilder.NewColumn("close_at", sqlbuilder.GetField(NewClosedAt)),
	sqlbuilder.NewColumn("expire", sqlbuilder.GetField(NewExpire)),
	sqlbuilder.NewColumn("created_at", sqlbuilder.GetField(NewCreatedAt)),
	sqlbuilder.NewColumn("any_amount", sqlbuilder.GetField(NewAnyAmount)),
	sqlbuilder.NewColumn("notify_url", sqlbuilder.GetField(NewNotifyUrl)),
	sqlbuilder.NewColumn("order_id", sqlbuilder.GetField(NewOrderId)),
	sqlbuilder.NewColumn("param", sqlbuilder.GetField(NewParam)),
	sqlbuilder.NewColumn("paid_at", sqlbuilder.GetField(NewPaidAt)),
	sqlbuilder.NewColumn("pay_id", sqlbuilder.GetField(NewPayId[string])),
	sqlbuilder.NewColumn("pay_url", sqlbuilder.GetField(NewPayUrl)),
	sqlbuilder.NewColumn("price", sqlbuilder.GetField(NewPrice)),
	sqlbuilder.NewColumn("paid_price", sqlbuilder.GetField(NewPaidPrice)),
	sqlbuilder.NewColumn("return_url", sqlbuilder.GetField(NewReturnUrl)),
	sqlbuilder.NewColumn("state", sqlbuilder.GetField(NewState)),
	sqlbuilder.NewColumn("paying_agent", sqlbuilder.GetField(NewPayingAgent)),
).AddIndexs(sqlbuilder.Index{
	IsPrimary: true,
	ColumnNames: func(tableColumns sqlbuilder.ColumnConfigs) (columnNames []string) {
		return []string{"id"}
	},
},
	sqlbuilder.Index{
		Unique: true,
		ColumnNames: func(tableColumns sqlbuilder.ColumnConfigs) (columnNames []string) {
			return []string{"pay_id"}
		},
	},
	sqlbuilder.Index{
		ColumnNames: func(tableColumns sqlbuilder.ColumnConfigs) (columnNames []string) {
			return []string{"order_id"}
		},
	},
)

type PayOrderRepository struct {
	repository sqlbuilder.Repository[PayOrderModel]
}

func NewPayOrderRepository() PayOrderRepository {
	return PayOrderRepository{
		repository: sqlbuilder.NewRepository[PayOrderModel](table_pay_order),
	}
}

/*
   $createDate = time();
   $data = array(
       "close_date" => 0,
       "create_date" => $createDate,
       "is_auto" => $isAuto,
       "notify_url" => $notify_url,
       "order_id" => $orderId,
       "param" => $param,
       "pay_date" => 0,
       "pay_id" => $payId,
       "pay_url" => $payUrl,
       "price" => $price,
       "really_price" => $reallyPrice,
       "return_url" => $return_url,
       "state" => 0,
       "type" => $type

   );
*/

type PayOrderAddIn struct {
	Expire      int    `json:"timeOut"`
	CreatedAt   string `json:"createDate"`
	AnyAmount   int    `json:"isAuto"`
	NotifyUrl   string `json:"notifyUrl"`
	OrderId     string `json:"orderId"`
	Param       string `json:"param"`
	PayId       string `json:"payId"`
	PayUrl      string `json:"payUrl"`
	Price       int    `json:"price"`
	PaidPrice   int    `json:"paidPrice"`
	ReturnUrl   string `json:"returnUrl"`
	State       string `json:"state"`
	PayingAgent string `json:"payingAgent"`
}

func (in PayOrderAddIn) Fields() sqlbuilder.Fields {
	return sqlbuilder.Fields{
		NewExpire(in.Expire),
		NewCreatedAt(in.CreatedAt),
		NewAnyAmount(in.AnyAmount),
		NewNotifyUrl(in.NotifyUrl),
		NewOrderId(in.OrderId),
		NewParam(in.Param),
		NewPayId(in.PayId),
		NewPayUrl(in.PayUrl),
		NewPrice(in.Price),
		NewPaidPrice(in.PaidPrice),
		NewReturnUrl(in.ReturnUrl),
		NewState(in.State),
		NewPayingAgent(in.PayingAgent),
	}
}

func (po PayOrderRepository) Add(in PayOrderAddIn) (err error) {
	err = po.repository.Insert(in.Fields())
	if err != nil {
		return err
	}
	return nil
}

func (po PayOrderRepository) GetByPayId(payId string) (model PayOrderModel, exists bool, err error) {
	fs := sqlbuilder.Fields{
		NewPayId(payId).AppendWhereFn(sqlbuilder.ValueFnForward),
	}
	model, exists, err = po.repository.First(fs)
	if err != nil {
		return model, exists, err
	}
	return model, exists, nil
}

func (po PayOrderRepository) GetByPayIdMust(payId string) (model PayOrderModel, err error) {
	model, exists, err := po.GetByPayId(payId)
	if !exists {
		err = sqlbuilder.ERROR_NOT_FOUND
		return model, err
	}
	return model, nil
}

func (po PayOrderRepository) GetByOrderId(orderId string) (models PayOrderModels, err error) {
	fs := sqlbuilder.Fields{
		NewOrderId(orderId).AppendWhereFn(sqlbuilder.ValueFnForward),
	}
	models, err = po.repository.All(fs)
	if err != nil {
		return models, err
	}
	return models, nil
}

type ChangeStatusIn struct {
	PayId       string
	NewState    string
	OldState    string
	ExtraFields sqlbuilder.Fields
}

func (in ChangeStatusIn) Fields() sqlbuilder.Fields {
	fs := sqlbuilder.Fields{
		NewPayId(in.PayId).SetRequired(true).ShieldUpdate(true).AppendWhereFn(sqlbuilder.ValueFnForward),
		NewState(in.NewState).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			//查询条件值使用旧状态值
			f.WhereFns.ResetSetValueFn(func(inputValue any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) {
				return in.OldState, nil
			})
		}),
	}
	fs.Add(in.ExtraFields...)
	return fs
}

func (po PayOrderRepository) ChangeStatus(in ChangeStatusIn) (err error) {
	fs := in.Fields()
	err = po.repository.Update(fs)
	if err != nil {
		return err
	}
	return nil
}

func (po PayOrderRepository) Pay(payId string, paidState string, oldState string) (err error) {
	in := ChangeStatusIn{
		PayId:       payId,
		NewState:    paidState,
		OldState:    oldState,
		ExtraFields: sqlbuilder.Fields{NewPaidAt(time.Now().Format(time.DateTime))},
	}
	err = po.ChangeStatus(in)
	if err != nil {
		return err
	}
	return nil
}

func (po PayOrderRepository) CloseByPayId(payId string, closeState string, oldState string) (err error) {
	in := ChangeStatusIn{
		PayId:       payId,
		NewState:    closeState,
		OldState:    oldState,
		ExtraFields: sqlbuilder.Fields{NewClosedAt(time.Now().Format(time.DateTime))},
	}
	err = po.ChangeStatus(in)
	return err
}

type CloseIn = ChangeStatusIn

// CloseBatch 批量关闭订单支付状态，如果存在多个支付流水号，则全部关闭。当订单关闭时，使用事务批量关闭。
func (po PayOrderRepository) CloseBatch(closeInArr ...CloseIn) (err error) {

	extraFields := sqlbuilder.Fields{NewClosedAt(time.Now().Format(time.DateTime))}
	for i := range closeInArr {
		closeInArr[i].ExtraFields = extraFields
	}

	po.repository.Transaction(func(txRepository sqlbuilder.Repository[PayOrderModel]) (err error) {
		for _, closeIn := range closeInArr {
			fs := closeIn.Fields()
			err = txRepository.Update(fs)
			if err != nil {
				return err
			}
		}
		return nil

	})

	return nil
}

func (po PayOrderRepository) Failed(payId string, failedState string, oldState string) (err error) {
	in := ChangeStatusIn{PayId: payId, NewState: failedState, OldState: oldState}
	err = po.ChangeStatus(in)
	return err
}
