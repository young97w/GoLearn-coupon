package main

import (
	"coupon/internal"
	"encoding/json"
	"fmt"
)

func main() {
	mrsConf, err := json.Marshal(internal.AppConf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(mrsConf))
}
