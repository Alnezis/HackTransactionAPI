package controllers

import (
	"HackTransactionAPI/merchant"
	"HackTransactionAPI/statistic"
	"HackTransactionAPI/transaction"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	Result interface{} `json:"result"`
	Error  *Error      `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Transactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := r.URL.Query()
	fmt.Println(params)
	//if err != nil {
	//	json.NewEncoder(w).Encode(&Response{
	//		Error: &Error{
	//			Message: err.Error(),
	//			Code:    9,
	//		},
	//	})
	//	return
	//}

	if params["offset"] == nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "offset отсутствует",
				Code:    201,
			},
		})
		return
	}

	if params["user_id"] == nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "user_id отсутствует",
				Code:    202,
			},
		})
		return
	}

	i := transaction.GetTr(params["user_id"][0], params["offset"][0])

	json.NewEncoder(w).Encode(&Response{
		Result: i,
	})
}

func Checks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := r.URL.Query()
	fmt.Println(params)

	if params["user_id"] == nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "user_id отсутствует",
				Code:    202,
			},
		})
		return
	}

	i := transaction.GetChecks(params["user_id"][0])

	json.NewEncoder(w).Encode(&Response{
		Result: i,
	})
}

func MerchantProductRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := r.URL.Query()
	fmt.Println(params)

	if params["merchant_name"] == nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "merchant_name отсутствует в запросе",
				Code:    202,
			},
		})
		return
	}

	if !merchant.Exist(params["merchant_name"][0]) {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: fmt.Sprintf("merchant_name (%s) отсутствует в базе", params["merchant_name"][0]),
				Code:    500,
			},
		})
		return
	}

	i := statistic.MerchantProductRating(params["merchant_name"][0])

	json.NewEncoder(w).Encode(&Response{
		Result: i,
	})
}

func Products(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := r.URL.Query()
	fmt.Println(params)

	if params["user_id"] == nil && params["user_id"][0] == "" {
		log.Println("202 \"user_id отсутствует\"")

		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "user_id отсутствует",
				Code:    202,
			},
		})
		return
	}
	if params["check_id"] == nil && params["check_id"][0] == "" {
		log.Println("202 \"check_id отсутствует\"")
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "check_id отсутствует",
				Code:    202,
			},
		})
		return
	}
	i := transaction.GetCheckContent(params["user_id"][0], params["check_id"][0])

	json.NewEncoder(w).Encode(&Response{
		Result: i,
	})
}

func ProductRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	i := statistic.ProductRating()

	json.NewEncoder(w).Encode(&Response{
		Result: i,
	})
}

func SummaryMerchantAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	i := statistic.SummaryMerchantAll()

	json.NewEncoder(w).Encode(&Response{
		Result: i,
	})
}

func SummaryByMerchant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	params := r.URL.Query()

	if params["merchant_name"] == nil {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: "merchant_name отсутствует в запросе",
				Code:    202,
			},
		})
		return
	}

	if !merchant.Exist(params["merchant_name"][0]) {
		json.NewEncoder(w).Encode(&Response{
			Error: &Error{
				Message: fmt.Sprintf("merchant_name (%s) отсутствует в базе", params["merchant_name"][0]),
				Code:    500,
			},
		})
		return
	}

	i := statistic.SummaryByMerchant(params["merchant_name"][0])

	json.NewEncoder(w).Encode(&Response{
		Result: i,
	})
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	//  Ensure our file does not exceed 5MB
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)

	file, handler, err := r.FormFile("image")

	// Capture any errors that may arise
	if err != nil {
		fmt.Fprintf(w, "Error getting the file")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded file name: %+v\n", handler.Filename)
	fmt.Printf("Uploaded file size %+v\n", handler.Size)
	fmt.Printf("File mime type %+v\n", handler.Header)

	// Get the file content type and access the file extension
	fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[1]

	// Create the temporary file name
	fileName := fmt.Sprintf("upload-*.%s", fileType)
	// Create a temporary file with a dir folder
	tempFile, err := ioutil.TempFile("images", fileName)

	if err != nil {
		fmt.Println(err)
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Successfully uploaded file")
}
