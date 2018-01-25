package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

type Data struct {
	WebsiteUrl         string 
	SessionId          string 
	ResizeFrom         Dimension 
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool// map[fieldId]true
	FormCompletionTime int// Seconds
}

type Dimension struct {
	Width  int 
	Height int 
}

// this struct is the general format of all the json request that came from the frontend side
// the Data struct will be built from the event struct (for each event that occur depending on the event type field )
type event struct {
	EventType		   string `json:"eventType"`
	WebsiteUrl		   string `json:"websiteUrl"`
	Session		   	   string `json:"sessionId"`
	ODimension         Dimension `json:"ODimension"`
	NDimension		   Dimension `json:"NDimension"`
	Copie			   bool `json:"copied"`
	Paste              bool `json:"pasted"`
	FormId			   string `json:"formId"`
	Time               int `json:"time"`
}
var mappedData = &Data{}
func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	
		http.ServeFile(w,r,"../client/index.html")
		time.Sleep(2 * time.Second)  //to handle concurrent requests
		/* the Data struct should customize a unique session,
		 that's why we declare it when the a new root (/) request is made
		which means that a new user opened a new session		
		and we define the CopyAndPaste map to false for each field of the form 
		*/
		mappedData = &Data{CopyAndPaste :map[string]bool {"inputEmail": false,"inputCardNumber": false,"inputCVV": false}} 	
		
		
}

/*
	this function will handle all the events that occurs (as defined in the readme) 

*/
func eventHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read body"))
		return
	}

	req :=  &event{} // the req is used to fetch data from the json that comes with the post request
	if err = json.Unmarshal(body, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to unmarshal JSON request"))
		return
	}
	// mappedData as defined in line 63 is the Data struct that will be printed for each stage of its construction
	mappedData.WebsiteUrl = req.WebsiteUrl
	mappedData.SessionId = req.Session
	switch event := req.EventType; event {

	case "screenResize" :
		mappedData.ResizeFrom = req.ODimension	
		mappedData.ResizeTo = req.NDimension
		log.Printf("\nEvent : screenResize , Current state of the data :%+v", mappedData)

	case "timeTaken": 
		mappedData.FormCompletionTime = req.Time
		mappedData.ResizeFrom = req.ODimension	
		mappedData.ResizeTo = req.NDimension
		log.Printf("\nEvent : Form submitted , Final state of the data :  %+v", mappedData)
	case "copyAndPaste":
		mappedData.CopyAndPaste[req.FormId] = req.Copie || req.Paste
		log.Printf("\nEvent : Copy and Paste detected , Current state of the data :  %+v", mappedData)

	}
	w.WriteHeader(http.StatusOK)
    

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
