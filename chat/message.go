package main

import (
	"time"
)

//mesage represents a single messsage
type message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
}
