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
	"github.com/Francesco149/kagami/common/player"
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

// MessageToChannel tells the worldserver that this packet must be relayed
// to a certain channel server
func MessageToChannel(channel int8, packet maplelib.Packet) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOMessageToChannel)
	p.Encode1s(channel)
	p.Append([]byte(packet))
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

// SyncWorldCharacterCreated returns a packet that notifies the world server that a character has been created
func SyncWorldCharacterCreated(charid int32) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncWorldCharacterCreated)
	p.Encode4s(charid)
	return
}

// SyncWorldCharacterDeleted returns a packet that notifies the world server that a character has been deleted
func SyncWorldCharacterDeleted(charid int32) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncWorldCharacterDeleted)
	p.Encode4s(charid)
	return
}

// SyncChannelCharacterCreated returns a packet that notifies the channel server that a character has been created
func SyncChannelCharacterCreated(char *player.Data) (p maplelib.Packet, err error) {
	p = packets.NewEncryptedPacket(IOSyncChannelCharacterCreated)
	err = char.Encode(&p)
	return
}

// SyncChannelCharacterDeleted returns a packet that notifies the channel server that a character has been deleted
func SyncChannelCharacterDeleted(charid int32) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncChannelCharacterDeleted)
	p.Encode4s(charid)
	return
}

// SyncChannelNewPlayer notifies the channel server that we're connecting to that channel
func SyncChannelNewPlayer(charId int32, ip []byte) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncChannelNewPlayer)
	p.Encode4s(charId)
	p.EncodeBuffer(ip)
	return
}

// SyncWorldPerformChangeChannel notifies the world server that the player can be
// transferred to the desired channel
func SyncWorldPerformChangeChannel(charId int32) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncWorldPerformChangeChannel)
	p.Encode4s(charId)
	return
}

// SyncWorldPerformChangeChannel notifies the channel server that the player has been
// transferred to another channel
func SyncChannelPerformChangeChannel(id int32, newchanid int8, ip []byte, port int16) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(IOSyncChannelPerformChangeChannel)
	p.Encode4s(id)
	p.Encode1s(newchanid)
	p.EncodeBuffer(ip)
	p.Encode2s(port)
	return
}

// SyncWorldLoadCharacter notifies the world server that a character joined the channel server and
// sends it the player's data
func SyncWorldLoadCharacter(data *player.Data) (p maplelib.Packet, err error) {
	p = packets.NewEncryptedPacket(IOSyncWorldLoadCharacter)
	err = data.Encode(&p)
	return
}

// SyncChannelUpdatePlayer sends a character's data to the channel server,
// telling it to update its local cache of the data
func SyncChannelUpdatePlayer(data *player.Data) (p maplelib.Packet, err error) {
	p = packets.NewEncryptedPacket(IOSyncChannelUpdatePlayer)
	err = data.Encode(&p)
	return
}
