package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"time"
)

func (fairyMQ *FairyMQ) SetupMemberListCluster() (func() error, error) {
	config := memberlist.DefaultWANConfig()
	config.Name = fmt.Sprintf("%s:%d", fairyMQ.Config.BindAddress, fairyMQ.Config.MemberlistPort)
	config.ProtocolVersion = memberlist.ProtocolVersionMax
	config.BindAddr = fairyMQ.Config.BindAddress
	config.BindPort = int(fairyMQ.Config.MemberlistPort)
	config.AdvertisePort = int(fairyMQ.Config.MemberlistPort)
	config.EnableCompression = true
	config.Delegate = &Delegate{
		fairyMQ: fairyMQ,
	}

	list, err := memberlist.Create(config)

	if err != nil {
		return nil, err
	}

	if len(fairyMQ.Config.JoinAddresses) > 0 {
		_, err = list.Join(fairyMQ.Config.JoinAddresses)
		if err != nil {
			return nil, err
		}
	}

	shutdownFunc := func() error {
		if err := list.Leave(200 * time.Millisecond); err != nil {
			return fmt.Errorf("memberlist shutdown - leave: %+v", err)
		}
		if err = list.Shutdown(); err != nil {
			return fmt.Errorf("memberlist shutdown - shutdown: %+v", err)
		}
		return nil
	}

	return shutdownFunc, nil
}
