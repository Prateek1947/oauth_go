package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	data, err := ioutil.ReadFile("/home/prateek/learn/Learn.py")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	ioutil.WriteFile("io/content.yaml", data, 0777)
	file, err := os.OpenFile("io/content.yaml", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	file.Write(data)
}
