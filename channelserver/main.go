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
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"github.com/Francesco149/kagami/channelserver/client"
	"github.com/Francesco149/kagami/channelserver/gamedata"
	"github.com/Francesco149/kagami/channelserver/players"
	"github.com/Francesco149/kagami/channelserver/status"
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/maplelib"
	"github.com/Francesco149/maplelib/wz"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Kagami Pre-Alpha")
	fmt.Println("Initializing ChannelServer...")

	fmt.Println("Loading Map.wz and String.wz")
	mapProvider, err := wz.NewMapleDataProvider("wz/Map.wz")
	checkError(err)
	stringProvider, err := wz.NewMapleDataProvider("wz/String.wz")
	checkError(err)
	factory, err := gamedata.NewMapleMapFactory(mapProvider, stringProvider)
	checkError(err)

	status.Init()
	st := <-status.Get
	st.SetMapProvider(mapProvider)
	st.SetStringsProvider(stringProvider)
	st.SetMapFactory(factory)
	status.Get <- st

	// disconnect all players if the server panics or is closed non-gracefully
	fnCleanup := func() {
		fmt.Println("Attempting emergency cleanup...")
		err := players.Execute(func(con *client.Connection) error {
			return con.SetDBOnline(false)
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Success!")
		}
	}

	// handle panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r, "\n\nRecovered from panic")
			fnCleanup()
		}
	}()

	// handle SIGINT
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGKILL)
	signal.Notify(sigint, syscall.SIGHUP)
	signal.Notify(sigint, syscall.SIGTRAP)
	signal.Notify(sigint, syscall.SIGQUIT)
	go func() {
		sig := <-sigint
		fmt.Println("Captured signal", sig)
		fnCleanup()
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}()

	// connect to loginserver
	fmt.Println("Waiting for the loginserver to assign a worldserver...")
	common.Connect("loginserver", fmt.Sprintf("%s:%d", consts.LoginIp, consts.LoginInterserverPort),
		func(con common.Connection, p maplelib.Packet) (bool, error) {
			scon, ok := con.(*common.InterserverClient)
			if !ok {
				return false, errors.New("Loginserver handler failed type assertion")
			}
			return HandleInter(scon, p)
		},
		func(con net.Conn) common.Connection {
			c := common.NewInterserverClient(con, consts.InterServerPassword,
				interserver.ChannelServer)
			st := <-status.Get
			defer func() { status.Get <- st }()
			st.SetLoginConn(c)
			return c
		})
}
