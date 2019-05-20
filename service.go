package main

import (
	"company/companyConstants"
	"fmt"
	"time"
)

func service(brokenMachineChan chan int, repairedMachineChan chan int,
	workers []serviceWorker, isInteractive bool) {

	brokenMachines := make([]int, 0, companyConstants.AdditionMachinesCount+companyConstants.AdditionMachinesCount)
	workingMachines := make([]bool, companyConstants.AdditionMachinesCount+companyConstants.AdditionMachinesCount)

	for i := 0; i < cap(workingMachines); i++ {
		workingMachines[i] = true
	}

	for {
		select {
		case brokenMachine := <-brokenMachineChan:
			if !isInteractive {
				fmt.Println("Service receives report about broken machine #", brokenMachine)
			}
			workingMachines[brokenMachine] = false
			add := true
			for i := range brokenMachines {
				if brokenMachines[i] == brokenMachine {
					add = false
					break
				}
			}
			if add && !workingMachines[brokenMachine] {
				brokenMachines = append(brokenMachines, brokenMachine)
			}
		case repairedMachine := <-repairedMachineChan:
			if !isInteractive {
				fmt.Println("Service repaired machine #", repairedMachine)
			}
			workingMachines[repairedMachine] = true
		default:
			for i := range workers {
				if len(brokenMachines) == 0 {
					break
				}
				select {
				case workers[i].brokenMachineId <- brokenMachines[0]:
					if !isInteractive {
						fmt.Println("Service sends worker #", i, "to broken machine #", brokenMachines[0])
					}
					brokenMachines = append(brokenMachines[:0], brokenMachines[1:]...)
					break
				default:
					break
				}
			}
		}
	}
}

type serviceWorker struct {
	id              int
	brokenMachineId chan int
	isBusy          bool
}

func (sw serviceWorker) work(machines []machine, repairedMachineChan chan<- int, isInteractive bool) {

	for {
		brokenMachine := <-sw.brokenMachineId
		sw.isBusy = true
		if !isInteractive {
			fmt.Println("Service Worker #", sw.id, "starts repairing machine #", brokenMachine)
		}
		backDoor := machines[brokenMachine].getBackDoor()
		time.Sleep(companyConstants.AccessTime * time.Second)
		backDoor <- true
		repairedMachineChan <- brokenMachine
		sw.isBusy = false
	}

}

func createServiceWorkers() []serviceWorker {
	serviceWorkers := make([]serviceWorker, 0, companyConstants.ServiceWorkersCount)

	for i := 0; i < companyConstants.ServiceWorkersCount; i++ {
		sw := serviceWorker{
			id:              i,
			brokenMachineId: make(chan int),
			isBusy:          false,
		}
		serviceWorkers = append(serviceWorkers, sw)
	}

	return serviceWorkers
}
