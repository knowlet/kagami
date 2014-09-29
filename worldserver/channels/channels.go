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

// Package channels contains data structures and utilities to keep track of the server's channels
package channels

import "sync"

import (
	"github.com/Francesco149/kagami/worldserver/status"
	"github.com/Francesco149/maplelib"
)

var mut sync.Mutex
var channels = make(map[int8]*Channel) // channels mapped by id

// Lock locks the channels mutex.
// Must be called before performing any operation on
// the channels
func Lock() {
	mut.Lock()
}

// Unlock unlocks the channels mutex.
func Unlock() {
	mut.Unlock()
}

// channels.Add creates and adds a new channel to the list
func Add(con *Connection, chanid int8, chanip []byte, port int16) {
	// TODO: resolve newchan's external ip
	channels[chanid] = NewChannel(con, chanid, port)
}

// channels.Remove removes a channel from the list
func Remove(chanid int8) {
	delete(channels, chanid)
}

// channels.Get gets a channel by id. Returns nil if the id doesn't exist.
func Get(chanid int8) *Channel {
	return channels[chanid]
}

func clonePacket(p maplelib.Packet) maplelib.Packet {
	bytes := []byte(p)
	clone := make([]byte, len(bytes))
	copy(clone, bytes)
	return maplelib.Packet(clone)
}

// channels.SendToChannelList sends a packet to a list of channel id's
func SendToChannelList(channelids []int8, p maplelib.Packet) (err error) {
	for _, id := range channelids {
		err = Get(id).Conn().SendPacket(clonePacket(p))
		if err != nil {
			return
		}
	}

	return
}

// channels.SendToAllChannels sends a packet to all of the channels
func SendToAllChannels(p maplelib.Packet) (err error) {
	for _, ch := range channels {
		err = ch.Conn().SendPacket(clonePacket(p))
		if err != nil {
			return
		}
	}

	return
}

// channels.GetFirstAvailableId returns the first available channel id
// returns -1 if there are no more available channel id's
func GetFirstAvailableId() int8 {
	var id, max int8 = -1, int8(status.Conf().MaxChannels())

	for i := int8(0); i < max; i++ {
		if Get(i) == nil { // find channel id's that are still not mapped
			id = i
			break
		}
	}

	return id
}
