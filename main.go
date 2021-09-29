package main

import (
	"fmt"
	"sync"

	Utils "github.com/bayronaz/LTV-Tecnhical-test/Helpers"
)

// albums slice to seed record album data.
var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	fmt.Println("test")
	var wg sync.WaitGroup
	wg.Add(1)

	go getResult(&wg)

	wg.Wait()
	//var res = Utils.PrintText()
	//fmt.Println(res)
}

func getResult(wg *sync.WaitGroup) {
	var res = Utils.PrintText()
	fmt.Println(res)
	wg.Done()
}
