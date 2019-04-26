package main

import (
	"company/companyConstants"
	"fmt"
)

type buyProductOp struct {
	product chan int
}

func productAddFilter(addChan chan int, products []int) chan int {
	if len(products) < cap(products) {
		return addChan
	}
	return nil
}

func buyFilter(buyChan chan *buyProductOp, products []int) chan *buyProductOp {
	if len(products) > 0 {
		return buyChan
	}
	return nil
}

func magazineServer(productPutChan chan int, productBuyChan chan *buyProductOp, infoChan chan bool) {
	var products = make([]int, 0, companyConstants.SizeOfMagazine)

	for {
		select {
		case newProduct := <-productAddFilter(productPutChan, products):
			products = append(products, newProduct)
		case buy := <-buyFilter(productBuyChan, products):
			buy.product <- products[len(products)-1]
			products = products[:len(products)-1]
		case <-infoChan:
			displayMagazine(products)
			infoChan <- true
		}
	}
}

func displayMagazine(products []int) {
	fmt.Println()
	fmt.Println("Magazine:")
	for _, product := range products {
		fmt.Println(product)
	}
	fmt.Println()
}
