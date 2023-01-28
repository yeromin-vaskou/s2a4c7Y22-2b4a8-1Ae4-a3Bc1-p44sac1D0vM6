package main

import (
	// "context"
	// "encoding/json"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/thedevsaddam/gojsonq/v2"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://localhost:27017/"
const url = "https://api.covid19api.com/summary"

var client *http.Client

type General struct {
	ID        string      `json:"ID"`
	Message   string      `json:"Message"`
	Global    Global      `json:"Global"`
	Countries []Countries `json:"Countries"`
	Date      time.Time   `json:"Date"`
}
type Global struct {
	NewConfirmed   int       `json:"NewConfirmed"`
	TotalConfirmed int       `json:"TotalConfirmed"`
	NewDeaths      int       `json:"NewDeaths"`
	TotalDeaths    int       `json:"TotalDeaths"`
	NewRecovered   int       `json:"NewRecovered"`
	TotalRecovered int       `json:"TotalRecovered"`
	Date           time.Time `json:"Date"`
}
type Premium struct {
}
type Countries struct {
	ID             string    `json:"ID"`
	Country        string    `json:"Country"`
	CountryCode    string    `json:"CountryCode"`
	Slug           string    `json:"Slug"`
	NewConfirmed   int       `json:"NewConfirmed"`
	TotalConfirmed int       `json:"TotalConfirmed"`
	NewDeaths      int       `json:"NewDeaths"`
	TotalDeaths    int       `json:"TotalDeaths"`
	NewRecovered   int       `json:"NewRecovered"`
	TotalRecovered int       `json:"TotalRecovered"`
	Date           time.Time `json:"Date"`
	Premium        Premium   `json:"Premium"`
}

func GetDecode(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func GetJson(data interface{}) string {
	val, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		fmt.Println("err: ", err)
	}
	return string(val)
}

func GetRecoveredAmount(url string) {
	var General General
	err := GetDecode(url, &General)
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Printf("Here is how many people recovered from Covid: %d!\n", General.Global.TotalRecovered)
	}
}

func GetCountryByCode(obj interface{}, code string) {
	var res string = GetJson(&obj)
	jq := gojsonq.New().FromString(res)

	cool_json := jq.From("Countries").Where("CountryCode", "=", code).Get()

	b, err := json.MarshalIndent(cool_json, "", "  ")
	if err != nil {
		fmt.Println("marshal err: ", err)
	}

	fmt.Println(string(b))
}

func main() {
	client = &http.Client{Timeout: 10 * time.Second}
	var general General

	GetDecode(url, &general)
	GetCountryByCode(&general, "DZ")
}
