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

// clientLoop sends the handshake and handles packets for a single client in a loop
func clientLoop(basecon net.Conn) {
	defer basecon.Close()
	con := client.NewConnection(basecon, false)

	for {
		inpacket, err := con.RecvPacket()
		if err != nil {
			fmt.Println(err)
			break
		}

		handled, err := common.Handle(con, inpacket)
		if err != nil {
			fmt.Println(err)
			break
		}

		if !handled {
			handled, err = Handle(con, inpacket)
			if err != nil {
				fmt.Println(err)
				break
			}
		}

		if !handled {
			fmt.Println("Unhandled packet", inpacket)
			//break
		}
	}

	fmt.Println("Dropping: ", con.Conn().RemoteAddr())
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Kagami Pre-Alpha")
	fmt.Println("Initializing LoginServer...")

	fmt.Println("Loading worlds...")
	loadDefaultWorlds()

	sock, err := common.NewTcpServer(fmt.Sprintf(":%d", consts.LoginPort))
	if err != nil {
		fmt.Println("Failed to create socket: ", err)
		return
	}

	fmt.Println("Listening on port", consts.LoginPort)

	for {
		con, err := sock.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err)
			return
		}

		fmt.Println("Accepted: ", con.RemoteAddr())
		go clientLoop(con)
	}
}
