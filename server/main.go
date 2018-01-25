package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", handleHttpRequest)
	http.HandleFunc("/helloworld", helloWorldHandler)
	http.HandleFunc("/api/event", eventHandler)
	fmt.Println("Server now running on localhost:3000")
	fmt.Println(`Try running: curl -X POST -d '{"hello":"test123"}' http://localhost:3000/helloworld`)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type helloWorldRequest struct {
	Hello string `json:"hello"`
}



func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	
		http.ServeFile(w,r,"../client/index.html")
}

type test_struct struct {
	Test string
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("a post detected")
	decoder := json.NewDecoder(r.Body)
	var t test_struct
	err := decoder.Decode(&t)
	
		if err != nil {
			panic(err)
		}
		log.Println(t.Test)
    

}



func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read body"))
		return
	}

	req := &helloWorldRequest{}

	if err = json.Unmarshal(body, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to unmarshal JSON request"))
		return
	}

	log.Printf("Request received %+v", req)

	w.WriteHeader(http.StatusOK)
}
