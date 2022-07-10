package user

import (
	"HackTransactionAPI/api"
	"HackTransactionAPI/app"
	"HackTransactionAPI/etc"
)

func Exist(merchantName string) bool {
	var exist bool
	err := app.DB.Get(&exist, `SELECT exists(select id FROM transaction where  user_id = $1);`, merchantName)
	api.CheckErrInfo(err, "User Exist")
	return exist
}

type Summary struct {
	MerchantName string  `json:"merchant_name,omitempty" db:"merchant_name,omitempty"`
	ChecksSum    float64 `json:"checks_sum" db:"checks_sum"`
	ChecksCount  int     `json:"checks_count" db:"checks_count"`
	MidCheck     float64 `json:"mid_check" db:"mid_check"`

	InterchangeSum         float64 `json:"interchange_sum" db:"interchange_sum"`
	TurnoverProductPercent float64 `json:"turnover_product_percent" db:"turnover_product_percent"`

	Mcc     int    `json:"mcc,omitempty" db:"mcc,omitempty"`
	MccName string `json:"mcc_name,omitempty" db:"mcc_name,omitempty"`
}

func SummaryUserInMerchantByName(userID, merchantName string) Summary {
	var i Summary
	err := app.DB.Get(&i, `select
    merchant_name,
    sum(product_cost) as checks_sum,
    count(distinct check_id_global) as checks_count,
    round(sum(product_cost)/count(distinct check_id_global)::numeric,2) as mid_check,
    round(sum(interchange_sum)::numeric,2) as interchange_sum,
    round(sum(interchange_sum)/ sum(product_cost)*100::numeric,2) as turnover_product_percent,
    mcc
from transaction where user_id=$1 and merchant_name=$2
group by merchant_name, mcc order by checks_sum desc;`, userID, merchantName)
	api.CheckErrInfo(err, "SummaryUserInMerchantByName 1")
	i.MccName = etc.MccName(i.Mcc)
	return i
}

func SummaryUserInMerchantAllSum(userID string) Summary {
	var i Summary

	err := app.DB.Get(&i, `select
    sum(product_cost) as checks_sum,
    count(distinct check_id_global) as checks_count,
    round(sum(product_cost)/count(distinct check_id_global)::numeric,2) as mid_check,
    round(sum(interchange_sum)::numeric,2) as interchange_sum,
    round(sum(interchange_sum)/ sum(product_cost)*100::numeric,2) as turnover_product_percent

from transaction where user_id=$1;`, userID)
	api.CheckErrInfo(err, "SummaryUserInMerchantByName 1")
	i.MccName = etc.MccName(i.Mcc)
	return i
}

func SummaryUserInMerchantAll(userID string) []Summary {
	rows, err := app.DB.Queryx(`select
    merchant_name,
    sum(product_cost) as checks_sum,
    count(distinct check_id_global) as checks_count,
    round(sum(product_cost)/count(distinct check_id_global)::numeric,2) as mid_check,
    round(sum(interchange_sum)::numeric,2) as interchange_sum,
    round(sum(interchange_sum)/ sum(product_cost)*100::numeric,2) as turnover_product_percent,
    mcc
from transaction where user_id=$1
group by merchant_name, mcc order by checks_sum desc;`, userID)
	api.CheckErrInfo(err, "SummaryUserInMerchantAll 1")

	i := []Summary{}

	for rows.Next() {
		var item Summary
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "SummaryUserInMerchantAll 2")
		item.MccName = etc.MccName(item.Mcc)
		i = append(i, item)
	}
	_ = rows.Close()
	return i
}
