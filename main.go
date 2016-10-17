package main

import (
	"fmt"
	"log"

	"github.com/Tradeshift/powerdns/client"
)

func main() {
	c := client.NewClient("http://localhost:8081", "changeme")
	err := c.AddZone("dk.cvr.tradeshift.eu.", []string{"ns.amazon.com.", "ns1.amazon.com."})
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("OK..")
	}
	fmt.Println("-------------------------------------")
	fmt.Println("About to create a new record")
	err = c.AddRecord("12345678.dk.cvr.tradeshift.eu.", "A", "127.0.0.1", 30, "dk.cvr.tradeshift.eu.")

	if err != nil {
		log.Println("TIS", err)
	} else {
		fmt.Println("Record created OK..")
	}

}
