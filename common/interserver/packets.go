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

// ConnectingToChannel notifies the world server that we're connecting to a channel
func ConnectingToChannel(channel int8, charId int32, ip []byte) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOConnectingToChannel)
	p.Encode1s(channel)
	p.Encode4s(charId)
	p.Append(ip)
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
