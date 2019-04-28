package main

import (
	"company/companyConstants"
	"time"
)

type machine interface {
	run()
	getTaskInsertChan() chan *insertTaskOp
}

type multiplicationMachine struct {
	id             int
	taskInsertChan chan *insertTaskOp
}

type additionMachine struct {
	id             int
	taskInsertChan chan *insertTaskOp
}

type insertTaskOp struct {
	task       task
	solvedChan chan task
}

func (machine multiplicationMachine) getTaskInsertChan() chan *insertTaskOp {
	return machine.taskInsertChan
}

func (machine multiplicationMachine) run() {
	sleepTime := time.Duration(companyConstants.MultiplicationTime) * time.Second

	for {
		insertedTask := <-machine.taskInsertChan
		time.Sleep(sleepTime)
		insertedTask.task = insertedTask.task.solve()
		insertedTask.solvedChan <- insertedTask.task
	}
}

func (machine additionMachine) getTaskInsertChan() chan *insertTaskOp {
	return machine.taskInsertChan
}

func (machine additionMachine) run() {
	sleepTime := time.Duration(companyConstants.AdditionTime) * time.Second

	for {
		insertedTask := <-machine.taskInsertChan
		time.Sleep(sleepTime)
		insertedTask.task = insertedTask.task.solve()
		insertedTask.solvedChan <- insertedTask.task
	}
}
