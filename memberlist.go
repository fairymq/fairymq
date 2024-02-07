package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
)

func (fairyMQ *FairyMQ) SetupMemberListCluster() error {
	config := memberlist.DefaultWANConfig()
	config.Name = fmt.Sprintf("%s:%d", fairyMQ.config.BindAddress, fairyMQ.config.MemberlistPort)
	config.BindAddr = fairyMQ.config.BindAddress
	config.BindPort = int(fairyMQ.config.MemberlistPort)
	config.AdvertisePort = int(fairyMQ.config.MemberlistPort)
	config.Delegate = &Delegate{
		fairyMQ: fairyMQ,
	}

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
