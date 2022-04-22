package main

import "github.com/digitalcircle-com-br/local-agent/lib/server"

func main() {
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
