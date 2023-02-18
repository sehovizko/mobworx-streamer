package main

import (
	"fmt"
	"hls.streaming.com/streaming"
)

func main() {
	streaming.InitConfig()

	conf := streaming.LoadedConfig()
	fmt.Println(conf)

}
