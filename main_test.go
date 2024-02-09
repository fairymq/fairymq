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
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestFairyMQ_SignalListener(t *testing.T) {
	type fields struct {
		UDPAddr                *net.UDPAddr
		Conn                   *net.UDPConn
		Wg                     *sync.WaitGroup
		SignalChannel          chan os.Signal
		Queues                 map[string]*Queue
		ContextCancel          context.CancelFunc
		Context                context.Context
		Config                 Config
		MemberlistShutdownFunc func() error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: false,
			fields: fields{
				Wg:                     &sync.WaitGroup{},
				SignalChannel:          make(chan os.Signal),
				MemberlistShutdownFunc: func() error { return nil },
			},
		},
	}
	for _, tt := range tests {
		var err error
		t.Run(tt.name, func(t *testing.T) {
			fairyMQ := &FairyMQ{
				UDPAddr:                tt.fields.UDPAddr,
				Conn:                   tt.fields.Conn,
				Wg:                     tt.fields.Wg,
				SignalChannel:          tt.fields.SignalChannel,
				Queues:                 tt.fields.Queues,
				ContextCancel:          tt.fields.ContextCancel,
				Context:                tt.fields.Context,
				Config:                 tt.fields.Config,
				MemberlistShutdownFunc: tt.fields.MemberlistShutdownFunc,
			}

			fairyMQ.UDPAddr, err = net.ResolveUDPAddr("udp", "0.0.0.0:5991")
			if err != nil {
				t.Errorf("TestFairyMQ_SignalListener() error = %v", err)
			}

			// Start listening for UDP packages on the given address
			fairyMQ.Conn, err = net.ListenUDP("udp", fairyMQ.UDPAddr)
			if err != nil {
				t.Errorf("TestFairyMQ_SignalListener() error = %v", err)
			}

			fairyMQ.Context, fairyMQ.ContextCancel = context.WithCancel(context.Background())

			fairyMQ.Wg.Add(1)
			go fairyMQ.SignalListener()

			fairyMQ.Wg.Add(1)
			go func() {
				defer fairyMQ.Wg.Done()
				tt.fields.SignalChannel <- os.Interrupt
			}()

			fairyMQ.Wg.Wait()
		})
	}
}

func TestFairyMQ_StartUDPListener(t *testing.T) {
	type fields struct {
		UDPAddr                *net.UDPAddr
		Conn                   *net.UDPConn
		Wg                     *sync.WaitGroup
		SignalChannel          chan os.Signal
		Queues                 map[string]*Queue
		ContextCancel          context.CancelFunc
		Context                context.Context
		Config                 Config
		MemberlistShutdownFunc func() error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: false,
			fields: fields{
				Wg:            &sync.WaitGroup{},
				SignalChannel: make(chan os.Signal),
				Config: Config{
					BindAddress: "0.0.0.0",
					BindPort:    5991,
				},
				MemberlistShutdownFunc: func() error { return nil },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fairyMQ := &FairyMQ{
				UDPAddr:                tt.fields.UDPAddr,
				Conn:                   tt.fields.Conn,
				Wg:                     tt.fields.Wg,
				SignalChannel:          tt.fields.SignalChannel,
				Queues:                 tt.fields.Queues,
				ContextCancel:          tt.fields.ContextCancel,
				Context:                tt.fields.Context,
				Config:                 tt.fields.Config,
				MemberlistShutdownFunc: tt.fields.MemberlistShutdownFunc,
			}

			fairyMQ.Wg.Add(1)
			go fairyMQ.StartUDPListener()

			fairyMQ.Context, fairyMQ.ContextCancel = context.WithCancel(context.Background())

			fairyMQ.Wg.Add(1)
			go func() {
				defer fairyMQ.Wg.Done()
				udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", fairyMQ.Config.BindAddress, fairyMQ.Config.BindPort))
				if err != nil {
					t.Errorf("TestFairyMQ_StartUDPListener() error = %v", err)
				}

				conn, err := net.DialUDP("udp", nil, udpAddr)
				if err != nil {
					t.Errorf("TestFairyMQ_StartUDPListener() error = %v", err)
				}

				_, err = conn.Write([]byte("testing, 1, 2, 3\n"))
				if err != nil {
					t.Errorf("TestFairyMQ_StartUDPListener() error = %v", err)
				}

				data, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					t.Errorf("TestFairyMQ_StartUDPListener() error = %v", err)
				}

				if strings.HasPrefix(data, "NACK") {
					fairyMQ.ContextCancel()
				} else {
					t.Errorf("TestFairyMQ_StartUDPListener() error = incorrect response.  expecting NACK")
				}
			}()

			fairyMQ.Wg.Wait()
		})
	}
}
