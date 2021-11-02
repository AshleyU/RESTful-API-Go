package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Article struct {
	Id string `json:"Id"`
	Title string `json:"Title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

// Declaring global Articles array
// to populate the main function
// to simulate a database
var Articles []Article

func returnAllArticles(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    json.NewEncoder(w).Encode(Articles)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
    // gets the body of our POST request
    // returns the string response containing the request body    
    reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article 
    json.Unmarshal(reqBody, &article)
    // updates our global Articles array to include
    // our new Article
    Articles = append(Articles, article)

    json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
    // parses the path parameters
    vars := mux.Vars(r)
    // extracts the `id` of the article we to delete
    id := vars["id"]

    // loop through all articles
    for index, article := range Articles {
        // if our id path parameter matches one of our
        // articles
        if article.Id == id {
            // updates our Articles array to remove the 
            // article
            Articles = append(Articles[:index], Articles[index+1:]...)
        }
    }

	fmt.Println("Deleted")
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
    // parse the path parameters
    vars := mux.Vars(r)
    // extract the `id` of the article to update
    id := vars["id"]

	var updatedEvent Article 
	
	reqBody, _ := ioutil.ReadAll(r.Body) 
	json.Unmarshal(reqBody, &updatedEvent) 
	for i, article := range Articles { 
		
		if article.Id == id {

		article.Title = updatedEvent.Title
		article.Desc = updatedEvent.Desc
		article.Content = updatedEvent.Content
		Articles[i] = article
		json.NewEncoder(w).Encode(article)
        }
    }
	fmt.Println("Updated")
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    key := vars["id"]

	// Loop over all Articles
    // if the article.Id equals the key passed in
    // return the article encoded as JSON
    for _, article := range Articles {
        if article.Id == key {
            json.NewEncoder(w).Encode(article)
        }
    }
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint hit: homePage")
}

func handleRequest() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
	myRouter.HandleFunc("/articles", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/articles/{id}", returnSingleArticle).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	Articles = []Article{
        Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
        Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
    }
	handleRequest()
}
