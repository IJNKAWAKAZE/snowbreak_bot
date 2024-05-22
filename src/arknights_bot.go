package main

import (
	"log"
	"snowbreak_bot/cmd"
)

func main() {
	cmd.Execute()
}

// 设置日志格式
func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}
