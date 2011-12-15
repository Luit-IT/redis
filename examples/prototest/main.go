package main

import (
	"log"
	"luit.it/redis.git/redis/proto"
)

func main() {
	conn, err := proto.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal("Error dialing: ", err)
	}
	o, err := conn.Command("HVALS", "a")
	if err != nil {
		log.Fatal("Error in conn.Command(): ", err)
	}
	l, err := proto.ObjectList(o)
	if err != nil {
		log.Fatal("Error in ObjectList(): ", err)
	}
	for _, li := range l {
		s, err := proto.ObjectString(li)
		if err != nil {
			log.Fatal("Error in ObjectList(): ", err)
		}
		log.Print("List item: ", s)
	}
}
