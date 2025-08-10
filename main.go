package main

import (
	"os"

	"github.com/genkaieng/nicopush/cmd"
)

func main() {
	if len(os.Args) < 2 {
		println("サブコマンドを指定してください: subscribe, genkeys")
		os.Exit(1)
	}

	var result int
	switch os.Args[1] {
	case "subscribe":
		result = cmd.Subscribe(os.Args[2:])
	case "genkeys":
		result = cmd.Genkeys(os.Args[2:])
	default:
		println("不明なサブコマンド: " + os.Args[1])
		result = 1
	}
	os.Exit(result)
}
