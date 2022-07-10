package statistic

import (
	"HackTransactionAPI/api"
	"HackTransactionAPI/app"
	"HackTransactionAPI/etc"
	"fmt"
	"math"
	"sort"
)

type T struct {
	ProductName string  `json:"product_name" db:"product_name"`
	Sum         float64 `json:"sum" db:"sum"`
}

type Result struct {
	ProductName string  `json:"product_name,omitempty" db:"product_name"`
	A           bool    `json:"a,omitempty" db:"a"`
	PsA         float64 `json:"psa" db:"psa"`
	B           bool    `json:"b,omitempty" db:"b"`
	PsB         float64 `json:"psb" db:"psb"`
}

func ProductRating() []Result {

	var m = map[string]*Result{}

	rows, err := app.DB.Queryx("SELECT product_name, sum(product_cost)::float / (select totalSumTr())*100 as sum FROM transaction GROUP BY product_name  ORDER BY sum DESC;")
	api.CheckErrInfo(err, "GetTr 1")

	var i float64

	for rows.Next() {
		var item T
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "GetChecks")
		if i < 80 {

			m[item.ProductName] = &Result{A: true, PsA: math.Round(item.Sum*100) / 100}
			i += item.Sum
		} else {

			m[item.ProductName] = &Result{A: false, PsA: item.Sum}
			continue
		}

	}
	fmt.Println(i)
	_ = rows.Close()

	rowsb, err := app.DB.Queryx("SELECT product_name, count(*)::float / (select totalCountTr())*100 as sum FROM transaction GROUP BY product_name  ORDER BY sum DESC;")
	api.CheckErrInfo(err, "GetTr 1")

	i = 0

	for rowsb.Next() {
		var item T
		err = rowsb.StructScan(&item)
		api.CheckErrInfo(err, "GetChecks")
		if i < 80 {
			m[item.ProductName].PsB = math.Round(item.Sum*100) / 100
			m[item.ProductName].B = true
			i += item.Sum
		} else {
			fmt.Println(i)
			break
			//continue
		}
	}
	_ = rowsb.Close()

	var res = []Result{}
	for pn, v := range m {
		if v.A {
			if v.B {
				//		fmt.Println(fmt.Sprintf("%s - A: %f, B: %f", pn, v.PsA, v.PsB))
				res = append(res, Result{
					ProductName: pn,
					PsA:         v.PsA,
					PsB:         v.PsB,
				})
			}
		}
	}

	sort.Slice(res, func(i, j int) (less bool) {
		return res[i].ProductName < res[j].ProductName
	})

	return res
}

func MerchantProductRating(merchantName string) []Result {

	var m = map[string]*Result{}

	var sumCost float64
	err := app.DB.Get(&sumCost, `SELECT sum(product_cost) as sum FROM transaction where  merchant_name = $1;`, merchantName)
	api.CheckErrInfo(err, "MerchantProductRating 1")

	rows, err := app.DB.Queryx(`SELECT product_name, sum(product_cost)::float/$1*100 as sum FROM transaction
                                                            where merchant_name = $2 GROUP BY product_name  ORDER BY sum DESC;`, sumCost, merchantName)
	api.CheckErrInfo(err, "MerchantProductRating 2")

	var i float64

	for rows.Next() {
		var item T
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "MerchantProductRating 3")
		if i < 80 {

			m[item.ProductName] = &Result{A: true, PsA: math.Round(item.Sum*100) / 100}
			i += item.Sum
		} else {

			m[item.ProductName] = &Result{A: false, PsA: item.Sum}
			continue
		}

	}
	fmt.Println(i)
	_ = rows.Close()

	var sumCount float64
	err = app.DB.Get(&sumCount, `SELECT count(*) as sum FROM transaction where  merchant_name = $1;`, merchantName)
	api.CheckErrInfo(err, "MerchantProductRating 4")

	rows, err = app.DB.Queryx(`SELECT product_name, count(product_cost)::float/$1*100 as sum FROM transaction
                                                            where merchant_name = $2 GROUP BY product_name  ORDER BY sum DESC;`, sumCount, merchantName)
	api.CheckErrInfo(err, "MerchantProductRating 5")

	i = 0
	for rows.Next() {
		var item T
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "MerchantProductRating 7")
		if i < 80 {
			m[item.ProductName].PsB = math.Round(item.Sum*100) / 100
			m[item.ProductName].B = true
			i += item.Sum
		} else {
			fmt.Println(i)
			break
			//continue
		}
	}
	_ = rows.Close()

	var res = []Result{}
	for pn, v := range m {
		if v.A {
			if v.B {
				//		fmt.Println(fmt.Sprintf("%s - A: %f, B: %f", pn, v.PsA, v.PsB))
				res = append(res, Result{
					ProductName: pn,
					PsA:         v.PsA,
					PsB:         v.PsB,
				})
			}
		}
	}
	sort.Slice(res, func(i, j int) (less bool) {
		return res[i].ProductName < res[j].ProductName
	})
	fmt.Println(res, len(res))
	return res
}

type Summary struct {
	MerchantName  string  `json:"merchant_name" db:"merchant_name"`
	Sum           float64 `json:"sum" db:"sum"`
	Count         int     `json:"count" db:"count"`
	Turnover      float64 `json:"turnover" db:"turnover"`
	Users         int     `json:"users" db:"users"`
	TurnoverUsers float64 `json:"turnover_users" db:"turnover_users"`

	InterchangeSum  float64 `json:"interchange_sum" db:"interchange_sum"`
	TurnoverProduct float64 `json:"turnover_product" db:"turnover_product"`

	Mcc     int    `json:"mcc" db:"mcc"`
	MccName string `json:"mcc_name" db:"mcc_name"`
}

func SummaryByMerchant(merchantName string) Summary {
	var i Summary

	err := app.DB.Get(&i, `select
    merchant_name,
    sum(product_cost) as sum,
    count(distinct check_id_global) as count,
    round(sum(product_cost)/count(distinct check_id_global)::numeric,2) as turnover,
    count(distinct user_id) as users,
    round(sum(product_cost)/count(distinct user_id)::numeric,2) as turnover_users,
    round(sum(interchange_sum)::numeric,2) as interchange_sum,
    round(sum(interchange_sum)/ sum(product_cost)*100::numeric,2) as turnover_product,
    mcc
from transaction where merchant_name=$1
group by merchant_name, mcc order by count desc;`, merchantName)
	api.CheckErrInfo(err, "SummaryByMerchant 1")
	i.MccName = etc.MccName(i.Mcc)
	return i
}

func SummaryMerchantAll() []Summary {
	rows, err := app.DB.Queryx(`select
    merchant_name,
    sum(product_cost) as sum,
    count(distinct check_id_global) as count,
    round(sum(product_cost)/count(distinct check_id_global)::numeric,2) as turnover,
    count(distinct user_id) as users,
    round(sum(product_cost)/count(distinct user_id)::numeric,2) as turnover_users,
    round(sum(interchange_sum)::numeric,2) as interchange_sum,
    round(sum(interchange_sum)/ sum(product_cost)*100::numeric,2) as turnover_product,
    mcc
from transaction
group by merchant_name, mcc order by turnover_product desc;`)
	api.CheckErrInfo(err, "SummaryMerchantAll 1")

	i := []Summary{}

	for rows.Next() {
		var item Summary
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "SummaryMerchantAll 2")
		item.MccName = etc.MccName(item.Mcc)
		i = append(i, item)
	}
	_ = rows.Close()
	return i
}
