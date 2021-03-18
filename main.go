package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type rintisRequest struct {
	Pan                  string `json:"pan"`
	ProcessingCode       string `json:"processingCode"`
	TotalAmount          int    `json:"totalAmount"`
	TransmissionDateTime string `json:"transmissionDateTime"`
	Stan                 string `json:"stan"`
	LocalTransactionTime string `json:"localTransactionTime"`
	LocalTransactionDate string `json:"localTransactionDate"`
	CaptureDate          string `json:"captureDate"`
	AcquirerID           string `json:"acquirerId"`
	Track2Data           string `json:"track2Data"`
	Refnum               string `json:"refnum"`
	TerminalID           string `json:"terminalId"`
	CardAcceptorData     string `json:"cardAcceptorData"`
	AdditionalData       string `json:"additionalData"`
	Currency             string `json:"currency"`
	PIN                  string `json:"personalIdentificationNumber"`
	TerminalData         string `json:"terminalData"`
	AccountTo            string `json:"accountTo"`
	TokenData            string `json:"tokenData"`
}

func main() {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	i := 0
	for i < 10 {
		wg.Add(1)

		go func(wg *sync.WaitGroup, num int) {
			nice := rand.Intn(3000)
			time.Sleep(time.Duration(nice) * time.Millisecond)
			data := generateDummy(strconv.Itoa(nice))
			fmt.Printf("wait %dms for msg no %d \n", nice, num)
			request(data)
			wg.Done()
		}(&wg, i)
		i++
	}

	wg.Wait()
	fmt.Println("Done testing!")

}

func generateDummy(num string) rintisRequest {
	data := rintisRequest{
		Pan:                  num + "4900000008883",
		ProcessingCode:       "253010",
		TotalAmount:          11550000,
		TransmissionDateTime: "0921082022",
		Stan:                 num,
		LocalTransactionTime: "082022",
		LocalTransactionDate: "0921",
		CaptureDate:          "0921",
		AcquirerID:           "12345678901",
		Track2Data:           "1234567890123456789012345678901234567",
		Refnum:               "678615554461",
		Currency:             "123",
		CardAcceptorData:     "PI04Q001CD30SUSAENMC03UMI",
		AdditionalData:       "PI04Q001CD30SUSAENMC03UMI",
		PIN:                  "1234567890123456",
		TerminalData:         "123456789012345",
		TerminalID:           "1234567890123456",
		AccountTo:            "1234567890123456789",
		TokenData:            "1234567890123456789012345678901234567890",
	}
	return data
}

func request(num interface{}) {
	client := &http.Client{}
	bodyReq, _ := json.Marshal(num)
	req, err := http.NewRequest("POST", "http://localhost:6010/epay/rintis", bytes.NewBuffer(bodyReq))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatalf("Failed to sent request to https://tiruan.herokuapp.com/biller. Error: %v\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response from https://tiruan.herokuapp.com/biller. Error: %v\n", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("%s \n", body)
}
