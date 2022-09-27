package main

import (
	"account/internal"
	"fmt"
)

func main() {
	//fmt.Println(internal.AppConf)
	fmt.Println(internal.GRPC == "grpc")
}
