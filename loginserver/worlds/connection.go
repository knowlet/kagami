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

package worlds

import "net"
import "github.com/knowlet/kagami/common"

// worlds.Connection represent a connection accepted by the channel server or world server.
// It's a wrapper around EncryptedConnection specialized for inter-server communication.
// It handles authentification through the internal password.
type Connection struct {
	*common.InterserverConnection // underlying encrypted connection
	worldId                       int8
}

func (c *Connection) SetWorldId(worldId int8) { c.worldId = worldId }
func (c *Connection) WorldId() int8           { return c.worldId }

// NewConnection initializes a new inter-server world connection around a basic net.Conn
func NewConnection(con net.Conn, passwd string) *Connection {
	res := &Connection{
		InterserverConnection: common.NewInterserverConnection(con, passwd),
		worldId:               -1,
	}

	return res
}
