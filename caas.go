//http://thenewstack.io/building-a-web-server-in-go/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"

)

type Message struct {
	Text string
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")

}

// display random cat pictures 
func catpic(w http.ResponseWriter, r *http.Request) {

	// url for image search
	catapi := "http://localhost:8888//image?search=cat"
	request, err := http.Get(catapi)
	if err != nil {
		panic(err)	
	}
	defer request.Body.Close()

	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	data := make(map[string][]map[string]string)
	//var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err.Error())
	}
	
	// get random image from the data returned
	var image string 
	randomIndex := rand.Intn(len(data["items"])-1)    	
	randocat := data["items"][randomIndex] 
	for _, url := range randocat {
		image = url
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<img src=%s>", image)
}

func cat(w http.ResponseWriter, r *http.Request) {
	m := Message{"MEOW, purrrrr"}
	b, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}

	w.Write(b)
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/cat", cat)
	http.HandleFunc("/catpic", catpic)
	http.ListenAndServe(":8080", nil)
}
