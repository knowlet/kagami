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

import "fmt"

// Connect waits for and connects to a tcp server on a given port.
// handler is the function that will handle this connection's packets, see PacketHandler for the signature.
// makeConnection is a connection factory function that must return a connection that implements common.Connection.
// Once a connection is estabilished, a loop will run to handle its packets, blocking the current thread.
func Connect(name, ipport string, handler PacketHandler, makeConnection ConnectionFactory) {
	fmt.Println("Connecting to", name, ipport)
	con, err := Dial(ipport)
	if err != nil {
		fmt.Println("Failed to connect: ", err)
		return
	}

	fmt.Println("Connected to", name, con.RemoteAddr())
	HandleLoop(name, con, handler, makeConnection, nil)
}
