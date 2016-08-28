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

import (
	"fmt"
	"net"
)

import (
	"github.com/knowlet/kagami/common/utils"
	"github.com/knowlet/maplelib"
)

// A PacketHandler is a generic packet handling function signature
type PacketHandler func(con Connection, p maplelib.Packet) (handled bool, err error)

// A ConnectionFactory is a callback that must allocate a struct that implements common.Connection
type ConnectionFactory func(net.Conn) Connection

// A DisconnectCallback is a callback that will be called when the connection is closed
type DisconnectCallback func(con Connection)

// handleLoop is the packet handling / sending loop for a single connected client
func HandleLoop(name string, basecon net.Conn, handler PacketHandler,
	makeConnection ConnectionFactory, onDisconnect DisconnectCallback) {

	// handle panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r, "\n\nRecovered from panic")
		}
	}()

	defer basecon.Close()
	con := makeConnection(basecon)

	for {
		inpacket, err := con.RecvPacket()
		if err != nil {
			fmt.Println(err)
			break
		}

		handled, err := Handle(con, inpacket)
		if err != nil {
			fmt.Println(err)
			break
		}

		if !handled {
			handled, err = handler(con, inpacket)
			if err != nil {
				fmt.Println(utils.MakeError(err.Error()))
				break
			}
		}

		if !handled {
			fmt.Println(utils.MakeWarning("Unhandled ", name, " packet ", inpacket))
			//break
		}
	}

	if onDisconnect != nil {
		onDisconnect(con)
	}
	fmt.Println("Dropping", name, con.Conn().RemoteAddr())
}

// Accept waits and accepts connections on a given port.
// handler is the function that will handle this connection's packets, see PacketHandler for the signature.
// makeConnection is a connection factory function that must return a connection that implements common.Connection.
// onDisconnect is a callback that will be called once the connection is dropped. This is optional, pass nil to ignore it.
// Once a connection is accepted, a new thread will be started to handle its packets.
func Accept(name string, port int16, handler PacketHandler,
	makeConnection ConnectionFactory, onDisconnect DisconnectCallback) {
	sock, err := Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Failed to create socket: ", err)
		return
	}

	fmt.Println("Listening for", name, "on port", port)

	for {
		con, err := sock.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err)
			return
		}

		fmt.Println("Accepted", name, con.RemoteAddr())
		go HandleLoop(name, con, handler, makeConnection, onDisconnect)
	}
}
