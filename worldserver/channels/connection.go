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

package channels

import "net"
import "github.com/Francesco149/kagami/common"

// A channls.Connection is a connection accepted by the channel server.
// It's a wrapper around EncryptedConnection specialized for inter-server communication.
// It handles authentification through the internal password.
type Connection struct {
	*common.InterserverConnection // underlying encrypted connection
	channelId                     int8
}

func (c *Connection) SetChannelId(channelId int8) { c.channelId = channelId }
func (c *Connection) ChannelId() int8             { return c.channelId }

// channls.Connection initializes a new inter-server world connection around a basic net.Conn
func NewConnection(con net.Conn, passwd string) *Connection {
	res := &Connection{
		InterserverConnection: common.NewInterserverConnection(con, passwd),
		channelId:             -1,
	}

	return res
}
