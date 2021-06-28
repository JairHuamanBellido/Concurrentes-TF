package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Packet struct {
	Message int `json:"name"`
}

type Host struct {
	Name   string      `json:"host"`
	IsBusy bool        `json:"isBusy"`
	Memory MemoryStats `json:"memory"`
}

type MemoryStats struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Free  float64 `json:"free"`
}

var hosts []*Host = []*Host{
	&Host{Name: "http://192.168.1.15:3000/", Memory: MemoryStats{}, IsBusy: false},
	&Host{Name: "http://192.168.1.16:3000/", Memory: MemoryStats{}, IsBusy: false},
	&Host{Name: "http://192.168.1.17:3000/", Memory: MemoryStats{}, IsBusy: false},
	&Host{Name: "http://192.168.1.18:3000/", Memory: MemoryStats{}, IsBusy: false},
	&Host{Name: "http://192.168.1.19:3000/", Memory: MemoryStats{}, IsBusy: false},
}

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

func ocuppyHost(host chan Host, freeHost Host) {
	host <- freeHost
}

func freeHost(host Host) {
	for _, v := range hosts {
		if v.Name == host.Name {
			*&v.IsBusy = false
		}
	}
}

func checkAllBusy() {
	var c int = 0

	for _, v := range hosts {
		if v.IsBusy {
			c++
		}
	}

	if len(hosts) == c {
		fmt.Println("Todos ocupados")
		for _, v := range hosts {
			if v.IsBusy {
				*&v.IsBusy = false

			}

		}
	}
}

func packet(res http.ResponseWriter, req *http.Request) {

	enableCors(&res)
	host := make(chan Host)
	path := req.URL

	for _, v := range hosts {
		if v.IsBusy == false {
			*&v.IsBusy = true
			go ocuppyHost(host, *v)
			break
		}
	}
	p := <-host

	fmt.Println("Solicitud enviada a", p.Name)
	actual_host_recepted := fmt.Sprintf("%v%v", p.Name, path.String())
	resp, _ := http.Get(actual_host_recepted)

	defer (resp).Body.Close()
	// time.Sleep(1 * time.Second)

	json.NewEncoder(res).Encode(getBodyResponse(Packet{}, *resp))
	freeHost(p)
}

func handleRequest() {
	http.HandleFunc("/", packet)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {

	handleRequest()

}
