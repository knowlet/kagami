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
	"github.com/Francesco149/kagami/channelserver/client"
	"github.com/Francesco149/kagami/channelserver/players"
	"github.com/Francesco149/kagami/channelserver/status"
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/config"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/maplelib"
)

// Handle handles inter-server packets exchanged between the channel server and the login/world server
func HandleInter(con *common.InterserverClient, p maplelib.Packet) (handled bool, err error) {
	handled = false
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	switch header {
	case interserver.IOLoginChannelConnect:
		return handleLoginChannelConnect(con, it)

	case interserver.IOChannelConnect:
		return handleChannelConnect(con, it)

	case interserver.IOPlayerJoiningChannel:
		return handlePlayerJoiningChannel(con, it)
	}

	return false, nil
}

// handleLoginChannelConnect handles a loginserver channel connect packet
// which tells the channel server which world it should connect to
func handleLoginChannelConnect(con *common.InterserverClient, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	worldId, err := it.Decode1s()
	if err != nil {
		return
	}

	// worldId will be -1 if there are no more worlds to handle
	if worldId == -1 {
		err = errors.New("No world server available")
		return
	}

	fmt.Println("Handling world", worldId, "'s channels")

	// decode ip as a byte array (this is the worldserver ip)
	ip := make([]byte, 4)

	for i := 0; i < 4; i++ {
		var tmp byte
		tmp, err = it.Decode1()
		if err != nil {
			return
		}
		ip[i] = tmp
	}

	// decode worldserver port
	port, err := it.Decode2s()
	if err != nil {
		return
	}

	// connect to worldserver
	go common.Connect("worldserver", fmt.Sprintf("%s:%d", common.BytesToIpString(ip), port),
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*common.InterserverClient)
			if !ok {
				return false, errors.New("Worldserver handler failed type assertion")
			}
			return HandleInter(scon, p)
		},
		func(con net.Conn) common.Connection {
			c := common.NewInterserverClient(con, consts.InterServerPassword, interserver.ChannelServer)
			status.Lock()
			defer status.Unlock()
			status.SetWorldConn(c)
			return c
		})

	handled = err == nil
	return
}

// handleChannelConnect handles a worldserver channel connect packet
// which tells the channel server which channel it will be handling
func handleChannelConnect(con *common.InterserverClient, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false

	chanid, err := it.Decode1s()
	if err != nil {
		return
	}

	if chanid == -1 {
		err = errors.New("No channel to handle")
		return
	}

	port, err := it.Decode2s()
	conf, err := config.DecodeWorldConf(&it)
	if err != nil {
		return
	}

	fmt.Println("Handling channel", chanid, "on port", port)
	status.Lock()
	status.SetChanId(chanid)
	status.SetPort(port)
	status.SetWorldConf(conf)
	status.Unlock()
	// TODO: set map unload time

	// accept client connections in a new thread
	go common.Accept("client", port,
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*client.Connection)
			if !ok {
				return false, errors.New("Client handler failed type assertion")
			}
			return Handle(scon, p)
		},
		func(con net.Conn) common.Connection {
			return client.NewConnection(con, false)
		},
		func(con common.Connection) {
			scon, ok := con.(*client.Connection)
			if !ok {
				panic(errors.New(common.MakeError("Client handler failed " +
					"type assertion on disconnect")))
			}
			status.WorldConn().SendPacket(interserver.SyncPlayerLeftChannel(status.ChanId()))
			err = scon.SetDBOnline(false)
			if err != nil {
				fmt.Println(common.MakeError(
					fmt.Sprint("Failed to disconnect ", scon.Name(), ": ", err)))
			}

			err = scon.Save()
			if err != nil {
				fmt.Println(common.MakeError(
					fmt.Sprint("Failed to save ", scon.Name(), ": ", err)))
			}
		})

	fmt.Println("Channel server is running!")

	handled = err == nil
	return
}

// handlePlayerJoiningChannel adds the pending player connection to the player pool
func handlePlayerJoiningChannel(con *common.InterserverClient, it maplelib.PacketIterator) (handled bool, err error) {
	charid, err := it.Decode4s()
	ip, err := it.DecodeBuffer()
	if err != nil {
		return
	}

	fmt.Println("Added pending player connection from", common.BytesToIpString(ip))
	players.AddPendingIp(charid, ip)
	handled = err == nil
	return
}
