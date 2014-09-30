/*
   Copyright 2014 Franc[e]sco (lolisamurai@tfwno.gf)
   This file is part of kagami.
   kagami is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   kagami is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with kagami. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
	"fmt"
	"net"
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/config"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/kagami/worldserver/channels"
	"github.com/Francesco149/kagami/worldserver/status"
	"github.com/Francesco149/maplelib"
)

// HandleLogin handles packets exchanged between the worldserver and the loginserver
func HandleLogin(con *common.InterserverClient, p maplelib.Packet) (handled bool, err error) {
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	switch header {
	case interserver.IOWorldConnect:
		return handleWorldConnect(con, it)
	}

	return false, nil
}

// handleWorldConnect handles a world connect packet from the login server, which tells the worldserver
// which world it will handle and provide the world configuration
func handleWorldConnect(con *common.InterserverClient, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	worldid, err := it.Decode1s()
	if err != nil {
		return
	}

	if worldid == -1 {
		fmt.Println("No worlds to handle!")
		return
	}

	port, err := it.Decode2s()
	conf, err := config.DecodeWorldConf(&it)
	if err != nil {
		return
	}

	handled = true
	fmt.Println("Handling world", worldid)
	status.Lock()
	defer status.Unlock()
	status.SetConf(conf)
	status.SetPort(port)
	status.SetLoginConn(con)
	status.SetWorldId(worldid)

	// TODO: check if I need to store the loginserver's external ip address

	// accept interserver chan connections in a separate thread
	go common.Accept("chan", status.Port(),
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*channels.Connection)
			if !ok {
				return false, errors.New("Channel handler failed type assertion")
			}
			return HandleChan(scon, p)
		},
		func(con net.Conn) common.Connection {
			return channels.NewConnection(con, consts.InterServerPassword)
		},
		func(con common.Connection) {
			scon, ok := con.(*channels.Connection)
			if !ok {
				panic(errors.New("Channel handler failed type assertion on disconnect"))
			}
			deletechanid := scon.ChannelId()

			if deletechanid == -1 {
				return
			}

			fmt.Println("Removing channel", deletechanid)
			status.Lock()
			channels.Lock()
			defer status.Unlock()
			defer channels.Unlock()
			if status.LoginConn() != nil {
				status.LoginConn().SendPacket(interserver.RemoveChannel(deletechanid))
			}

			// TODO: disconnect players
			channels.Remove(deletechanid)
		})

	fmt.Println("World server is running!")
	return
}
