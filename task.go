package main

import (
	"company/companyConstants"
	"fmt"
)

type AdditionTask struct {
	arg1 int
	op   string
	arg2 int
}

type MultiplicationTask struct {
	arg1 int
	op   string
	arg2 int
}

type SubtractionTask struct {
	arg1 int
	op   string
	arg2 int
}

type Task interface {
	solve() int
}

func (t AdditionTask) solve() int {
	return t.arg1 + t.arg2
}

func (t MultiplicationTask) solve() int {
	return t.arg1 * t.arg2
}

func (t SubtractionTask) solve() int {
	return t.arg1 - t.arg2
}

type getTaskOp struct {
	response chan Task
}

func taskAddFilter(addChan chan Task, taskList []Task) chan Task {
	if len(taskList) < cap(taskList) {
		return addChan
	}
	return nil
}

func taskGetFilter(getChan chan *getTaskOp, taskList []Task) chan *getTaskOp {
	if len(taskList) > 0 {
		return getChan
	}
	return nil
}

func taskListServer(taskAddChan chan Task, taskGetChan chan *getTaskOp, infoChan chan bool) {
	var taskList = make([]Task, 0, companyConstants.SizeOfTaskList)

	for {
		select {
		case newTask := <-taskAddFilter(taskAddChan, taskList):
			taskList = append(taskList, newTask)
		case get := <-taskGetFilter(taskGetChan, taskList):
			get.response <- taskList[0]
			taskList = append(taskList[:0], taskList[1:]...)
		case <-infoChan:
			displayTaskList(taskList)
		}
	}
}

func displayTaskList(taskList []Task) {
	fmt.Println()
	fmt.Println("Task list:")
	for _, task := range taskList {
		fmt.Println(task)
	}
	fmt.Println()
}
