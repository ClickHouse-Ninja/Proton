package main

import (
	"log"
	"time"

	"github.com/mkevac/gopinba"
)

func main() {
	pc, err := gopinba.NewClient("127.0.0.1:30002")
	if err != nil {
		log.Fatalf("NewClient() returned error: %v", err)
	}

	req := gopinba.Request{
		Tags: map[string]string{
			"A": "B",
			"C": "D",
		},
	}

	for i := 0; i < 50000000; i++ {
		req.Hostname = "hostname"
		req.ServerName = "servername"
		req.ScriptName = "scriptname"
		req.Schema = "https"
		req.Status = 200
		req.RequestCount = 1
		req.RequestTime = 145987 * time.Microsecond
		req.DocumentSize = 1024
		err = pc.SendRequest(&req)
		if err != nil {
			log.Fatalf("SendRequest() returned error: %v", err)
		}
	}
}
