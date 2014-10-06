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
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

import (
	"github.com/Francesco149/kagami/channelserver/client"
	"github.com/Francesco149/kagami/channelserver/gamedata"
	"github.com/Francesco149/kagami/channelserver/players"
	"github.com/Francesco149/kagami/channelserver/status"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/kagami/common/utils"
	"github.com/Francesco149/maplelib"
)

// Handle handles channelserver packets
func Handle(con *client.Connection, p maplelib.Packet) (handled bool, err error) {
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	// Refuse any packet except the one for loading the character until the player is connected
	if !con.Connected() {
		if header == packets.ILoadCharacter {
			return handleLoadCharacter(con, it)
		}
	}

	switch header {
	case packets.IUnknownPlsIgnore1:
		return true, nil

	case packets.IUnknownPlsIgnore2:
		return true, nil

	case packets.IPlayerUpdate:
		return handlePlayerUpdate(con)

	case packets.IChangeMapSpecial:
		return handleChangeMapSpecial(con, it)

	case packets.IChangeMap:
		return handleChangeMap(con, it)
	}

	return false, nil // forward packet to next handler
}

// connectData returns a packet that sends the initial character data when
// a player connects to the channelserver
func connectData(con *client.Connection, chanid int8) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(packets.OWarpToMap)
	// TODO: add all missing data
	p.Encode4s(int32(chanid))
	p.Encode1(0x01)          // what the hell is this
	p.Encode1(0x01)          // what the hell is this
	p.Encode2(0x0000)        // what the hell is this
	p.Encode4s(rand.Int31()) // rng seed
	whatthehellisthis := []byte{0xF8, 0x17, 0xD7, 0x13, 0xCD, 0xC5, 0xAD, 0x78}
	p.Append(whatthehellisthis) // what the hell is this
	p.Encode8s(-1)              // what the hell is this

	con.Stats().Encode(&p)
	p.Encode1(con.BuddylistSize())
	p.Encode4s(con.Meso())

	// TODO: get real inv slots
	p.Encode1(100) // equip slots
	p.Encode1(100) // use slots
	p.Encode1(100) // set-up slots
	p.Encode1(100) // etc slots
	p.Encode1(100) // cash slots

	// TODO: encode equips
	p.Encode2(0x0000) // inventories are zero-terminated lists
	// TODO: encode equip inventory
	p.Encode1(0x00)
	// TODO: encode use inventory
	p.Encode1(0x00)
	// TODO: encode set-up inventory
	p.Encode1(0x00)
	// TODO: encode etc inventory
	p.Encode1(0x00)
	// TODO: encode cash inventory
	p.Encode1(0x00)
	p.Encode2(0x0000) // 0 skills for now (placeholder)
	// TODO: encode skills id's here
	p.Encode2(0x0000)
	con.EncodeQuestInfo(&p)
	// TODO: encode rings
	p.Encode8(0x0000000000000000)

	magic := []byte{0xFF, 0xC9, 0x9A, 0x3B}
	for i := 0; i < 15; i++ {
		p.Append(magic)
	}

	p.Encode4(0x00000000)
	p.Encode8s(time.Now().UnixNano() / 1000000) // time in millisecs
	return
}

// handleLoadCharacter handles the packet for loading a player's character when the player first
// connects to the channelserver
func handleLoadCharacter(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false

	charid, err := it.Decode4s()
	if err != nil {
		return
	}

	charip := utils.RemoteAddrToBytes(con.Conn().RemoteAddr().String())

	// look for the character in the pending connections list then check the ip
	// if the ips don't match, then someone is trying to remote hack
	starttime := time.Now().Unix()
	for {
		if time.Now().Unix()-starttime > 30 {
			err = errors.New("Pending connection timed out")
			break
		}

		players.Lock()
		expectedip := players.PendingIp(charid)
		players.Unlock()
		if expectedip != nil {
			if !bytes.Equal(charip, expectedip) {
				err = errors.New(fmt.Sprint(utils.BytesToIpString(charip),
					"tried to remote hack",
					utils.BytesToIpString(expectedip)))
			}
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return
	}

	players.Lock()
	players.RemovePendingIp(charid)
	players.Unlock()

	err = con.LoadFromDB(charid)
	if err != nil {
		return
	}

	// TODO: check forced return map
	// TODO: check if the player is dead and repawn him

	// TODO: init position, stance and foothold

	stts := <-status.Get
	defer func() { status.Get <- stts }()

	con.SendPacket(connectData(con, stts.ChanId()))

	if len(stts.WorldConf().ScrollingHeader()) != 0 {
		err = con.SendPacket(packets.ScrollingHeader(stts.WorldConf().ScrollingHeader()))
		if err != nil {
			return
		}
	}

	// TODO: init pets
	// TODO: send keymaps
	// TODO: send update buddylist
	// TODO: check for pending buddylist requests
	// TODO: send skill macros

	// TODO: add player to player list
	// TODO: add player to map's player list

	fmt.Println(con.Conn().RemoteAddr().String(), "connected as", con.Stats().Name())

	err = con.SetDBOnline(true)
	if err != nil {
		return
	}

	con.SetConnected(true)
	fmt.Println(con.String())

	stts.WorldConn().SendPacket(interserver.SyncPlayerJoinedChannel(stts.ChanId()))

	// TODO: add to player pool

	handled = err == nil
	return
}

// handlePlayerUpdate handles a request to save the player's data to the database
func handlePlayerUpdate(con *client.Connection) (handled bool, err error) {
	// TODO: rate check on this packet to prevent clients from spamming it to flood the database
	err = con.Save()
	handled = err == nil
	return
}

// handleChangeMapSpecial handles a special map change packet
func handleChangeMapSpecial(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	_, err = it.Decode1()
	portalname, err := it.DecodeString()
	_, err = it.Decode1()
	_, err = it.Decode1()
	if err != nil {
		return
	}

	fmt.Println(con.Stats().Name(), "entered", portalname, "in map", con.Stats().MapId())
	portal := con.Map().Portal(portalname)
	if portal == nil {
		fmt.Println("Enabled actions for", con.Stats().Name())
		err = con.SendPacket(packets.EnableActions())
	} else {
		fmt.Println("Sending portal enter packet")
		err = con.Enter(portal.(gamedata.IMapleGenericPortal))
	}

	handled = err == nil
	return
}

// handleChangeMap handles a map change or revival packet
func handleChangeMap(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	_, err = it.Decode1()
	target, err := it.Decode4s()
	portalname, err := it.DecodeString()
	portal := con.Map().Portal(portalname)

	if target != -1 && !con.Alive() {
		fmt.Println(con.Stats().Name(), "died")
	} else {
		fmt.Println(con.Stats().Name(), "is entering portal",
			portalname)
	}

	switch {
	case target != -1 && !con.Alive():
		fmt.Println("TODO: revive player")

	case target != -1 && con.GmLevel() > 2:
		// TODO: check chalkboard
		oldmap := con.Stats().MapId()
		err = con.SetMapId(target)
		if err != nil {
			con.SetMapId(oldmap)
			err = nil
		} else {
			portal = con.Map().PortalById(0)
			con.WarpToMap(con.Map(), portal)
		}

	case target != -1 && con.GmLevel() <= 2:
		fmt.Println(con.Stats().Name(), "tried to map warp without gm powers.")

	default:
		if portal != nil {
			err = con.Enter(portal.(gamedata.IMapleGenericPortal))
		} else {
			err = con.SendPacket(packets.EnableActions())
		}
	}

	handled = err == nil
	return
}
