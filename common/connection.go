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

package common

import "github.com/Francesco149/maplelib"
import "net"

// A connection is a generic wrapper for a client connected to a packet oriented protocol
type Connection interface {
	// RecvPacket listens for an incoming packet and reads it.
	// The implementation must retrieve the packet length from the incoming data.
	RecvPacket() (packet maplelib.Packet, err error)

	// SendPacket sends a packet through this connection.
	// The implementation must provide its own way to let the receiver figure out
	// the length of the incoming packet
	SendPacket(packet maplelib.Packet) error

	// Ping sends a ping packet to the client
	// if the connection is a client, this will send a pong instead
	Ping() error

	// OnPong must be called when a pong is received
	OnPong() error

	// Conn returns the underlying connection
	Conn() net.Conn

	// IsClient returns true if the connection is a client connected to a server
	IsClient() bool
}
