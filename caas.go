package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Message struct {
	Text string
}

type Service struct {
	Node                     string   `json:"Node"`
	Address                  string   `json:"Address"`
	ServiceID                string   `json:"ServiceID"`
	ServiceName              string   `json:"ServiceName"`
	ServiceTags              []string `json:"ServiceTags"`
	ServiceAddress           string   `json:"ServiceAddress"`
	ServicePort              int      `json:"ServicePort"`
	ServiceEnableTagOverride bool     `json:"ServiceEnableTagOverride"`
	CreateIndex              int      `json:"CreateIndex"`
	ModifyIndex              int      `json:"ModifyIndex"`
}

// this should be it's own package but for now it's not
func getConsulService(url string) []Service {

	var service []Service
	// disable security check (at your own risk)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	request, err := client.Get(url)
	if err != nil {
		log.Println("WARNING: ", err)
		return service
	}
	defer request.Body.Close()

	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &service)
	if err != nil {
		panic(err.Error())
	}

	return service
}

func downloadCatpic(url string) string {
	extension := strings.Split(url, ".")
	filename := "cat." + extension[len(extension)-1]

	// create filename
	file, err := os.Create("/tmp/" + filename)
	if err != nil {
		fmt.Println("Error while creating", filename, "-", err)
		panic(err)
	}
	defer file.Close()

	request, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		panic(err)
	}
	defer request.Body.Close()

	_, err = io.Copy(file, request.Body)
	if err != nil {
		fmt.Println("Error copying file", file, "-", err)
		panic(err)
	}

	return file.Name()
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

// display random cat pictures
func pic(w http.ResponseWriter, r *http.Request) {

	// testing ELK
	log.Println("YAY I AM WORKING")

	// set default imageURL in case something fails
	imageURL := "http://i.dailymail.co.uk/i/pix/2014/08/05/1407225932091_wps_6_SANTA_MONICA_CA_AUGUST_04.jpg"
	// set default image api
	catapi := "http://image-service.apps.ciscocloud.io/image?search=cat"

	consulURL := "http://consul.service.consul:8500/v1/catalog/service/image-service"
	service := getConsulService(consulURL)

	if service != nil {
		// construct the catapi
		catapi = "http://image-service.service.consul:" + strconv.Itoa(service[0].ServicePort) + "/image?search=beer"
	}

	request, err := http.Get(catapi)
	if err != nil {
		log.Println("problem with fetching", catapi, ":", err)
	} else {
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
		randomIndex := rand.Intn(len(data["items"]) - 1)
		randocat := data["items"][randomIndex]
		for _, url := range randocat {
			imageURL = url
		}
	}

	image := downloadCatpic(imageURL)
	w.Header().Set("Content-Type", "image/jpg")
	http.ServeFile(w, r, image)
	//fmt.Fprintf(w, "<img src=%s width=\"500\" height=\"500\">", image)
}

func cat(w http.ResponseWriter, r *http.Request) {
	m := Message{"MEOW, purrrrr"}
	b, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func static(w http.ResponseWriter, r *http.Request) {
	urlSplit := strings.Split(r.URL.Path, "/")
	filename := urlSplit[len(urlSplit)-1]
	if len(filename) == 0 {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	dirFilename := "/tmp/" + filename

	w.Header().Set("Content-Type", "image/jpg")
	http.ServeFile(w, r, dirFilename)
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/cat", cat)
	http.HandleFunc("/pic", pic)
	http.HandleFunc("/static/", static)
	http.ListenAndServe(":8080", nil)
}
