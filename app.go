package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/data", HandlerGetData)
	http.ListenAndServe(":8080", router)
}

func HandlerGetData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	keyword := r.FormValue("keyword")
	safeText := url.QueryEscape(keyword)
	url := fmt.Sprintf("https://api.bukalapak.com/v2/products.json?keywords=%s&page=1&per_page=24", safeText)
	//Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	var record ReturnAPI
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(record); err != nil {
		panic(err)
	}
}

type Task struct {
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

type Tasks []Task

type ReturnAPI struct {
	Product []Product `json:"products"`
}
type Product struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	City           string `json:"city"`
	Price          int    `json:"price"`
	Category       string `json:"category"`
	SellerUsername string `json:"seller_username"`
	SellerName     string `json:"seller_name"`
	Province       string `json:"province"`
	Url            string `json:"url"`
	Weight         int    `json:"weight"`
	Stock          int    `json:"stock"`
}
