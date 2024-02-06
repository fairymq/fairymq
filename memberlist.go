package main

import (
	"github.com/hashicorp/memberlist"
)

func (fairyMQ *FairyMQ) SetupMemberListCluster() error {
	config := memberlist.DefaultLocalConfig()
	config.BindAddr = fairyMQ.config.BindAddress
	config.BindPort = int(fairyMQ.config.MemberlistPort)
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
