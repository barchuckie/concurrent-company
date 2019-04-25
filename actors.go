package main

import (
	"company/companyConstants"
	"fmt"
	"math/rand"
	"time"
)

func executiveOfficer(taskAddChan chan Task, isInteractive bool) {
	sleepTime := time.Duration(companyConstants.CEORate) * time.Second

	for {
		time.Sleep(sleepTime)

		n := rand.Intn(3)
		arg1 := rand.Intn(100000)
		arg2 := rand.Intn(100000)

		var task Task

		switch n {
		case 0:
			task = AdditionTask{arg1, "+", arg2}
		case 1:
			task = SubtractionTask{arg1, "-", arg2}
		default:
			task = MultiplicationTask{arg1, "*", arg2}
		}

		taskAddChan <- task

		if !isInteractive {
			fmt.Println("CEO adds new task:", task)
		}
	}
}

func worker(taskGetChan chan *getTaskOp, productAddChan chan int, isInteractive bool) {
	sleepTime := time.Duration(companyConstants.WorkerRate) * time.Second
	for {
		time.Sleep(sleepTime)

		newTaskChan := &getTaskOp{
			response: make(chan Task),
		}
		taskGetChan <- newTaskChan
		newTask := <-newTaskChan.response
		product := newTask.solve()

		if !isInteractive {
			fmt.Println("Worker creates product:", product)
		}

		productAddChan <- product

		if !isInteractive {
			fmt.Println("Worker puts product:", product, "in magazine")
		}
	}
}

func customer(productBuyChan chan *buyProductOp, isInteractive bool) {
	sleepTime := companyConstants.CustomerRate * time.Second

	for {
		time.Sleep(sleepTime)

		productBuyStruct := &buyProductOp{
			product: make(chan int),
		}

		productBuyChan <- productBuyStruct
		boughtProduct := <-productBuyStruct.product

		if !isInteractive {
			fmt.Println("Customer buys product:", boughtProduct)
		}
	}
}
