package personalqrcodepayment_test

import (
	"testing"

	"github.com/suifengpiao14/personalqrcodepayment"
	"github.com/suifengpiao14/sqlbuilder"
)

func init() {
	sqlbuilder.CreateTableIfNotExists = true
}

func TestCrate(t *testing.T) {
	cfg := personalqrcodepayment.Config{}
	payOrderService := personalqrcodepayment.NewPayOrderService(cfg)
	payOrderIn := personalqrcodepayment.PayOrderCreateIn{
		OrderId:     "pId2654981",
		PayingAgent: personalqrcodepayment.PayingAgent_Wechat,
		OrderAmount: 10000,
		Sign:        "54845",
		Param:       `{"userId":"1234","userName":"hah","prodcutId":"157485"}`,
	}
	payOrderService.Create(payOrderIn)

}
