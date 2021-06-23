package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"

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

type AtmosphereQuality struct {
	MonoxidoCarbono    float64
	AcidoSulfridico    float64
	DioxidoDeNitrogeno float64
	Ozono              float64
	PM10               float64
	PM25               float64
	DioxidoDeAzufre    float64
	Ruido              float64
	UV                 float64
	Humedad            float64
	Presion            float64
	Temperatura        float64
	distancia          float64
}

const K int = 1

var dataset = []AtmosphereQuality{}

func ecludian(model, input AtmosphereQuality) float64 {
	monoxidoCarbono_difference := math.Pow(input.MonoxidoCarbono-model.MonoxidoCarbono, 2)
	acidoSulfridico_difference := math.Pow(input.AcidoSulfridico-model.AcidoSulfridico, 2)
	dioxidoDeNitrogeno_difference := math.Pow(input.DioxidoDeNitrogeno-model.DioxidoDeNitrogeno, 2)
	ozono_difference := math.Pow(input.Ozono-model.Ozono, 2)
	pm10_difference := math.Pow(input.PM10-model.PM10, 2)
	pm25_difference := math.Pow(input.PM25-model.PM25, 2)
	dioxidoDeAzufre_difference := math.Pow(input.DioxidoDeAzufre-model.DioxidoDeAzufre, 2)
	ruido_difference := math.Pow(input.Ruido-model.Ruido, 2)
	UV_difference := math.Pow(input.UV-model.UV, 2)
	humedad_difference := math.Pow(input.Humedad-model.Humedad, 2)
	presion_difference := math.Pow(input.Presion-model.Presion, 2)

	return math.Sqrt(monoxidoCarbono_difference +
		acidoSulfridico_difference +
		dioxidoDeNitrogeno_difference +
		ozono_difference +
		pm10_difference +
		pm25_difference +
		dioxidoDeAzufre_difference +
		ruido_difference +
		UV_difference +
		humedad_difference +
		presion_difference)
}

func knn(dataset []AtmosphereQuality, input AtmosphereQuality, result chan float64) {

	var res []AtmosphereQuality
	for _, v := range dataset {
		distance := ecludian(v, input)
		v.distancia = distance
		res = append(res, v)
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].distancia < res[j].distancia
	})

	result <- res[:K][0].Temperatura

}

func parseRowToFloat(value string, err error) float64 {
	if err != nil {
		log.Fatal(err)
	}

	floatValue, _ := strconv.ParseFloat(value, 64)

	return floatValue
}

func getTemperature(res http.ResponseWriter, req *http.Request) {
	var model AtmosphereQuality = AtmosphereQuality{
		MonoxidoCarbono:    parseRowToFloat(req.FormValue("monoxidoCarbono"), nil),
		AcidoSulfridico:    parseRowToFloat(req.FormValue("acidoSulfridico"), nil),
		DioxidoDeNitrogeno: parseRowToFloat(req.FormValue("dioxidoDeNitrogeno"), nil),
		Ozono:              parseRowToFloat(req.FormValue("ozono"), nil),
		PM10:               parseRowToFloat(req.FormValue("pm10"), nil),
		PM25:               parseRowToFloat(req.FormValue("pm25"), nil),
		DioxidoDeAzufre:    parseRowToFloat(req.FormValue("dioxidoDeAzufre"), nil),
		Ruido:              parseRowToFloat(req.FormValue("ruido"), nil),
		UV:                 parseRowToFloat(req.FormValue("uv"), nil),
		Humedad:            parseRowToFloat(req.FormValue("humedad"), nil),
		Presion:            parseRowToFloat(req.FormValue("presion"), nil),
	}

	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	res.Header().Set("Content-Type", "application/json")

	var temperature chan float64 = make(chan float64)

	go knn(dataset, model, temperature)

	model.Temperatura = <-temperature

	json.NewEncoder(res).Encode(model)

}

func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
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

func handleRequest(port string) {
	http.HandleFunc("/", getTemperature)
	http.HandleFunc("/status", getOsInformation)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	var port string = ":3000"

	if len(os.Args) > 1 {
		port = fmt.Sprintf(":%v", os.Args[2])
	}
	csvFile, _ := readCSVFromUrl("https://raw.githubusercontent.com/JairHuamanBellido/dataset-concurrente/main/dataset.csv")

	for i := 2; i < len(csvFile); i++ {
		monoxidoCarbono := parseRowToFloat(csvFile[i][6], nil)
		acidoSulfridico := parseRowToFloat(csvFile[i][7], nil)
		dioxidoDeNitrogeno := parseRowToFloat(csvFile[i][8], nil)
		ozono := parseRowToFloat(csvFile[i][9], nil)
		pm10 := parseRowToFloat(csvFile[i][10], nil)
		pm25 := parseRowToFloat(csvFile[i][11], nil)
		dioxidoDeAzufre := parseRowToFloat(csvFile[i][12], nil)
		ruido := parseRowToFloat(csvFile[i][13], nil)
		UV := parseRowToFloat(csvFile[i][14], nil)
		humedad := parseRowToFloat(csvFile[i][15], nil)
		presion := parseRowToFloat(csvFile[i][18], nil)
		temperatura := parseRowToFloat(csvFile[i][19], nil)

		var model AtmosphereQuality = AtmosphereQuality{
			MonoxidoCarbono:    monoxidoCarbono,
			AcidoSulfridico:    acidoSulfridico,
			DioxidoDeNitrogeno: dioxidoDeNitrogeno,
			Ozono:              ozono,
			PM10:               pm10,
			PM25:               pm25,
			DioxidoDeAzufre:    dioxidoDeAzufre,
			Ruido:              ruido,
			UV:                 UV,
			Humedad:            humedad,
			Presion:            presion,
			Temperatura:        temperatura,
		}

		dataset = append(dataset, model)

	}

	handleRequest(port)

}
