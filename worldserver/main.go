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

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/config"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/maplelib"
)

// TODO: everything

var worldconf *config.WorldConf = nil
var worldport int16 = -1

func Handle(con *common.InterserverClient, p maplelib.Packet) (handled bool, err error) {
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	switch header {
	case interserver.IOWorldConnect:
		return handleWorldConnect(con, it)
	}

	return false, nil
}

// handleWorldConnect handles a world connect packet
func handleWorldConnect(con *common.InterserverClient, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	worldid, err := it.Decode1s()
	if err != nil {
		return
	}

	if worldid == -1 {
		fmt.Println("No worlds to handle!")
		return
	}

	port, err := it.Decode2s()
	conf, err := config.DecodeWorldConf(&it)
	if err != nil {
		return
	}

	handled = true
	fmt.Println("Handling world", worldid)
	worldconf = conf
	worldport = port

	// TODO: start accepting connections

	fmt.Println("World server is running!")
	return
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Kagami Pre-Alpha")
	fmt.Println("Initializing WorldServer...")

	// TODO: accept channelserver connections
	common.Connect("loginserver", fmt.Sprintf("%s:%d", consts.LoginIp, consts.LoginInterserverPort),
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*common.InterserverClient)
			if !ok {
				return false, errors.New("Loginserver handler failed type assertion")
			}
			return Handle(scon, p)
		},
		func(con net.Conn) common.Connection {
			return common.NewInterserverClient(con, consts.InterServerPassword)
		})
}
