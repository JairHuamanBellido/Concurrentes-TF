package main

import (
	"fmt"
	"math"
	"sort"
)

type Class struct {
	pointX   float64
	pointY   float64
	label    string
	distance float64
}

const K int = 3

var target Class = Class{pointX: 53, pointY: 1.68, label: "sin definir"}

func ecludian(x1, y1 float64) float64 {
	x_difference := math.Pow(target.pointX-x1, 2)
	y_difference := math.Pow(target.pointY-y1, 2)
	fmt.Println(x_difference)
	return math.Sqrt(x_difference + y_difference)
}

func knn(dataset []Class) {

	var res []Class
	for _, v := range dataset {
		distance := ecludian(v.pointX, v.pointY)
		res = append(res, Class{pointX: v.pointX, pointY: v.pointY, label: v.label, distance: distance})
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].distance < res[j].distance
	})

	var countFather int = 0
	var countChild int = 0

	for _, v := range res[:K] {
		if v.label == "niño" {
			countChild++
		} else {
			countFather++
		}
	}

	if countChild > countFather {
		target.label = "niño"
	} else {
		target.label = "adulto"

	}

}

func main() {
	var dataset = []Class{
		{pointX: 49, pointY: 1.43, label: "niño", distance: 0},
		{pointX: 51, pointY: 1.55, label: "niño", distance: 0},
		{pointX: 57, pointY: 1.58, label: "niño", distance: 0},
		{pointX: 47, pointY: 1.55, label: "niño", distance: 0},
		{pointX: 54, pointY: 1.60, label: "niño", distance: 0},
		{pointX: 56, pointY: 1.58, label: "niño", distance: 0},
		{pointX: 59, pointY: 1.64, label: "niño", distance: 0},
		{pointX: 53, pointY: 1.61, label: "niño", distance: 0},
		{pointX: 58, pointY: 1.63, label: "niño", distance: 0},
		{pointX: 52, pointY: 1.60, label: "adulto", distance: 0},
		{pointX: 75, pointY: 1.73, label: "adulto", distance: 0},
		{pointX: 80, pointY: 1.75, label: "adulto", distance: 0},
		{pointX: 75, pointY: 1.69, label: "adulto", distance: 0},
		{pointX: 65, pointY: 1.71, label: "adulto", distance: 0},
		{pointX: 75, pointY: 1.79, label: "adulto", distance: 0},
		{pointX: 77, pointY: 1.76, label: "adulto", distance: 0},
		{pointX: 65, pointY: 1.71, label: "adulto", distance: 0},
		{pointX: 70, pointY: 1.70, label: "adulto", distance: 0},
		{pointX: 78, pointY: 1.81, label: "adulto", distance: 0},
		{pointX: 70, pointY: 1.67, label: "adulto", distance: 0},
	}

	knn(dataset)

	fmt.Println(target)

}
