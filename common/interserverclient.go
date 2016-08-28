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

import "net"
import "github.com/knowlet/kagami/common/interserver"

// An InterserverClient is a connection to another component of the server.
// It's a wrapper around EncryptedConnection specialized for inter-server communication.
// It handles authentification through the internal password.
type InterserverClient struct {
	*EncryptedConnection // underlying encrypted connection
}

// NewInterserverClient initializes a new inter-server connection around a basic net.Conn
func NewInterserverClient(con net.Conn, passwd string, serverType byte) *InterserverClient {
	res := &InterserverClient{
		EncryptedConnection: NewEncryptedConnection(con, false, true),
	}

	auth := interserver.Auth(passwd, serverType)
	res.SendPacket(auth)
	return res
}
