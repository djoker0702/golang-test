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
	//http.HandleFunc("/post", handlepost)
	http.HandleFunc("/helloworld", helloWorldHandler)
	fmt.Println("Server now running on localhost:8080")
	fmt.Println(`Try running: curl -X POST -d '{"hello":"test123"}' http://localhost:8080/helloworld`)
	log.Fatal(http.ListenAndServe(":8088", nil))
}

type helloWorldRequest struct {
	Hello string `json:"hello"`
}



func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	
		http.ServeFile(w,r,"../client/index.html")
}
/*func handlepost(w http.ResponseWriter, r *http.Request)  {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			log.Printf(string(body))

}*/

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
