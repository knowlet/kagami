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

// Package players contains the channelserver player pool and utilities to manage them
package players

import (
	"errors"
	"sync"
	"time"
)

import (
	"github.com/Francesco149/kagami/channelserver/client"
	"github.com/Francesco149/kagami/common/player"
)

// ConnectingPlayer holds information of a pending player connection to
// the channel server
type ConnectingPlayer struct {
	ip            []byte
	unixtimestamp int64
	mapid         int32
}

var mut sync.Mutex
var players = make(map[int32]*client.Connection)        // player connections mapped by id
var playersByName = make(map[string]*client.Connection) // player connections mapped by name
var playersData = make(map[int32]*player.Data)          // player data mapped by id
var playersDataByName = make(map[string]*player.Data)   // player data mapped by name
var connecting = make(map[int32]*ConnectingPlayer)      // connecting players mapped by id

// Lock locks the players pool mutex. Call it before performing any operation.
func Lock() {
	mut.Lock()
}

// Unlock unlocks the players pool mutex.
func Unlock() {
	mut.Unlock()
}

// Add adds a player's connection to the player pool
func Add(data *client.Connection) {
	players[data.CharId()] = data
	playersByName[data.Name()] = data
}

// Add adds a player's data to the player pool
func AddData(data *player.Data) {
	playersData[data.CharId()] = data
	playersDataByName[data.Name()] = data
}

// Get returns a player's connection by id
func Get(id int32) *client.Connection {
	return players[id]
}

// GetData returns a player's data by id
func GetData(id int32) *player.Data {
	return playersData[id]
}

// RegisterConnection registers a pending player connection
func RegisterConnection(charid int32, charip []byte) (err error) {
	player := &ConnectingPlayer{
		ip:            charip,
		unixtimestamp: time.Now().Unix(),
	}

	if len(charip) != 4 {
		err = errors.New("ipv6 not supported")
		return
	}

	connecting[charid] = player
	return
}
