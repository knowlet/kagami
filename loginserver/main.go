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
	"github.com/Francesco149/kagami/loginserver/client"
	"github.com/Francesco149/kagami/loginserver/worlds"
	"github.com/Francesco149/maplelib"
)

// loadDefaultWorlds loads and adds the default world list to the loginserver
func loadDefaultWorlds() {
	// TODO: config files
	configs := config.DefaultWorldConf()

	for i, config := range configs {
		id := consts.WorldId[i]

		world := worlds.Get(id)
		if world != nil {
			world.SetConf(config)
			continue // only refresh config
		}

		// add new world
		world = &worlds.World{}
		world.SetConf(config)
		world.SetId(id)
		world.SetPort(consts.WorldListenPort[i])
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Kagami Pre-Alpha")
	fmt.Println("Initializing LoginServer...")

	fmt.Println("Loading worlds...")
	loadDefaultWorlds()

	go common.Accept("world", consts.LoginInterserverPort,
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*common.InterserverConnection)
			if !ok {
				return false, errors.New("World handler failed type assertion")
			}
			return HandleInter(scon, p)
		},
		func(con net.Conn) common.Connection {
			return common.NewInterserverConnection(con, consts.InterServerPassword)
		})

	common.Accept("client", consts.LoginPort,
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*client.Connection)
			if !ok {
				return false, errors.New("Client handler failed type assertion")
			}
			return Handle(scon, p)
		},
		func(con net.Conn) common.Connection {
			return client.NewConnection(con, false)
		})
}
