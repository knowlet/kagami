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

// Package players contains the worldserver player pool and utilities to manage them
package players

import (
	"errors"
	"fmt"
	"sync"
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/player"
)

var mut sync.Mutex
var players = make(map[int32]*player.Data)        // player data mapped by id
var playersByName = make(map[string]*player.Data) // player data mapped by name
var pendingCC = make(map[int32]*int8)             // pending cc requests mapped by player id

// Lock locks the players pool mutex. Call it before performing any operation.
func Lock() {
	mut.Lock()
}

// Unlock unlocks the players pool mutex.
func Unlock() {
	mut.Unlock()
}

// LoadWorld loads all characters from the given world
func LoadWorld(worldid int8) (err error) {
	fmt.Println("Initializing player pool...")

	db := common.GetDB()
	st, err := db.Prepare("SELECT character_id, name FROM characters " +
		"WHERE world_id = ?")
	if err != nil {
		return
	}
	res, err := st.Run(worldid)
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	colcharid := res.Map("character_id")
	colname := res.Map("name")

	for _, row := range rows {
		data := &player.Data{}
		data.SetIp(make([]byte, 4))
		data.SetCharId(int32(row.Int(colcharid)))
		data.SetName(row.Str(colname))
		Add(data)
	}

	return
}

// Add adds a player's data to the player pool
func Add(data *player.Data) {
	players[data.CharId()] = data
	playersByName[data.Name()] = data
}

// Get returns a player's data by id
func Get(id int32) *player.Data {
	return players[id]
}

// Load lods the given character id's data into the players pool
func Load(charid int32) (err error) {
	if players[charid] != nil {
		return
	}

	db := common.GetDB()
	st, err := db.Prepare("SELECT character_id, name FROM characters " +
		"WHERE character_id = ?")
	if err != nil {
		return
	}
	res, err := st.Run(charid)
	rows, err := res.GetRows()
	if err != nil {
		return
	}
	if len(rows) < 1 {
		err = errors.New("character not found")
	}

	row := rows[0]
	colcharid := res.Map("character_id")
	colname := res.Map("name")

	data := &player.Data{}
	data.SetIp(make([]byte, 4))
	data.SetCharId(int32(row.Int(colcharid)))
	data.SetName(row.Str(colname))
	Add(data)
	return
}

// GetPendingCC gets the given player's target channel.
// if the player isn't ccing, the function returns -1
func GetPendingCC(charid int32) int8 {
	c := pendingCC[charid]
	if c == nil {
		return -1
	}
	return *c
}

// RemovePendingCC removes a player's pending cc request
func RemovePendingCC(charid int32) {
	if pendingCC[charid] == nil {
		return
	}

	delete(pendingCC, charid)
}
