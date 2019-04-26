package main

import "math/rand"

type taskMachineAdapter interface {
	getTask() task
	getRandMachine() machine
}

type multiplicationAdapter struct {
	task     multiplicationTask
	machines []multiplicationMachine
}

type additionAdapter struct {
	task     additionTask
	machines []additionMachine
}

func (ma multiplicationAdapter) getTask() task {
	return ma.task
}

func (ma multiplicationAdapter) getRandMachine() machine {
	idx := rand.Intn(len(ma.machines))
	return ma.machines[idx]
}

func (aa additionAdapter) getTask() task {
	return aa.task
}

func (aa additionAdapter) getRandMachine() machine {
	idx := rand.Intn(len(aa.machines))
	return aa.machines[idx]
}
