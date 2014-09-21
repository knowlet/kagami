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

// Package worlds contains utilities to manage and send information about the
// worlds that are currently connected to the login server
package worlds

import (
	"errors"
	"fmt"
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/kagami/loginserver/client"
	"github.com/Francesco149/maplelib"
)

// TODO: make this thread safe?

var worlds = make(map[int8]*World)

// Add adds the given world to the world list
func Add(w *World) {
	worlds[w.Id()] = w
}

// Get returns the world associated with the given id
func Get(worldId int8) *World {
	return worlds[worldId]
}

// worldConnect returns a world connect inter-server packet
func worldConnect(w *World) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(interserver.IOWorldConnect)
	p.Encode1s(w.Id())
	p.Encode2s(w.Port())
	w.Conf().Encode(&p)
	return
}

// AddWorldServer assigns a world to the given world server connection
func AddWorldServer(con *Connection) error {
	var bindworld *World = nil
	var bindworldid int8 = -1

	for _, world := range worlds {
		if !world.Connected() { // we need to find a world that still isn't connected
			bindworld = world
			break
		}
	}

	if bindworld == nil { // no more worlds available
		con.SendPacket(interserver.NoMoreWorlds())
		return errors.New("No more worlds to assign.")
	}

	// assign world to conenction
	bindworldid = bindworld.Id()
	con.SetWorldId(bindworldid)
	bindworld.SetWorldCon(con)
	bindworld.SetConnected(true)

	// TODO: store external ip of the worldserver? check vana

	err := con.SendPacket(worldConnect(bindworld))
	if err != nil {
		return err
	}
	fmt.Println(con.Conn().RemoteAddr(), "assigned to world", bindworldid)
	return nil
}

// AddChannelServer assigns a channel to the given channel server connection
func AddChannelServer(con *Connection) error {
	var targetworld *World = nil
	var targetworldid int8 = -1

	worldIp := make([]byte, 4)

	// find a connected world that has room for channels
	for _, world := range worlds {
		if world.ChannelCount() < world.Conf().MaxChannels() && world.Connected() {
			targetworld = world
			break
		}
	}

	if targetworld == nil {
		con.SendPacket(interserver.LoginChannelConnect(targetworldid, worldIp, 0))
		return errors.New("No more channels to assign.")
	}

	targetworldid = targetworld.Id()
	// TODO: resolve this to the external ip address to actually make it work online
	// FIXME
	worldIp = common.RemoteAddrToBytes(con.Conn().RemoteAddr().String())
	if len(worldIp) != 4 {
		return errors.New("Ipv6 not supported")
	}

	err := con.SendPacket(interserver.LoginChannelConnect(targetworldid, worldIp, targetworld.Port()))
	if err != nil {
		return err
	}
	fmt.Println(con.Conn().RemoteAddr(), "assigned to world", targetworldid, "'s channels")
	return nil
}

// showWorld returns a packet to send world data to a client for the world list
// Send one for each world followed by a WorldListEnd() packet.
func showWorld(w *World) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(packets.OServerList)
	p.Encode1s(w.Id())
	p.EncodeString(fmt.Sprintf("%s World %d", w.Conf().Name(), w.Id()))
	p.Encode1(w.Conf().Ribbon())
	p.EncodeString(w.Conf().EventMsg())
	p.Encode2(100) // exp rate % (event message)
	p.Encode2(100) // drop rate % (event message)
	p.Encode1(0x00)
	p.Encode1(w.Conf().MaxChannels())

	for i := byte(0); i < w.Conf().MaxChannels(); i++ {
		p.EncodeString(fmt.Sprintf("%s World %d-%d", w.Conf().Name(), w.Id(), i+1))

		ch := w.Channel(int8(i))
		if ch != nil { // append channel population
			p.Encode4s(ch.Population())
		} else { // channel doesn't exist / crashed
			p.Encode4(0x00000000)
		}

		p.Encode1s(w.Id())
		p.Encode2(uint16(i))
	}

	p.Encode2(0x0000)
	return
}

// Show sends the world and channel list to the given client
func Show(con *client.Connection) (err error) {
	for _, world := range worlds {
		if !world.Connected() {
			continue
		}

		err = con.SendPacket(showWorld(world))
		if err != nil {
			return
		}
	}

	return con.SendPacket(packets.WorldListEnd())
}
