package main

import "github.com/tonsV2/event-rooster-api/di"

func main() {
	s := di.BuildServer()
	s.Run()
}
