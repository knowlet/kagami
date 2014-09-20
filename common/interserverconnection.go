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
	"errors"
	"fmt"
	"net"
)

import "github.com/Francesco149/maplelib"

// InterserverConnection represent a connection accepted from another component of the server.
// It's a wrapper around EncryptedConnection specialized for inter-server communication.
// It handles authentification through the internal password.
type InterserverConnection struct {
	*EncryptedConnection        // underlying encrypted connection
	password             string // inter-server password
	authenticated        bool   // true if the other end has corrently sent it's internal password
}

// NewInterserverConnection initializes a new inter-server connection around a basic net.Conn
func NewInterserverConnection(con net.Conn, passwd string) *InterserverConnection {
	res := &InterserverConnection{
		EncryptedConnection: NewEncryptedConnection(con, false, false),
		password:            passwd,
		authenticated:       false,
	}

	return res
}

func (c *InterserverConnection) Authenticated() bool         { return c.authenticated }
func (c *InterserverConnection) SetPassword(password string) { c.password = password }

// CheckAuth parses an inter-server auth packet and sets the connection as authenticated if successful
func (c *InterserverConnection) CheckAuth(it maplelib.PacketIterator) (err error) {
	password, err := it.DecodeString()
	if err != nil {
		return
	}

	if c.password != password {
		err = errors.New(fmt.Sprintf("%s is not a valid inter-server password.", password))
		return
	}

	// TODO: check that the ip is allowed to connect to the server

	c.authenticated = true
	fmt.Println(c.Conn().RemoteAddr(), "authenticated inter-server connection")
	return nil
}
