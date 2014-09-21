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
	"math/rand"
	"net"
	"time"
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/kagami/common/config"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/kagami/worldserver/status"
	"github.com/Francesco149/kagami/worldserver/channels"
	"github.com/Francesco149/maplelib"
)

// TODO: everything, this is just a temporary main that will get reoganized into multiple files as I add stuff

func makeChannelPort(chanid int8) int16 {
        return status.Port() + int16(chanid) + 1        
}

// channelConnect returns a packet that respons to a channelserver connection request
func channelConnect(channelId int8, port int16) (p maplelib.Packet) {
        p = packets.NewEncryptedPacket(interserver.IOChannelConnect)
        p.Encode1s(channelId)
        p.Encode2s(port)
        status.Conf().Encode(&p)
        return
}

// HandleChan handles packets exchanged between the worldserver and the channelserver
func HandleChan(con *channels.Connection, p maplelib.Packet) (handled bool, err error) {
        handled = false
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	// check auth
	if !con.Authenticated() {
		if header != interserver.IOAuth {
			err = errors.New(fmt.Sprintf("Tried to send %v without being authenticated", p))
			return
		}

		var servertype byte = 255
		servertype, err = con.CheckAuth(it)
		if err != nil {
			return
		}

		switch servertype {
		case interserver.ChannelServer: // we're only accepting channel serv connections here
			available := channels.GetFirstAvailableId()
			con.SetChannelId(available)
			
			if available == -1 {
			        con.SendPacket(channelConnect(-1, 0))
			        err = errors.New("No more channels available to assign.")
			        return
			}
			
			chanport := makeChannelPort(available)
			// TODO: get external ip
			ipbytes := common.RemoteAddrToBytes(con.Conn().RemoteAddr().String())
			channels.Add(con, available, ipbytes, chanport)
			con.SendPacket(channelConnect(available, chanport))
			// TODO: send sync channel connect
			status.LoginConn().SendPacket(interserver.RegisterChannel(available, ipbytes, chanport))
		default:
			err = errors.New("Unknown server type")
		}

		return true, nil
	}

	switch header {
	// TODO
	}

	return false, nil
}

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
	status.SetConf(conf)
	status.SetPort(port)
	status.SetLoginConn(con)

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
			if status.LoginConn() != nil {
				status.LoginConn().SendPacket(interserver.RemoveChannel(deletechanid))
			}

			// TODO: disconnect players
			channels.Remove(deletechanid)
		})

	fmt.Println("World server is running!")
	return
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Kagami Pre-Alpha")
	fmt.Println("Initializing WorldServer...")

	// TODO: accept channelserver connections
	common.Connect("loginserver", fmt.Sprintf("%s:%d", consts.LoginIp, consts.LoginInterserverPort),
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*common.InterserverClient)
			if !ok {
				panic(errors.New("Login handler failed type assertion"))
			}
			return HandleLogin(scon, p)
		},
		func(con net.Conn) common.Connection {
			return common.NewInterserverClient(con, consts.InterServerPassword, interserver.WorldServer)
		})
}
