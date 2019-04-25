package main

import (
	"company/companyConstants"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {

	isInteractive := false
	rand.Seed(time.Now().UnixNano())

	taskAddChan := make(chan Task)
	taskGetChan := make(chan *getTaskOp)
	taskListInfoChan := make(chan bool)

	productPutChan := make(chan int)
	productBuyChan := make(chan *buyProductOp)
	magazineInfoChan := make(chan bool)

	done := make(chan bool, 1)

	if len(os.Args) > 1 && os.Args[1] == "-i" {
		isInteractive = true
		go interact(taskListInfoChan, magazineInfoChan, done)
		go taskListServer(taskAddChan, taskGetChan, taskListInfoChan)
		go magazineServer(productPutChan, productBuyChan, magazineInfoChan)
	} else {
		go taskListServer(taskAddChan, taskGetChan, nil)
		go magazineServer(productPutChan, productBuyChan, nil)
	}

	go executiveOfficer(taskAddChan, isInteractive)

	for i := 0; i < companyConstants.NumberOfWorkers; i++ {
		go worker(taskGetChan, productPutChan, isInteractive)
	}

	go customer(productBuyChan, isInteractive)

	if isInteractive {
		<-done
	} else {
		_, _ = fmt.Scanln()
	}
}
