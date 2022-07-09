package statistic

import (
	"HackTransactionAPI/api"
	"HackTransactionAPI/app"
	"fmt"
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

			m[item.ProductName] = &Result{A: true, PsA: item.Sum}
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
			m[item.ProductName].PsB = item.Sum
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
	return res
}
