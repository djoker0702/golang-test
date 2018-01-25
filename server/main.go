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
var mappedData = &Data{CopyAndPaste :map[string]bool {"inputEmail": false,"inputCardNumber": false,"inputCVV": false}} 	
func eventHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to read body"))
		return
	}

	req :=  &event{}
	if err = json.Unmarshal(body, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to unmarshal JSON request"))
		return
	}
	mappedData.WebsiteUrl = req.WebsiteUrl
	mappedData.SessionId = req.Session
	switch event := req.EventType; event {

	case "screenResize" :
		mappedData.ResizeFrom = req.ODimension	
		mappedData.ResizeTo = req.NDimension
		log.Printf("\nEvent : screenResize , Current state of the data :%+v", 
		"\nWebsiteUrl :%+v " , mappedData.WebsiteUrl,
		"\nSessionId :%+v" , mappedData.SessionId,
		"\nResizeFrom :%+v" , mappedData.ResizeFrom,
		"\nResizeTo :%+v" , mappedData.ResizeTo)
	case "timeTaken": 
		mappedData.FormCompletionTime = req.Time
		mappedData.ResizeFrom = req.ODimension	
		mappedData.ResizeTo = req.NDimension
		log.Printf("\nEvent : Form submitted , Final state of the data :  %+v", 
			"\nWebsiteUrl : ", mappedData.WebsiteUrl,
			"\nSessionId : ", mappedData.SessionId,
			"\nResizeFrom : ", mappedData.ResizeFrom,
			"\nResizeTo : ", mappedData.ResizeTo,
			"\nCopyAndPaste : ", mappedData.CopyAndPaste,
			"\nFormCompletionTime : ", mappedData.FormCompletionTime)	
	case "copyAndPaste":
		mappedData.CopyAndPaste[req.FormId] = req.Copie || req.Paste
		log.Printf("\nEvent : Copy and Paste detected , Current state of the data :  %+v", 
			"\nWebsiteUrl : ", mappedData.WebsiteUrl,
			"\nSessionId : ", mappedData.SessionId,
			"\nCopyAndPaste : ", mappedData.CopyAndPaste)


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
