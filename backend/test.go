package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Mesero struct {
	Name    string
	Ocupado bool
}

var meseros []*Mesero = []*Mesero{
	&Mesero{Name: "Mesero 1", Ocupado: false},
	&Mesero{Name: "Mesero 2", Ocupado: false},
	&Mesero{Name: "Mesero 3", Ocupado: false},
	&Mesero{Name: "Mesero 4", Ocupado: false},
}

func verificar(d, ch chan bool) {
	val := <-ch
	// d <- true
	fmt.Println("xd: ", val)

}
func seleccionarMesero(mes chan Mesero, mesero Mesero) {
	mes <- mesero
}
func liberarMesero(mesero Mesero) {
	for _, v := range meseros {
		if v.Name == mesero.Name {
			*&v.Ocupado = false
		}
	}
}

func checkAllBusy() {
	var c int = 0

	for _, v := range meseros {
		if v.Ocupado {
			c++
		}
	}

	if len(meseros) == c {
		fmt.Println("Todos ocupados")
		for _, v := range meseros {
			if v.Ocupado {
				*&v.Ocupado = false

			}

		}
	}
}

func getTemperature(res http.ResponseWriter, req *http.Request) {

	mes := make(chan Mesero)
	checkAllBusy()
	for _, v := range meseros {
		if v.Ocupado == false {
			*&v.Ocupado = true
			go seleccionarMesero(mes, *v)
			break
		}

	}

	actual_mes := <-mes
	// time.Sleep(1 * time.Second)

	var message string = fmt.Sprintf("Atendido por %v", actual_mes.Name)

	json.NewEncoder(res).Encode(message)
	liberarMesero(actual_mes)

}

func handleRequest() {
	http.HandleFunc("/", getTemperature)
	log.Fatal(http.ListenAndServe(":3030", nil))
}

func main() {

	handleRequest()
}
