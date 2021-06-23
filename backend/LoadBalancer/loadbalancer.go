package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type Packet struct {
	Message int `json:"name"`
}

type Host struct {
	Name   string `json:"host"`
	Memory MemoryStats
}

type MemoryStats struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Free  float64 `json:"free"`
}

var hosts []string = []string{"http://192.168.1.15:3000/", "http://192.168.1.16:3000/", "http://192.168.1.17:3000/", "http://192.168.1.18:3000/", "http://192.168.1.19:3000/"}

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
		resp, _ := http.Get(v + "/status")
		defer (resp).Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyString := string(body)

		var response MemoryStats

		json.Unmarshal([]byte(bodyString), &response)
		targets = append(targets, Host{Name: v, Memory: response})

	}

	sort.SliceStable(targets, func(i, j int) bool {
		return (targets[i].Memory).Free > targets[j].Memory.Free
	})

	fmt.Print("\n====\n")
	for _, v := range targets {
		fmt.Println(v.Name, " : ", v.Memory.Free)
	}

	return targets[0].Name

}

func packet(res http.ResponseWriter, req *http.Request) {

	enableCors(&res)
	host := make(chan string)
	path := req.URL
	resp, _ := http.Get(RoundRobin(host) + path.String())
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
