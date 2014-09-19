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

import "fmt"

import (
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/kagami/loginserver/client"
	"github.com/Francesco149/maplelib"
)

var worlds = make(map[byte]*World)

// Add adds the given world to the world list
func Add(w *World) {
	worlds[w.Id()] = w
}

// Get returns the world associated with the given id
func Get(worldId byte) *World {
	return worlds[worldId]
}

// showWorld returns a packet to send world data to a client for the world list
// Send one for each world followed by a WorldListEnd() packet.
func showWorld(w *World) (p maplelib.Packet) {
	p.Encode4(0x00000000)
	p.Encode2(packets.OServerList)
	p.Encode1(w.Id())
	p.EncodeString(fmt.Sprintf("%s World %d", w.Conf().Name(), w.Id()))
	p.Encode1(w.Conf().Ribbon())
	p.EncodeString(w.Conf().EventMsg())
	p.Encode2(100) // exp rate % (event message)
	p.Encode2(100) // drop rate % (event message)
	p.Encode1(0x00)
	p.Encode1(w.Conf().MaxChannels())

	for i := byte(0); i < w.Conf().MaxChannels(); i++ {
		p.EncodeString(fmt.Sprintf("%s World %d-%d", w.Conf().Name(), w.Id(), i+1))

		ch := w.Channel(i)
		if ch != nil { // append channel population
			p.Encode4s(ch.Population())
		} else { // channel doesn't exist / crashed
			p.Encode4(0x00000000)
		}

		p.Encode1(w.Id())
		p.Encode2(uint16(i)) // index on the server list?
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
