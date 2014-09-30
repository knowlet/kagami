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

package interserver

import (
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/maplelib"
)

// Server types for Auth()
const (
	WorldServer   = 0
	ChannelServer = 1
)

// Auth generates an inter-server authentication packet
func Auth(passwd string, serverType byte) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOAuth)
	p.EncodeString(passwd)
	p.Encode1(serverType)
	return
}

// NoMoreWorlds notifies a world server that there are no more worlds available to handle
func NoMoreWorlds() (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOWorldConnect)
	p.Encode1s(-1)
	return
}

// LoginChannelConnect returns a packet that notifies the channel
// server that has requested to connect to the loginserver
func LoginChannelConnect(worldId int8, ip []byte, port int16) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOLoginChannelConnect)
	p.Encode1s(worldId)
	p.Append(ip)
	p.Encode2s(port)
	return
}

// RemoveChannel returns a packet that requests the loginserver to remove a channel
func RemoveChannel(channelId int8) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IORemoveChannel)
	p.Encode1s(channelId)
	return
}

// RegisterChannel returns a packet that requests the loginserver to register a channel
func RegisterChannel(channelId int8, ipbytes []byte, channelPort int16) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IORegisterChannel)
	p.Encode1s(channelId)
	p.Append(ipbytes)
	p.Encode2s(channelPort)
	return
}

// SyncPlayerJoinedChannel returns a packet that notifies the worldserver that a player has joined a channel
func SyncPlayerJoinedChannel(channelid int8) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncPlayerJoinedChannel)
	p.Encode1s(channelid)
	return
}

// SyncPlayerLeftChannel returns a packet that notifies the worldserver that a player has left a channel
func SyncPlayerLeftChannel(channelid int8) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncPlayerLeftChannel)
	p.Encode1s(channelid)
	return
}

// SyncChannelPopulation returns a packet that notifies the loginserver to update a channel's population
func SyncChannelPopulation(worldid int8, channelid int8, population int32) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncChannelPopulation)
	p.Encode1s(worldid)
	p.Encode1s(channelid)
	p.Encode4s(population)
	return
	return
}
