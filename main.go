package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"os"
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

const clientID = "clientID"
const clientSecret = "clientSecret"

func main() {

	///oAuth start/////
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// We will be using `httpClient` to make external HTTP requests later in our code
	httpClient := http.Client{}

	// Create a new redirect route route
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		// Next, lets for the HTTP request to call the github oauth enpoint
		// to get our access token
		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		// We set this header since we want the response
		// as JSON
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer res.Body.Close()

		// Parse the request body into the `OAuthAccessResponse` struct
		var t OAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		// Finally, send a response to redirect the user to the "welcome" page
		// with the access token
		w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)

	///oAuth End//////

	//String to be hashed for hashing practice
	name := "Ashley"

	//Different types of hashing. 
	//Using the Sum function to computer a hash from a byte
	//Then we pass the string to be hashed into the byte
	hmd5 := md5.Sum([]byte(name))
	hsha1 := sha1.Sum([]byte(name))
	hsha2 := sha256.Sum256([]byte(name))

	//printing each of the different types of hash to the console
	fmt.Printf("   MD5: %x\n", hmd5)
	fmt.Printf("  SHA1: %x\n", hsha1)
	fmt.Printf("SHA256: %x\n", hsha2)
	
	//Creating json data
	Articles = []Article{
        Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
        Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
    }
	handleRequest()
}

////for oauth/////
type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
