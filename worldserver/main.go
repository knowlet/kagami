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
	"github.com/knowlet/kagami/common/consts"
	"github.com/knowlet/kagami/common/interserver"
	"github.com/knowlet/maplelib"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Kagami Pre-Alpha")
	fmt.Println("Initializing WorldServer...")

	// TODO: accept channelserver connections
	common.Connect("loginserver", fmt.Sprintf("%s:%d", consts.LoginIp, consts.LoginInterserverPort),
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*common.InterserverClient)
			if !ok {
				panic(errors.New("Login handler failed type assertion"))
			}
			return HandleLogin(scon, p)
		},
		func(con net.Conn) common.Connection {
			return common.NewInterserverClient(con, consts.InterServerPassword, interserver.WorldServer)
		})
}
