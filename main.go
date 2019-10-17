package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Object struct {
	ID    string `json:"ID"`
	Name  string `json:"Name"`
	Title string `json:"Title"`
}

/*
this array of Objects is used to similulate a database
can be easily added later
*/
var Objects []Object

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!\n")
	fmt.Println("Endpoint hit: Home")
}

//returns all objects from array formatted in JSON
func returnAllObjects(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: Returned all objects")
	json.NewEncoder(w).Encode(Objects)
}

//returns a single object based on ID from the array formatted in JSON
func returnSingleObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["ID"]
	for _, object := range Objects {
		if object.ID == key {
			json.NewEncoder(w).Encode(object)
		}
	}
}

func createObject(w http.ResponseWriter, r *http.Request) {
	//get POST request and return response containing the request body
	reqBody, err := ioutil.ReadAll(r.Body)
	var object Object
	if err != nil {
		fmt.Println(w, "Invalid")
	}
	json.Unmarshal(reqBody, &object)
	Objects = append(Objects, object)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func deleteObject(w http.ResponseWriter, r *http.Request) {
	/*
		vars := mux.Vars(r)
		id := vars["ID"]
	*/
	objectID := mux.Vars(r)["ID"]
	for i, object := range Objects {
		if object.ID == objectID {
			Objects = append(Objects[:i], Objects[i+1:]...)
			fmt.Fprintf(w, "Object %v has been successfully deleted.", objectID)
		}
	}
}

func updateObject(w http.ResponseWriter, r *http.Request) {
	objectID := mux.Vars(r)["ID"]
	var objectUpdate Object
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid")
	}
	json.Unmarshal(reqBody, &objectUpdate)
	for i, object := range Objects {
		if object.ID == objectID {
			object.Name = objectUpdate.Name
			object.Title = objectUpdate.Title
			Objects = append(Objects[:i], object)
		}
	}
}

func requestHandler() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/object", createObject).Methods("POST")
	myRouter.HandleFunc("/object", returnAllObjects).Methods("GET")
	myRouter.HandleFunc("/objects/{ID}", returnSingleObject).Methods("GET")
	myRouter.HandleFunc("/objects/{ID}", updateObject).Methods("POST")
	myRouter.HandleFunc("/objects/{ID}", deleteObject).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	//populate array with dummy data
	Objects = []Object{
		Object{ID: "1", Name: "John Doe", Title: "Vice President of Sales"},
		Object{ID: "2", Name: "Jane Doe", Title: "Vice President of Marketing"},
	}
	requestHandler()
}
