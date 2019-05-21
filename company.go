package main

import (
	"company/companyConstants"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	var isInteractive bool
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "-i" {
		isInteractive = true
	} else {
		isInteractive = false
	}

	additionMachines, multiplicationMachines := createMachines()

	taskAddChan := make(chan taskMachineAdapter)
	taskGetChan := make(chan *getTaskOp)
	taskListInfoChan := make(chan bool)

	productPutChan := make(chan int)
	productBuyChan := make(chan *buyProductOp)
	magazineInfoChan := make(chan bool)

	brokenMachinesChan := make(chan int)

	machines := make([]machine, 0, companyConstants.AdditionMachinesCount+companyConstants.MultiplicationMachinesCount)

	for i := range additionMachines {
		machines = append(machines, additionMachines[i])
	}
	for i := range multiplicationMachines {
		machines = append(machines, multiplicationMachines[i])
	}

	go service(brokenMachinesChan, machines, isInteractive)

	workersInfoChan := make([]chan int, 0, companyConstants.WorkersCount)

	done := make(chan bool, 1)

	for i := 0; i < companyConstants.WorkersCount; i++ {
		workerInfoChan := make(chan int)
		workersInfoChan = append(workersInfoChan, workerInfoChan)
		go worker(i, taskGetChan, productPutChan, brokenMachinesChan, workerInfoChan, isInteractive)
	}

	if isInteractive {
		go interact(taskListInfoChan, magazineInfoChan, workersInfoChan, done)
		go taskListServer(taskAddChan, taskGetChan, taskListInfoChan)
		go magazineServer(productPutChan, productBuyChan, magazineInfoChan)
	} else {
		go taskListServer(taskAddChan, taskGetChan, nil)
		go magazineServer(productPutChan, productBuyChan, nil)
	}

	go executiveOfficer(taskAddChan, isInteractive, multiplicationMachines, additionMachines)

	go customer(productBuyChan, isInteractive)

	if isInteractive {
		<-done
	} else {
		_, _ = fmt.Scanln()
	}
}
