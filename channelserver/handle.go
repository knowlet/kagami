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
	"github.com/Francesco149/kagami/channelserver/players"
	"github.com/Francesco149/kagami/channelserver/status"
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/kagami/common/packets"
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
	}

	return false, nil // forward packet to next handler
}

// connectData returns a packet that sends the initial character data when
// a player connects to the channelserver
func connectData(con *client.Connection) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(packets.OConnectData)
	// TODO: add all missing data
	p.Encode4s(int32(status.ChanId()))
	p.Encode1(0x01)          // what the hell is this
	p.Encode1(0x01)          // what the hell is this
	p.Encode2(0x0000)        // what the hell is this
	p.Encode4s(rand.Int31()) // rng seed
	whatthehellisthis := []byte{0xF8, 0x17, 0xD7, 0x13, 0xCD, 0xC5, 0xAD, 0x78}
	p.Append(whatthehellisthis) // what the hell is this
	p.Encode8s(-1)              // what the hell is this

	err := con.EncodeStats(&p)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	p.Encode1(100)  // TODO: get real buddylist capacity
	p.Encode4(1337) // TODO: get real meso
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

	charip := common.RemoteAddrToBytes(con.Conn().RemoteAddr().String())

	// look for the character in the pending connections list then check the ip
	// if the ips don't match, then someone is trying to remote hack
	starttime := time.Now().Unix()
	for {
		if time.Now().Unix()-starttime > 30 {
			break
		}

		players.Lock()
		expectedip := players.PendingIp(charid)
		players.Unlock()
		if expectedip != nil {
			if !bytes.Equal(charip, expectedip) {
				err = errors.New(fmt.Sprint(common.BytesToIpString(charip), "tried to remote hack",
					common.BytesToIpString(expectedip)))
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
	con.SetCharId(charid)

	// get char data from db
	db := common.GetDB()
	st, err := db.Prepare("SELECT c.*, a.gm_level, a.admin FROM `characters` c " +
		"INNER JOIN `accounts` a ON c.user_id = a.id " +
		"WHERE c.character_id = ?")
	if err != nil {
		fmt.Println("Unexpected invalid query in handleLoadCharacter")
		return
	}
	if st == nil {
		err = errors.New("handleLoadCharacter: wat")
		fmt.Println("handleLoadCharacter: wat")
		return
	}
	res, err := st.Run(charid)
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	if len(rows) < 1 {
		err = errors.New("Character not found.")
		return
	}

	row := rows[0]

	colname := res.Map("name")
	coluserid := res.Map("user_id")
	colmap := res.Map("map")
	colgmlevel := res.Map("gm_level")
	coladmin := res.Map("admin")
	colface := res.Map("face")
	colhair := res.Map("hair")
	colworldid := res.Map("world_id")
	colgender := res.Map("gender")
	colskin := res.Map("skin")
	colpos := res.Map("pos")

	con.SetName(row.Str(colname))
	con.SetUserId(int32(row.Int(coluserid)))
	con.SetMapId(int32(row.Int(colmap)))
	con.SetGmLevel(int32(row.Int(colgmlevel)))
	con.SetAdmin(row.Int(coladmin) > 0)
	con.SetFace(int32(row.Int(colface)))
	con.SetHair(int32(row.Int(colhair)))
	con.SetWorldId(int8(row.Int(colworldid)))
	con.SetGender(byte(row.Int(colgender)))
	con.SetSkin(int8(row.Int(colskin)))
	con.SetMapPos(int8(row.Int(colpos)))

	// TODO: get buddylist size
	// TODO: get stats
	// TODO: get max inventory slots & meso

	// TODO: do not reset uptime if the player is just xfering

	con.SetUptime(0)
	con.SetGmChat(con.GmChat() && con.GmLevel() > 0)

	// TODO: get book cover (wtf is a book cover)
	// TODO: init keymaps
	// TODO: init hpmp

	// TODO: check forced return map
	// TODO: check if the player is dead and repawn him

	// TODO: init position, stance and foothold

	status.Lock()
	con.SendPacket(connectData(con))

	conf := status.WorldConf()
	if len(conf.ScrollingHeader()) != 0 {
		err = con.SendPacket(packets.ScrollingHeader(conf.ScrollingHeader()))
		if err != nil {
			return
		}
	}
	status.Unlock()

	// TODO: init pets
	// TODO: send keymaps
	// TODO: send update buddylist
	// TODO: check for pending buddylist requests
	// TODO: send skill macros

	// TODO: add player to player list
	// TODO: add player to map's player list

	fmt.Println(con.Conn().RemoteAddr().String(), "connected as", con.Name())

	err = con.SetDBOnline(true)
	if err != nil {
		return
	}

	con.SetConnected(true)
	fmt.Println(con.String())

	status.Lock()
	defer status.Unlock()
	status.WorldConn().SendPacket(interserver.SyncPlayerJoinedChannel(status.ChanId()))

	// TODO: add to player pool

	handled = err == nil
	return
}
