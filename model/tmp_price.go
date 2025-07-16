package model

import (
	"github.com/pkg/errors"
	"github.com/suifengpiao14/sqlbuilder"
)

/*
CREATE TABLE `tmp_price` (

	`price` varchar(255) NOT NULL,
	`oid` varchar(255) NOT NULL

) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*/
type TmpPriceModel struct {
	Price     string `json:"price" gorm:"column:price"`
	Oid       string `json:"oid" gorm:"column:oid"`
	DeletedAt string `json:"deletedAt" gorm:"column:deleted_at"`
}

func getTableTmpPrice() sqlbuilder.TableConfig {
	var table_setting = sqlbuilder.NewTableConfig("tmp_price").WithHandler(DBHander).AddColumns(
		sqlbuilder.NewColumn("price", sqlbuilder.GetField(NewPriceStr)),
		sqlbuilder.NewColumn("oid", sqlbuilder.GetField(NewOrderId)),
		sqlbuilder.NewColumn("deleted_at", NewDeletedAt()),
	)
	return table_setting
}

type TmpPriceService struct {
	//keyvalue.KeyValueService
	repository sqlbuilder.Repository[TmpPriceModel]
}

func NewTmpPriceSerivce() *TmpPriceService {
	// dbColumnRefer := keyvalue.KeyValueDbColumnRefer{
	// 	Key:   "vkey",
	// 	Value: "vvalue",
	// }
	tableConfig := getTableTmpPrice()
	//kvService := keyvalue.NewKeyValueService(getTableTmpPrice(), dbColumnRefer)
	service := &TmpPriceService{
		//KeyValueService: kvService,
		repository: sqlbuilder.NewRepository[TmpPriceModel](tableConfig),
	}
	return service
}

var Err_TmpPriceAlreadyExist = errors.New("tmp_price already exist")

func (s TmpPriceService) InsertIgnore(price string, orderId string) (err error) {
	fs := sqlbuilder.Fields{
		NewPriceStr(price),
		NewOrderId(orderId),
	}
	_, rowsAffected, err := s.repository.InsertWithLastId(fs, func(p *sqlbuilder.InsertParam) {
		p.WithInsertIgnore(true)
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return Err_TmpPriceAlreadyExist
	}
	return nil
}

func (s TmpPriceService) Deleted(price string) (err error) {
	fs := sqlbuilder.Fields{
		NewPriceStr(price),
		NewDeletedAt(),
	}
	err = s.repository.Delete(fs)
	if err != nil {
		return err
	}
	return nil
}
