package main

import "intent/cmd"

func main() {
	go cmd.StartWorker()
	cmd.StartServer()
}
