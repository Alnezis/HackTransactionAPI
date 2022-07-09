package merchant

import (
	"HackTransactionAPI/api"
	"HackTransactionAPI/app"
)

func Exist(merchantName string) bool {
	var exist bool
	err := app.DB.Get(&exist, `SELECT exists(select id FROM transaction where  merchant_name = $1);`, merchantName)
	api.CheckErrInfo(err, "Merchant Exist")
	return exist
}
