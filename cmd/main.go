package main

import (
	s "github.com/NeMoSmile/Jokes.server.git/internal"
)

var host = ":8081"

func main() {
	s.Start(host)
}
