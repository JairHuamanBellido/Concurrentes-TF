package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mackerelio/go-osstat/memory"
)

type HttpResponse struct {
	Message int `json:"name"`
}

type MemoryStats struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
	Free  float64 `json:"free"`
}

func getOsInformation(res http.ResponseWriter, req *http.Request) {
	enableCors(&res)

	switch req.Method {
	case http.MethodGet:
		memory, err := memory.Get()
		if err != nil {
			log.Fatal(err)
			return
		}

		var memoryStats MemoryStats = MemoryStats{
			Total: float64(memory.Total),
			Used:  float64(memory.Used),
			Free:  float64(memory.Free),
		}
		json.NewEncoder(res).Encode(memoryStats)
	}

}

func enableCors(res *http.ResponseWriter) {
	(*res).Header().Set("Access-Control-Allow-Origin", "*")
	(*res).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD")
	(*res).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	(*res).Header().Set("Content-Type", "application/json")
}

func packet(res http.ResponseWriter, req *http.Request) {

	enableCors(&res)

	fmt.Println("Recibiendo peticion")
	json.NewEncoder(res).Encode(HttpResponse{Message: 2})

}

func handleRequest(port string) {
	http.HandleFunc("/", packet)
	http.HandleFunc("/status", getOsInformation)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	var port string = ":3000"

	if len(os.Args) > 1 {
		port = fmt.Sprintf(":%v", os.Args[2])
	}

	handleRequest(port)

}
