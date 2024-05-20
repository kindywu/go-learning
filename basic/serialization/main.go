package main

import (
	"encoding/json"
	"log"
	"time"
)

type Person struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Birthday time.Time `json:"birthday"`
}

func main() {
	var jack = Person{Id: 1, Name: "Jack", Birthday: time.Date(1979, 10, 19, 0, 0, 0, 0, time.Local)}
	str, err := json.Marshal(jack)
	if err != nil {
		log.Fatal("marshal error")
	}
	println(string(str))
	var jack2 Person
	err = json.Unmarshal([]byte(str), &jack2)
	if err != nil {
		log.Fatal("marshal error")
	}
	println(jack2.Name)
}
