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
	"github.com/knowlet/kagami/common"
	"github.com/knowlet/kagami/common/config"
	"github.com/knowlet/kagami/common/consts"
	"github.com/knowlet/kagami/loginserver/client"
	"github.com/knowlet/kagami/loginserver/worlds"
	"github.com/Francesco149/maplelib"
)

// loadDefaultWorlds loads and adds the default world list to the loginserver
func loadDefaultWorlds() {
	worlds.Lock()
	defer worlds.Unlock()

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
		world = worlds.NewWorld(config, id, consts.WorldListenPort[i])
		worlds.Add(world)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Kagami Pre-Alpha")
	fmt.Println("Initializing LoginServer...")

	fmt.Println("Loading worlds...")
	loadDefaultWorlds()

	// accept interserver world connections in a separate thread
	go common.Accept("world/chan", consts.LoginInterserverPort,
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*worlds.Connection)
			if !ok {
				return false, errors.New("World handler failed type assertion")
			}
			return HandleInter(scon, p)
		},
		func(con net.Conn) common.Connection {
			return worlds.NewConnection(con, consts.InterServerPassword)
		},
		func(con common.Connection) {
			worlds.Lock()
			defer worlds.Unlock()

			scon, ok := con.(*worlds.Connection)
			if !ok {
				panic(errors.New("World handler failed type assertion on disconnect"))
			}
			deleteworldid := scon.WorldId()

			if deleteworldid == -1 {
				return
			}

			fmt.Println("Removing world", deleteworldid)
			deleteworld := worlds.Get(deleteworldid)

			if deleteworld == nil {
				fmt.Println("Could not find world", deleteworldid)
				return
			}

			deleteworld.SetConnected(false)
			deleteworld.ClearChannels()
		})

	// accept client connections in this thread
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
		},
		nil)
}
