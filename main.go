package shipping

import (
	"fmt"
	"sync"
)

var (
	min = float64(9999999999999)
)

type Product struct {
	Name string
	Len  int
	Wid  int
	Hei  int
}

type Box struct {
	Len int
	Wid int
	Hei int
}

type TotalProductDimension struct {
	TotalLen int
	TotalWid int
	TotalHei int
}

type Result struct {
	Box     Box
	Fitness float64
}

func main() {
	availableBoxes := []Box{
		{Len: 5, Wid: 10, Hei: 7},
		{Len: 9, Wid: 10, Hei: 37},
		{Len: 1, Wid: 1, Hei: 2},
		{Len: 2, Wid: 7, Hei: 2},
		{Len: 6, Wid: 40, Hei: 3},
		{Len: 3, Wid: 4, Hei: 5},
		{Len: 3, Wid: 4, Hei: 2},
		{Len: 3, Wid: 3, Hei: 3},
	}
	products := []Product{
		{Name: "product 1", Len: 1, Wid: 1, Hei: 3},
		{Name: "product 2", Len: 2, Wid: 1, Hei: 1},
		{Name: "product 3", Len: 1, Wid: 1, Hei: 2},
		{Name: "product 4", Len: 2, Wid: 3, Hei: 3},
	}

	result := getBestBox(availableBoxes, products)
	fmt.Println(result)
}

func getBestBox(availableBoxes []Box, products []Product) Box {

	// 1. calculate the total dimension of all the products
	t := calculateAggregateDimension(products)
	input := make(chan Result)
	output := make(chan Box)

	defer close(output)
	var wg sync.WaitGroup

	go handleResult(&wg, input, output)

	// 2. find the most fit
	for _, box := range availableBoxes {
		wg.Add(1)
		go checkFittness(box, t, input)
	}

	wg.Wait()
	close(input)

	result := <-output

	return result
}

func checkFittness(b Box, t float64, output chan Result) {
	vB := b.Len * b.Wid * b.Hei

	d := float64(vB) - t
	
	output <- Result{
		Box:     b,
		Fitness: d,
	}
}

func calculateAggregateDimension(products []Product) float64 {
	var result float64
	for _, product := range products {
		vP := product.Len * product.Hei * product.Wid
		result += float64(vP)
	}
	return result
}

func handleResult(wg *sync.WaitGroup, input chan Result, output chan Box) {
	var result Box
	
	for incomingEvent := range input {
		if incomingEvent.Fitness >= 0 {
			if incomingEvent.Fitness <= min {
				min = incomingEvent.Fitness
				result = incomingEvent.Box
			}
		}
		wg.Done()
	}

	output <- result
}
