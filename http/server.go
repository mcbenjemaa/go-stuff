package main

import (
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcbenjemaa/go-stuff/internal/album"
)

func allAlbums(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: All Albums Endpoint")
	json.NewEncoder(w).Encode(album.Albums)
}

func saveAlbum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: saveAlbum Endpoint")
	// get the body of our POST request
	// return the string response containing the request body

	var alb album.Album
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&alb) // OR -
	// reqBody, _ := ioutil.ReadAll(r.Body)
	//err := json.Unmarshal(reqBody, &album)
	if err != nil {
		log.Printf("Error %w", err)
	}
	log.Printf("Result: %#v", alb)
	album.AddAlbum(alb)
	fmt.Fprintf(w, "%+v", alb)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage !")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/albums", saveAlbum).Methods("POST")
	router.HandleFunc("/albums", allAlbums)
	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	handleRequests()
}
