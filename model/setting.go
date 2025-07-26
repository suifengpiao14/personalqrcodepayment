package model

import (
	"github.com/suifengpiao14/sqlbuilder"
	"gitlab.huishoubao.com/gopackage/keyvalue"
)

/*
CREATE TABLE `setting` (
  `vkey` varchar(255) NOT NULL,
  `vvalue` varchar(255) DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
*/
/*
 INSERT INTO `setting` (`vkey`, `vvalue`) VALUES
('user', 'admin'),
('pass', 'admin'),
('notifyUrl', ''),
('returnUrl', ''),
('key', ''),
('lastheart', '0'),
('lastpay', '0'),
('jkstate', '-1'),
('close', '5'),
('payQf', '1'),
('wxpay', ''),
('zfbpay', '');
*/

var DBHander = sqlbuilder.NewGormHandler(sqlbuilder.GormDB)

func getTableSetting() sqlbuilder.TableConfig {
	var table_setting = sqlbuilder.NewTableConfig("setting").WithHandler(DBHander).AddColumns(
		sqlbuilder.NewColumn("vkey", sqlbuilder.GetField(keyvalue.NewKeyField)),
		sqlbuilder.NewColumn("vvalue", sqlbuilder.GetField(keyvalue.NewValueField)),
	)
	return table_setting
}

type SettingSerivce struct {
	keyvalue.KeyValueService
}

func NewSettingService() *SettingSerivce {
	dbColumnRefer := keyvalue.KeyValueDbColumnRefer{
		Key:   "vkey",
		Value: "vvalue",
	}
	kvservice := keyvalue.NewKeyValueService(getTableSetting(), dbColumnRefer)
	service := &SettingSerivce{kvservice}
	return service
}

func (s SettingSerivce) GetSignKey() (signKey string, err error) {
	value, err := s.KeyValueService.Get("key", nil)
	if err != nil {
		return "", err
	}
	signKey = value.String()
	return signKey, nil
}

func (s SettingSerivce) GetOrderExpire() (closeTime int, err error) {
	value, err := s.KeyValueService.Get("close", nil)
	if err != nil {
		return 0, err
	}
	closeTime = value.Int()
	return closeTime, nil
}

func init() {
	initSetting()
}

func initSetting() {
	records := keyvalue.KeyValueModels{
		{Key: "user", Value: "admin"},
		{Key: "pass", Value: "admin"},
		{Key: "notifyUrl", Value: ""},
		{Key: "returnUrl", Value: ""},
		{Key: "key", Value: ""},
		{Key: "lastheart", Value: "0"},
		{Key: "lastpay", Value: "0"},
		{Key: "jkstate", Value: "-1"},
		{Key: "close", Value: "5"},
		{Key: "payQf", Value: "1"},
		{Key: "wxpay", Value: ""},
		{Key: "zfbpay", Value: ""},
	}
	keyvalue.RegisterInitKeyValues(records...)

}
