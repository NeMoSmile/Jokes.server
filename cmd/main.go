package main

import (
	s "github.com/NeMoSmile/Jokes.server.git/internal"
)

var IP = "149.154.71.182"
var PORT = "8081"

var host = IP + ":" + PORT

func main() {
	s.Start(host)
}
