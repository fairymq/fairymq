package main

import (
	"github.com/hashicorp/memberlist"
)

func (fairyMQ *FairyMQ) SetupMemberListCluster() error {
	config := memberlist.DefaultLocalConfig()
	// TODO: Setup Delegate
	config.Delegate = &Delegate{}
	// TODO: Setup EventDelegate
	config.Events = &EventDelegate{}

	list, err := memberlist.Create(config)
	if err != nil {
		return err
	}

	// If JoinAddr is provided, then join the cluster
	// TODO: Replace this with check for JoinAddr in config
	if true {
		_, err := list.Join([]string{"JoinAddr here"})
		if err != nil {
			return err
		}
	}

	return nil
}
