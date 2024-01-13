/*
* fairyMQ
* Core Unit Tests
* ******************************************************************
* Originally authored by Alex Gaetano Padula
* Copyright (C) fairyMQ
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"context"
	"net"
	"os"
	"sync"
	"testing"
)

func TestFairyMQ_GenerateQueueKeypair(t *testing.T) {
	type fields struct {
		UDPAddr       *net.UDPAddr
		Conn          *net.UDPConn
		Wg            *sync.WaitGroup
		SignalChannel chan os.Signal
		Queues        map[string][]Message
		ContextCancel context.CancelFunc
		Context       context.Context
	}
	//type args struct {
	//	queue string
	//}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{name: "test", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fairyMQ := &FairyMQ{
				UDPAddr:       tt.fields.UDPAddr,
				Conn:          tt.fields.Conn,
				Wg:            tt.fields.Wg,
				SignalChannel: tt.fields.SignalChannel,
				Queues:        tt.fields.Queues,
				ContextCancel: tt.fields.ContextCancel,
				Context:       tt.fields.Context,
			}
			if err := fairyMQ.GenerateQueueKeypair("test"); (err != nil) != tt.wantErr {
				t.Errorf("GenerateQueueKeypair() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFairyMQ_SignalListener(t *testing.T) {
	type fields struct {
		UDPAddr       *net.UDPAddr
		Conn          *net.UDPConn
		Wg            *sync.WaitGroup
		SignalChannel chan os.Signal
		Queues        map[string][]Message
		ContextCancel context.CancelFunc
		Context       context.Context
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// todo
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fairyMQ := &FairyMQ{
				UDPAddr:       tt.fields.UDPAddr,
				Conn:          tt.fields.Conn,
				Wg:            tt.fields.Wg,
				SignalChannel: tt.fields.SignalChannel,
				Queues:        tt.fields.Queues,
				ContextCancel: tt.fields.ContextCancel,
				Context:       tt.fields.Context,
			}
			fairyMQ.SignalListener()
		})
	}
}

func TestFairyMQ_StartUDPListener(t *testing.T) {
	type fields struct {
		UDPAddr       *net.UDPAddr
		Conn          *net.UDPConn
		Wg            *sync.WaitGroup
		SignalChannel chan os.Signal
		Queues        map[string][]Message
		ContextCancel context.CancelFunc
		Context       context.Context
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// todo
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fairyMQ := &FairyMQ{
				UDPAddr:       tt.fields.UDPAddr,
				Conn:          tt.fields.Conn,
				Wg:            tt.fields.Wg,
				SignalChannel: tt.fields.SignalChannel,
				Queues:        tt.fields.Queues,
				ContextCancel: tt.fields.ContextCancel,
				Context:       tt.fields.Context,
			}
			fairyMQ.StartUDPListener()
		})
	}
}
