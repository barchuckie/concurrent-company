package main

import (
	"company/companyConstants"
	"fmt"
	"math/rand"
	"time"
)

func executiveOfficer(taskAddChan chan taskMachineAdapter, isInteractive bool,
	multiplicationMachines []multiplicationMachine, additionMachines []additionMachine) {

	sleepTime := time.Duration(companyConstants.CEORate) * time.Second
	x := 0

	for {
		time.Sleep(sleepTime)

		n := rand.Intn(2)
		arg1 := rand.Intn(100000)
		arg2 := rand.Intn(100000)
		x++

		var task taskMachineAdapter

		switch n {
		case 0:
			task = additionAdapter{
				additionTask{arg1, "+", arg2, nil},
				additionMachines,
			}
		default:
			task = multiplicationAdapter{
				multiplicationTask{arg1, "*", arg2, nil},
				multiplicationMachines,
			}
		}

		if !isInteractive {
			fmt.Print("CEO adds new task: ")
			task.getTask().print()
		}

		taskAddChan <- task
	}
}

func worker(id int, taskGetChan chan *getTaskOp, productPutChan chan int, brokenMachineChan chan int,
	infoChan chan int, isInteractive bool) {

	sleepTime := time.Duration(companyConstants.WorkerRate) * time.Second
	var isPatient bool
	var solvedTaskCounter = 0
	var solvedTask task

	if rand.Intn(2) == 1 {
		isPatient = true
	} else {
		isPatient = false
	}

	go func() {
		for {
			reqId := <-infoChan
			if reqId == id {
				printWorkerStatistics(id, isPatient, solvedTaskCounter)
			} else {
				infoChan <- reqId
			}
			infoChan <- reqId
		}
	}()

	for {
		time.Sleep(sleepTime)

		/* Getting task from task list */
		newTaskChan := &getTaskOp{
			response: make(chan taskMachineAdapter),
		}
		taskGetChan <- newTaskChan
		newTask := <-newTaskChan.response

		if !isInteractive {
			fmt.Print("Worker #", id, " gets task: ")
			newTask.getTask().print()
		}

		/* Inserting task into machine */

		machine := newTask.getRandMachine()
		insertTaskChan := machine.getTaskInsertChan()
		machineSolveChan := make(chan task)
		ito := &insertTaskOp{
			newTask.getTask(),
			machineSolveChan,
		}

		solved := false
		if isPatient {
			for !solved {
				insertTaskChan <- ito
				solvedTask = <-machineSolveChan
				if solvedTask.getResult() == nil {
					if !isInteractive {
						fmt.Println("Worker #", id, "reports broken machine #", machine.getId(), "to service")
					}
					brokenMachineChan <- machine.getId()
					machine = newTask.getRandMachine()
					insertTaskChan = machine.getTaskInsertChan()
				} else {
					solved = true
				}
			}
		} else {
			for !solved {
				select {
				case insertTaskChan <- ito:
					solvedTask = <-machineSolveChan
					if solvedTask.getResult() == nil {
						if !isInteractive {
							fmt.Println("Worker #", id, "reports broken machine #", machine.getId(), "to service")
						}
						brokenMachineChan <- machine.getId()
						machine = newTask.getRandMachine()
						insertTaskChan = machine.getTaskInsertChan()
					} else {
						solved = true
					}
				case <-time.After(companyConstants.ImpatientWorkerWaitTime * time.Second):
					machine = newTask.getRandMachine()
					insertTaskChan = machine.getTaskInsertChan()
				}
			}
		}

		/* Creating product */

		solvedTaskCounter++
		product := *solvedTask.getResult()

		if !isInteractive {
			fmt.Println("Worker #", id, "creates product:", product)
		}

		productPutChan <- product

		if !isInteractive {
			fmt.Println("Worker #", id, "puts product:", product, "in magazine")
		}
	}
}

func printWorkerStatistics(workerId int, isPatient bool, solvedTaskCount int) {
	fmt.Println()
	fmt.Println("Worker #", workerId)
	fmt.Println("Patient: ", isPatient)
	fmt.Println("# of solved task: ", solvedTaskCount)
	fmt.Println()
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
