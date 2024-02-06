package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"math/rand"
)

func (fairyMQ *FairyMQ) SetupMemberListCluster() error {
	config := memberlist.DefaultWANConfig()
	// TODO: Set a sensible unique name for this node
	config.Name = fmt.Sprintf("%d", rand.Int())
	config.BindAddr = fairyMQ.config.BindAddress
	config.BindPort = int(fairyMQ.config.MemberlistPort)
	config.AdvertisePort = int(fairyMQ.config.MemberlistPort)
	config.Delegate = &Delegate{}
	config.Events = &EventDelegate{}

	list, err := memberlist.Create(config)
	if err != nil {
		return err
	}

	if len(fairyMQ.config.JoinAddresses) > 0 {
		_, err = list.Join(fairyMQ.config.JoinAddresses)
		if err != nil {
			return err
		}
	}

	return nil
}
