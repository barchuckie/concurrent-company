package main

import (
	"bufio"
	"fmt"
	"os"
)

func interact(taskListInfoChan chan bool, magazineInfoChan chan bool, done chan bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("--- Company Simulator ---")
		fmt.Println("mag		<- display magazine")
		fmt.Println("tl		<- display task list")
		fmt.Println("quit	<- quit program")

		text, _ := reader.ReadString('\n')

		switch text {
		case "mag\n":
			magazineInfoChan <- true
		case "tl\n":
			taskListInfoChan <- true
		case "quit\n":
			done <- true
			return
		}
	}
}
