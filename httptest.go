//http://thenewstack.io/building-a-web-server-in-go/
//testing
package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")

}

type Message struct {
	Text string
}

// catapi NTMzNzg
func catpic(w http.ResponseWriter, r *http.Request) {
	catapi := "http://thecatapi.com/api/images/get?format=html"
	request, err := http.Get(catapi)
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()
	content, err := ioutil.ReadAll(request.Body)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(content))
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
