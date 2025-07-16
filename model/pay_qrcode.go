package model

import "github.com/suifengpiao14/sqlbuilder"

/*
CREATE TABLE `pay_qrcode` (
  `id` bigint(20) NOT NULL,
  `pay_url` varchar(255) DEFAULT NULL,
  `price` double NOT NULL,
  `type` int(11) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*/

type PayQRCodeModel struct {
	Id     int64  `json:"id" gorm:"column:id"`
	PayUrl string `json:"pay_url" gorm:"column:pay_url"`
	Price  int    `json:"price" gorm:"column:price"`
	Type   int    `json:"type" gorm:"column:type"`
}

var table_pay_qrcode = sqlbuilder.NewTableConfig("pay_qrcode").WithHandler(DBHander).AddColumns(
	sqlbuilder.NewColumn("id", sqlbuilder.GetField(NewId)),
	sqlbuilder.NewColumn("pay_url", sqlbuilder.GetField(NewPayUrl)),
	sqlbuilder.NewColumn("price", sqlbuilder.GetField(NewPrice)),
	sqlbuilder.NewColumn("type", sqlbuilder.GetField(NewPayingAgent)),
)

type PayQRCodeRepository struct {
	repository sqlbuilder.Repository[PayQRCodeModel]
}

func NewPayQRCodeRepository() PayQRCodeRepository {
	return PayQRCodeRepository{
		repository: sqlbuilder.NewRepository[PayQRCodeModel](table_pay_qrcode),
	}
}

func (s PayQRCodeRepository) GetByPrice(realPrice int, payingAgent string) (payQRCodeModelRef *PayQRCodeModel, exists bool, err error) {
	fs := sqlbuilder.Fields{
		NewPrice(realPrice).AppendWhereFn(sqlbuilder.ValueFnForward),
		NewPayingAgent(payingAgent).AppendWhereFn(sqlbuilder.ValueFnForward),
	}
	payQRCodeModel, exists, err := s.repository.First(fs)
	if err != nil {
		return nil, false, err
	}

	return &payQRCodeModel, exists, nil
}
