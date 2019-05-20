package main

import (
	"company/companyConstants"
	"fmt"
)

type task interface {
	solve(isBroken bool) task
	getResult() *int
	print()
}

type additionTask struct {
	arg1   int
	op     string
	arg2   int
	result *int
}

func (t additionTask) solve(isBroken bool) task {
	if isBroken {
		t.result = nil
	} else {
		t.result = new(int)
		*t.result = t.arg1 + t.arg2
	}
	return t
}

func (t additionTask) getResult() *int {
	return t.result
}

func (t additionTask) print() {
	fmt.Println(t.arg1, t.op, t.arg2)
}

type multiplicationTask struct {
	arg1   int
	op     string
	arg2   int
	result *int
}

func (t multiplicationTask) solve(isBroken bool) task {
	if isBroken {
		t.result = nil
	} else {
		t.result = new(int)
		*t.result = t.arg1 * t.arg2
	}
	return t
}

func (t multiplicationTask) getResult() *int {
	return t.result
}

func (t multiplicationTask) print() {
	fmt.Println(t.arg1, t.op, t.arg2)
}

type subtractionTask struct {
	arg1   int
	op     string
	arg2   int
	result *int
}

func (t subtractionTask) solve(isBroken bool) task {
	if isBroken {
		t.result = nil
	} else {
		t.result = new(int)
		*t.result = t.arg1 - t.arg2
	}
	return t
}

func (t subtractionTask) getResult() *int {
	return t.result
}

func (t subtractionTask) print() {
	fmt.Println(t.arg1, t.op, t.arg2)
}

type getTaskOp struct {
	response chan taskMachineAdapter
}

func taskAddFilter(addChan chan taskMachineAdapter, taskList []taskMachineAdapter) chan taskMachineAdapter {
	if len(taskList) < cap(taskList) {
		return addChan
	}
	return nil
}

func taskGetFilter(getChan chan *getTaskOp, taskList []taskMachineAdapter) chan *getTaskOp {
	if len(taskList) > 0 {
		return getChan
	}
	return nil
}

func taskListServer(taskAddChan chan taskMachineAdapter, taskGetChan chan *getTaskOp, infoChan chan bool) {
	var taskList = make([]taskMachineAdapter, 0, companyConstants.SizeOfTaskList)

	for {
		select {
		case newTask := <-taskAddFilter(taskAddChan, taskList):
			taskList = append(taskList, newTask)
		case get := <-taskGetFilter(taskGetChan, taskList):
			get.response <- taskList[0]
			taskList = append(taskList[:0], taskList[1:]...)
		case <-infoChan:
			displayTaskList(taskList)
			infoChan <- true
		}
	}
}

func displayTaskList(taskList []taskMachineAdapter) {
	fmt.Println()
	fmt.Println("task list:")
	for _, task := range taskList {
		fmt.Println(task.getTask())
	}
	fmt.Println()
}
