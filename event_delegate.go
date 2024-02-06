package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
)

type EventDelegate struct{}

func (event *EventDelegate) NotifyJoin(node *memberlist.Node) {
	// TODO: Handle node joining the cluster
	fmt.Printf("A new member has joined: %+v\n", node)
}

func (event *EventDelegate) NotifyLeave(node *memberlist.Node) {
	// TODO: Handle node leaving the cluster
}

func (event *EventDelegate) NotifyUpdate(node *memberlist.Node) {
	// TODO: Handle updating of the node meta data
}
