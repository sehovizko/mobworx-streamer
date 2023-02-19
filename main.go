package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"hls.streaming.com/streaming/helpers"
	"log"
)

func main() {
	mySession := session.Must(session.NewSession())
	sc := helpers.NewStreamerConnection(mySession)
	log.Println(sc)
}
