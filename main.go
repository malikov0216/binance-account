package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var baseURL = "https://api.binance.com"
var endpoint = "/api/v3/account"

var (
	apiKey    = "pzIhpuLbh8zTs9xFNWwDHWpep9Tq3ubP3RuZ6yOhFUViGCQ51i4uG532Ge9vtiHH"
	secretKey = "lx6JAefTtUXbuEr31bRZb5bQ0owDQ6QgOuySzJE4UdNttsoAtDsvwmuKkISo7qzJ"
)

type Balances struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type AccountInfo struct {
	MakerCommission  int        `json:"makerCommission"`
	TakerCommission  int        `json:"takerCommission"`
	BuyerCommission  int        `json:"buyerCommission"`
	SellerCommission int        `json:"sellerCommission"`
	CanTrade         bool       `json:"canTrade"`
	CanWithdraw      bool       `json:"canWithdraw"`
	CanDeposit       bool       `json:"canDeposit"`
	UpdateTime       int64      `json:"updateTime"`
	AccountType      string     `json:"accountType"`
	Balances         []Balances `json:"balances"`
}

func main() {
	accountInfo := new(AccountInfo)
	fullURL := makeFullURL(baseURL, endpoint, secretKey)

	accountInfo.getBalance(fullURL, apiKey)
}

func encodeSecretKey(secretKey, query string) (string, error) {
	key := []byte(secretKey)
	message := query

	sig := hmac.New(sha256.New, key)
	_, err := sig.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig.Sum(nil)), nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeFullURL(baseURL, endpoint, secretKey string) string {
	url := baseURL + endpoint

	timestamp := makeTimestamp()
	queryTimestamp := fmt.Sprintf("timestamp=%d", timestamp)

	signature, err := encodeSecretKey(secretKey, queryTimestamp)
	if err != nil {
		log.Println(err)
	}

	return fmt.Sprintf("%s?%s&signature=%s", url, queryTimestamp, signature)
}

func (acc *AccountInfo) getBalance(fullURL, apiKey string) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", fullURL, nil)
	request.Header.Set("X-MBX-APIKEY", apiKey)
	request.URL.Query()
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(body, &acc)

	if len(acc.Balances) != 0 {
		for i := 0; i < len(acc.Balances); i++ {
			fmt.Printf("%s balance: %s", acc.Balances[i].Asset, acc.Balances[i].Free)
		}
	} else {
		fmt.Println("You have no balance")
	}
}
