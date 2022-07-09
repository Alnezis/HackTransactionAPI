package transaction

import (
	"HackTransactionAPI/api"
	"HackTransactionAPI/app"
	"fmt"
	"github.com/CossackPyra/pyraconv"
	"github.com/tealeg/xlsx"
)

func Get(id string) Transaction {
	var i Transaction
	err := app.DB.Get(&i, `select * from transaction where id=$1`, id)
	api.CheckErrInfo(err, "GetT")
	return i
}

func New(i Transaction) int {
	var id int
	err := app.DB.Get(&id, `INSERT INTO transaction (user_id, check_id, product_name, product_cost, merchant_name, mcc) VALUES ($1,$2,$3,$4,$5,$6) returning id`,
		i.UserId, i.CheckId, i.ProductName, i.ProductCost, i.MerchantName, i.MCC)
	api.CheckErrInfo(err, "NewT")
	return id
}

type Transaction struct {
	ID      int `json:"id" db:"id"`
	UserId  int `json:"user_id" db:"user_id"`
	CheckId int `json:"check_id" db:"check_id"`

	ProductName string `json:"product_name" db:"product_name"`
	ProductCost int    `json:"product_cost" db:"product_cost"`

	MerchantName string `json:"merchant_name" db:"merchant_name"`
	MCC          int    `json:"mcc" db:"mcc"`
}

func GetTr(userID string, offset string) []Transaction {
	rows, err := app.DB.Queryx("SELECT * FROM transaction where user_id=$1 ORDER BY id ASC offset $2", userID, offset)
	api.CheckErrInfo(err, "GetTr 1")

	i := []Transaction{}

	for rows.Next() {
		var item Transaction
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "GetTr")
		i = append(i, item)
	}
	_ = rows.Close()
	return i
}

type Check struct {
	CheckId      int    `json:"check_id" db:"check_id"`
	MerchantName string `json:"merchant_name" db:"merchant_name"`
	Count        int    `json:"count" db:"count"`
	Sum          int    `json:"sum" db:"sum"`
	Mcc          int    `json:"mcc" db:"mcc"`
}

func GetChecks(userID string) []Check {
	rows, err := app.DB.Queryx("SELECT check_id, merchant_name, count(check_id), sum(product_cost), mcc FROM transaction where user_id=$1 GROUP BY check_id, merchant_name, mcc;", userID)
	api.CheckErrInfo(err, "GetChecks")

	i := []Check{}

	for rows.Next() {
		var item Check
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "GetChecks")
		i = append(i, item)
	}
	_ = rows.Close()
	return i
}

type Product struct {
	ProductName string `json:"product_name" db:"product_name"`
	ProductCost int    `json:"product_cost" db:"product_cost"`
}

func GetCheckContent(userID, checkID string) []Product {
	rows, err := app.DB.Queryx("select product_name, product_cost from transaction where user_id=$1 and check_id=$2;", userID, checkID)
	api.CheckErrInfo(err, "GetCheckContent")

	i := []Product{}

	for rows.Next() {
		var item Product
		err = rows.StructScan(&item)
		api.CheckErrInfo(err, "GetCheckContent")
		i = append(i, item)
	}
	_ = rows.Close()
	return i
}

func ParseFILE() {
	excelFileName := "./HistoryDataSet.csv.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println("err")
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			var t Transaction
			fmt.Println(row.Cells)
			t.UserId = int(pyraconv.ToFloat64(row.Cells[0]))
			t.CheckId = int(pyraconv.ToFloat64(row.Cells[1]))
			t.ProductName = pyraconv.ToString(row.Cells[2])

			t.ProductCost = int(pyraconv.ToFloat64(row.Cells[3]))

			t.MerchantName = pyraconv.ToString(row.Cells[4])

			t.MCC = int(pyraconv.ToFloat64(row.Cells[5]))

			fmt.Println(New(t))
		}
	}
}
