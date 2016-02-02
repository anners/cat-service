package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strconv"
	"net/http"
)

type Message struct {
	Text string
}

type Service struct {
	Node string `json:"Node"`
	Address string `json:"Address"`
	ServiceID string `json:"ServiceID"`
	ServiceName string `json:"ServiceName"`
	ServiceTags []string `json:"ServiceTags"`
	ServiceAddress string `json:"ServiceAddress"`
	ServicePort int `json:"ServicePort"`
	ServiceEnableTagOverride bool `json:"ServiceEnableTagOverride"`
	CreateIndex int `json:"CreateIndex"`
	ModifyIndex int `json:"ModifyIndex"`
}

// this should be it's own package but for now it's not
func getConsulService(url string) []Service{

	// disable security check (at your own risk)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	request, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()

	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
	panic(err)
	}

	var service []Service
	err = json.Unmarshal(content, &service)
	if err != nil {
		panic(err.Error())
	}

	return service
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")

}

// display random cat pictures 
func catpic(w http.ResponseWriter, r *http.Request) {

	consulURL := "http://consul.service.consul:8500/v1/catalog/service/image-service"
	service := getConsulService(consulURL)

	// set default image in case something fails and still makes it to the iamage
	image := "http://i.dailymail.co.uk/i/pix/2014/08/05/1407225932091_wps_6_SANTA_MONICA_CA_AUGUST_04.jpg"

	// url for image search
	catapi := "http://image-service.service.consul:" + strconv.Itoa(service[0].ServicePort) + "/image?search=grumpy+cat"
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
