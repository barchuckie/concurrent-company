package main

import (
	"bufio"
	"fmt"
	"os"
)

func interact(taskListInfoChan chan bool, magazineInfoChan chan bool, workersInfoChan []chan int, done chan bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("--- Company Simulator ---")
		fmt.Println("mag		<- display magazine")
		fmt.Println("tl		<- display task list")
		fmt.Println("w			<- display workers info")
		fmt.Println("quit	<- quit program")

		text, _ := reader.ReadString('\n')

		switch text {
		case "mag\n":
			magazineInfoChan <- true
			<-magazineInfoChan
		case "tl\n":
			taskListInfoChan <- true
			<-taskListInfoChan
		case "w\n":
			for i := 0; i < len(workersInfoChan); i++ {
				workersInfoChan[i] <- i
				<-workersInfoChan[i]
			}
		case "quit\n":
			done <- true
			return
		}
	}
}
