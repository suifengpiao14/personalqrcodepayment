package model

import (
	"github.com/suifengpiao14/sqlbuilder"
	"gitlab.huishoubao.com/gopackage/keyvalue"
)

/*
CREATE TABLE `tmp_price` (

	`price` varchar(255) NOT NULL,
	`oid` varchar(255) NOT NULL

) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*/
type TmpPriceModel struct {
	Price string `json:"price" gorm:"column:price"`
	Oid   string `json:"oid" gorm:"column:oid"`
}

func getTableTmpPrice() sqlbuilder.TableConfig {
	var table_setting = sqlbuilder.NewTableConfig("tmp_price").WithHandler(DBHander).AddColumns(
		sqlbuilder.NewColumn("price", sqlbuilder.GetField(keyvalue.NewKeyField)),
		sqlbuilder.NewColumn("oid", sqlbuilder.GetField(keyvalue.NewValueField)),
	)
	return table_setting
}

func NewTmpPriceSerivce() keyvalue.KeyValueService {
	dbColumnRefer := keyvalue.KeyValueDbColumnRefer{
		Key:   "vkey",
		Value: "vvalue",
	}
	service := keyvalue.NewKeyValueService(getTableTmpPrice(), dbColumnRefer)
	return service
}
