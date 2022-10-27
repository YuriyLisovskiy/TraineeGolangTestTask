package main

import (
	"log"

	"TraineeGolangTestTask/cli"
)

func main() {
	if err := cli.RunCLI(); err != nil {
		log.Fatal(err)
	}
}
