package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Invalid arguments count. Must be 2 or greater")
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(RunCmd(os.Args[2:], env, os.Stdin, os.Stdout, os.Stderr))
}
