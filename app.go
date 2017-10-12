package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/data", HandlerGetData)
	router.GET("/getpost", HandlerAllData)
	router.GET("/post", HandlerOneData)
	http.ListenAndServe(":8081", router)
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

func HandlerAllData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	posts, count := GetPostData()

	allposts := Posts{
		AllPosts: posts,
	}

	meta := ResultMeta{
		Count: count,
	}

	result := Result{
		Data: allposts,
		Meta: meta,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func HandlerOneData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	param := r.FormValue("id")
	id := url.QueryEscape(param)
	posts, count := GetPostWhereData(id)

	allposts := Posts{
		AllPosts: posts,
	}

	meta := ResultMeta{
		Count: count,
	}

	result := Result{
		Data: allposts,
		Meta: meta,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
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

type Post struct {
	IDPost      int     `json:"idpost"`
	Description string  `json:"description"`
	Latitude    float64 `json:"lat"`
	Langitude   float64 `json:"lang"`
	Image       string  `json:"image"`
}

type Posts struct {
	AllPosts []Post `json:"posts"`
}

type ResultMeta struct {
	Count int `json:"total_data"`
}

type Result struct {
	Data Posts      `json:"data"`
	Meta ResultMeta `json:"meta"`
}

const (
	DbUsername = "b665ecd09d7fe2"
	DbPassword = "a801aadc"
	DbName     = "kongko"
	DbHost     = "tcp(ap-cdbr-azure-southeast-b.cloudapp.net:3306)"
)

func OpenConnection() *sql.DB {
	conn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8", DbUsername, DbPassword, DbHost, DbName)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Println("Cannot Open Connection MySQL Database ", err)
		return nil
	} else {
		log.Println("Connected")
	}
	return db
}

func CloseConnection(db *sql.DB) {
	db.Close()
}

func GetPostData() ([]Post, int) {
	db := OpenConnection()
	count := 0
	defer CloseConnection(db)

	rows, err := db.Query("SELECT * FROM post")

	var posts []Post

	if err != nil {
		return posts, count
	}

	for rows.Next() {
		var post_id int
		var post_description string
		var post_lat float64
		var post_lang float64
		var post_image string

		err = rows.Scan(&post_id, &post_description, &post_image, &post_lang, &post_lat)
		if err != nil {
			log.Println("Cannot Query MySQL Database ", err)
		}
		posts = append(posts, Post{
			IDPost:      post_id,
			Description: post_description,
			Image:       post_image,
			Langitude:   post_lang,
			Latitude:    post_lat,
		})
		count++
	}

	return posts, count

}

func GetPostWhereData(id string) ([]Post, int) {
	db := OpenConnection()
	count := 0
	defer CloseConnection(db)
	q := ""
	if id != "" {
		q = fmt.Sprintf("SELECT * FROM post WHERE idpost=%s", id)
	} else {
		q = "SELECT * FROM post"
	}
	rows, err := db.Query(q)

	var posts []Post

	if err != nil {
		return posts, count
	}

	for rows.Next() {
		var post_id int
		var post_description string
		var post_lat float64
		var post_lang float64
		var post_image string

		err = rows.Scan(&post_id, &post_description, &post_image, &post_lang, &post_lat)
		if err != nil {
			log.Println("Cannot Query MySQL Database ", err)
		}
		posts = append(posts, Post{
			IDPost:      post_id,
			Description: post_description,
			Image:       post_image,
			Langitude:   post_lang,
			Latitude:    post_lat,
		})
		count++
	}

	return posts, count

}
