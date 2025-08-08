package model

import (
	"errors"

	"github.com/suifengpiao14/commonlanguage"
	paymentrecordrepository "github.com/suifengpiao14/paymentrecord/repository"
	"github.com/suifengpiao14/sqlbuilder"
)

/*
CREATE TABLE `pay_qrcode` (
  `id` bigint(20) NOT NULL,
  `pay_url` varchar(255) DEFAULT NULL,
  `price` double NOT NULL,
  `type` int(11) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*/

type PayQRCodeModel struct {
	Id               int64  `json:"id" gorm:"column:Fid"`
	RecipientAccount string `json:"recipientAccount" gorm:"column:Frecipient_account"`
	RecipientName    string `json:"recipientName" gorm:"column:Frecipient_name"`
	PayUrl           string `json:"pay_url" gorm:"column:Fpay_url"`
	PayAmount        int    `json:"payAmount" gorm:"column:Fpay_amount"`
	PayAgent         string `json:"payAgent" gorm:"column:Fpay_agent"`
	LockKey          string `json:"lockKey" gorm:"column:Flock_key"`
	CreatedAt        string `json:"createdAt" gorm:"column:Fcreated_at"`
	UpdatedAt        string `json:"updatedAt" gorm:"column:Fupdated_at"`
}

var table_pay_qrcode = sqlbuilder.NewTableConfig("pay_qrcode").WithHandler(DBHander).AddColumns(
	sqlbuilder.NewColumn("Fid", sqlbuilder.GetField(paymentrecordrepository.NewId)),
	sqlbuilder.NewColumn("Frecipient_account", sqlbuilder.GetField(paymentrecordrepository.NewRecipientAccount)),
	sqlbuilder.NewColumn("Frecipient_name", sqlbuilder.GetField(paymentrecordrepository.NewRecipientName)),
	sqlbuilder.NewColumn("Flock_key", sqlbuilder.GetField(NewLockKey)),
	sqlbuilder.NewColumn("Fpay_url", sqlbuilder.GetField(paymentrecordrepository.NewPayUrl)),
	sqlbuilder.NewColumn("Fpay_amount", sqlbuilder.GetField(paymentrecordrepository.NewPayAmount)),
	sqlbuilder.NewColumn("Fpay_agent", sqlbuilder.GetField(paymentrecordrepository.NewPayAgent)),
	sqlbuilder.NewColumn("Fcreated_at", sqlbuilder.GetField(paymentrecordrepository.NewCreatedAt)),
	sqlbuilder.NewColumn("Fupdated_at", sqlbuilder.GetField(paymentrecordrepository.NewUpdatedAt)),
)

type PayQRCodeRepository struct {
	repository sqlbuilder.Repository
}

func NewPayQRCodeRepository() PayQRCodeRepository {
	return PayQRCodeRepository{
		repository: sqlbuilder.NewRepository(table_pay_qrcode),
	}
}

func (s PayQRCodeRepository) LockQRCodeByOrderId(orderId string, recipientAccount string, amount int, payAgent string) (payQRCodeModelRef *PayQRCodeModel, err error) {
	fs := sqlbuilder.Fields{
		paymentrecordrepository.NewPayAgent(payAgent).AppendWhereFn(sqlbuilder.ValueFnForward),
		paymentrecordrepository.NewRecipientAccount(recipientAccount).AppendWhereFn(sqlbuilder.ValueFnForward),
		paymentrecordrepository.NewPayAmount(amount).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.WhereFns.Append(sqlbuilder.ValueFnBetween(nil, amount)) // 支付金额小于等于当前金额的二维码
			f.SetOrderFn(sqlbuilder.OrderFnDesc)                      // 倒序查询，确保锁定的二维码金额和实际金额最接近
		}).SetMinimum(1), // 最小金额为1分
		NewLockKey(orderId).Apply(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
			f.WhereFns.ResetSetValueFn(func(inputValue any, f *sqlbuilder.Field, fs ...*sqlbuilder.Field) (any, error) {
				return "", nil // 查询条件的值改成空字符串,即查找未锁定的记录，增加锁
			})
		}),
		commonlanguage.NewUpdateLimit(1), // 只锁定一条记录即可
	}
	err = s.repository.Update(fs)
	if err != nil {
		return nil, err
	}

	getFs := sqlbuilder.Fields{
		NewLockKey(orderId).AppendWhereFn(sqlbuilder.ValueFnForward),
		paymentrecordrepository.NewPayAgent(payAgent).AppendWhereFn(sqlbuilder.ValueFnForward),
	}
	payQRCodeModelRef = &PayQRCodeModel{}
	exists, err := s.repository.First(payQRCodeModelRef, getFs)
	if !exists {
		err = errors.New("符合条件的支付码全被占用，请稍后重试")
		return nil, err
	}
	return payQRCodeModelRef, nil
}
