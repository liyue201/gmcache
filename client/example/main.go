package main

import (
	"github.com/liyue201/gmcache/client"
	"log"
)

func main() {
	c := client.NewClient("127.0.0.1:8002")
	if err := c.Connect(); err != nil {
		log.Println(err)
		return
	}
	if err := c.Set("aaa", []byte("bbbbbbb"), 10); err != nil {
		log.Println(err)
		return
	}
	val, err := c.Get("aaa")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(val))

	err = c.Delete("aaa")
	if err != nil {
		log.Println(err)
		return
	}
	c.Disconnect()

}
