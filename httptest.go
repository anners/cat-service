//http://thenewstack.io/building-a-web-server-in-go/
package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")

}

type Message struct {
	Text string
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
	http.ListenAndServe(":8080", nil)
}
