package main

import (
	"company/companyConstants"
	"math/rand"
	"time"
)

type machine interface {
	run()
	getId() int
	getTaskInsertChan() chan *insertTaskOp
	getBackDoor() chan bool
}

type insertTaskOp struct {
	task       task
	solvedChan chan task
}

/*
	--- Addition machine ---
*/

type additionMachine struct {
	id             int
	taskInsertChan chan *insertTaskOp
	backDoor       chan bool
}

func (machine additionMachine) getId() int {
	return machine.id
}

func (machine additionMachine) getTaskInsertChan() chan *insertTaskOp {
	return machine.taskInsertChan
}

func (machine additionMachine) getBackDoor() chan bool {
	return machine.backDoor
}

func (machine additionMachine) run() {
	sleepTime := time.Duration(companyConstants.AdditionTime) * time.Second
	isBroken := false

	for {
		select {
		case insertedTask := <-machine.taskInsertChan:
			time.Sleep(sleepTime)
			if !isBroken {
				if rand.Float32() > companyConstants.MachineReliability {
					isBroken = true
				}
			}
			insertedTask.task = insertedTask.task.solve(isBroken)
			insertedTask.solvedChan <- insertedTask.task
		case <-machine.backDoor:
			isBroken = false
		}
	}
}

/*
	--- Multiplication machine ---
*/

type multiplicationMachine struct {
	id             int
	taskInsertChan chan *insertTaskOp
	backDoor       chan bool
}

func (machine multiplicationMachine) getId() int {
	return machine.id
}

func (machine multiplicationMachine) getTaskInsertChan() chan *insertTaskOp {
	return machine.taskInsertChan
}

func (machine multiplicationMachine) getBackDoor() chan bool {
	return machine.backDoor
}

func (machine multiplicationMachine) run() {
	sleepTime := time.Duration(companyConstants.MultiplicationTime) * time.Second
	isBroken := false

	for {
		select {
		case insertedTask := <-machine.taskInsertChan:
			time.Sleep(sleepTime)
			if !isBroken {
				if rand.Float32() > companyConstants.MachineReliability {
					isBroken = true
				}
			}
			insertedTask.task = insertedTask.task.solve(isBroken)
			insertedTask.solvedChan <- insertedTask.task
		case <-machine.backDoor:
			isBroken = false
		}
	}
}

/*
	--- Machine creating
*/

func createMachines() ([]additionMachine, []multiplicationMachine) {
	multiplicationMachines := make([]multiplicationMachine, 0, companyConstants.MultiplicationMachinesCount)
	additionMachines := make([]additionMachine, 0, companyConstants.AdditionMachinesCount)

	machineId := 0

	for i := 0; i < companyConstants.AdditionMachinesCount; i++ {
		machine := additionMachine{
			id:             machineId,
			taskInsertChan: make(chan *insertTaskOp),
			backDoor:       make(chan bool),
		}
		additionMachines = append(additionMachines, machine)
		go machine.run()
		machineId++
	}

	for i := 0; i < companyConstants.MultiplicationMachinesCount; i++ {
		machine := multiplicationMachine{
			id:             machineId,
			taskInsertChan: make(chan *insertTaskOp),
			backDoor:       make(chan bool),
		}
		multiplicationMachines = append(multiplicationMachines, machine)
		go machine.run()
		machineId++
	}

	return additionMachines, multiplicationMachines
}
