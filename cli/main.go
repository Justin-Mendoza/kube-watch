package main

import (
	"fmt/log"
	"net/http"
	"log"
	"os/exec"
	"os"
	"time"
	"context"
	"bufio"
	"fmt"
)

func main(){
	if len(os.Args) < 2{
		fmt.Println("Usage: k8s-cli [deploy|destroy|monitor]")
		return
	}

	switch os.Args[1]{
	case "deploy":
		deploy()
	case "destroy":
		destroy()
	case "monitor":
		monitor()
	default:
		fmt.Println("Unknown command")
	}
}