package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func main() {
	r := gin.Default()
	client = &http.Client{Timeout: 10 * time.Second}
	mongo_client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongo_client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	coll := mongo_client.Database("pro-ect").Collection("covid")
	var general General
	GetDecode(url, &general)
	doc := &general

	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted: %v\n", result.InsertedID)

	r.GET("/country/:code", func(c *gin.Context) {
		var res string = GetJson(&general)
		jq := gojsonq.New().FromString(res)
		var code = c.Params.ByName("code")

		cool_json := jq.From("Countries").Where("CountryCode", "=", code).Get()

		b, err := json.MarshalIndent(cool_json, "", "   ")
		if err != nil {
			fmt.Println("marshal err: ", err)
		}

		c.JSON(200, string(b))
	})

	r.Run()
}
