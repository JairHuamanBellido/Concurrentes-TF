package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Packet struct {
	Message int `json:"name"`
}

type Host struct {
	Name   string `json:"host"`
	Memory interface{}
}

type MemoryStats struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Free  float64 `json:"free"`
}

var hosts []string = []string{"http://192.168.1.8:3031/", "http://192.168.1.8:3032/", "http://192.168.1.8:3033/", "http://192.168.1.8:3034/", "http://192.168.1.8:3035/"}

func enableCors(res *http.ResponseWriter) {
	(*res).Header().Set("Access-Control-Allow-Origin", "*")
	(*res).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD")
	(*res).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	(*res).Header().Set("Content-Type", "application/json")
}

func getBodyResponse(model interface{}, resp http.Response) interface{} {
	body, _ := ioutil.ReadAll(resp.Body)

	bodyString := string(body)

	var response interface{}

	json.Unmarshal([]byte(bodyString), &response)
	return response
}

func RoundRobin(c chan string) string {

	var targets []Host
	for _, v := range hosts {
		resp, _ := http.Get(v + "status")
		defer (resp).Body.Close()
		memorys := getBodyResponse(MemoryStats{}, *resp)

		targets = append(targets, Host{Name: v, Memory: memorys})

	}

	return targets[3].Name

}

func packet(res http.ResponseWriter, req *http.Request) {

	enableCors(&res)
	host := make(chan string)

	resp, _ := http.Get(RoundRobin(host))
	defer (resp).Body.Close()

	json.NewEncoder(res).Encode(getBodyResponse(Packet{}, *resp))

}

func handleRequest() {
	http.HandleFunc("/", packet)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {

	handleRequest()

}
