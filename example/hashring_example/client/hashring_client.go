package main

import (
	"fmt"
	"log"

	"stathat.com/c/consistent"
)

func main() {
	c := consistent.New()
	c.NumberOfReplicas = 20
	c.Add("localhsot:9999")
	c.Add("192.168.0.1:2233")
	c.Add("localhost:4455")
	users := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "x", "y", "z"}
	for _, u := range users {
		server, err := c.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}

	c.Remove("localhost:4455")
	fmt.pri
	for _, u := range users {
		server, err := c.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}
}
