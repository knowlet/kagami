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
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/kagami/common/player"
	"github.com/Francesco149/kagami/worldserver/channels"
	"github.com/Francesco149/kagami/worldserver/players"
	"github.com/Francesco149/kagami/worldserver/status"
	"github.com/Francesco149/maplelib"
)

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
			status.Lock()
			channels.Lock()
			defer channels.Unlock()
			defer status.Unlock()

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
	case interserver.IOSyncWorldPerformChangeChannel:
		return syncPerformChangeChannel(con, it)

	case interserver.IOSyncWorldLoadCharacter:
		return syncLoadCharacter(con, it)
	}

	return false, nil
}

// syncPerformChangeChannel updates a ccing player's channel
func syncPerformChangeChannel(con *channels.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	id, err := it.Decode4s()
	if err != nil {
		return
	}

	players.Lock()
	defer players.Unlock()

	player := players.Get(id)
	curchan := channels.Get(player.Channel())
	if curchan == nil {
		return
	}

	newchanid := players.GetPendingCC(id)
	newchan := channels.Get(newchanid)
	ip := make([]byte, 4)
	port := int16(-1)

	if newchan != nil {
		cconn := newchan.Conn().Conn()
		if cconn != nil {
			ip = common.RemoteAddrToBytes(cconn.RemoteAddr().String())
			port = newchan.Port()
		}
	}

	curchan.Conn().SendPacket(interserver.SyncChannelPerformChangeChannel(id, newchanid, ip, port))
	players.RemovePendingCC(id)
	handled = err == nil
	return
}

// syncLoadCharacter adds a newly connected player's data to the player pool
func syncLoadCharacter(con *channels.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	data, err := player.DecodeData(&it)
	if err != nil {
		return
	}

	channels.Lock()
	players.Lock()
	defer players.Unlock()
	defer channels.Unlock()

	player := players.Get(data.CharId())
	*player = *data
	player.SetInitialized(true)
	player.SetTransferring(false)

	syncpacket, err := interserver.SyncChannelUpdatePlayer(player)
	if err != nil {
		return
	}
	err = channels.SendToAllChannels(syncpacket)
	handled = err == nil
	return
}
