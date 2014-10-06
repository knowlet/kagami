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

// Package players manages the character pool for the channelserver
package players

import "sync"
import "github.com/Francesco149/kagami/channelserver/client"

type ClientOperationCallback func(*client.Connection) error

var mut sync.Mutex
var pendingIps = make(map[int32][]byte) // pending connections mapped by charid
var characters = make(map[int32]*client.Connection)

// TODO: add a timeout on pendingIps entries so that they don't clog up the memory
// when players fail to connect to the channelserver for whatever reason

// Lock locks the player pool mutex.
// Must be called before performing any operation on
// the channelserver player pool
func Lock() {
	mut.Lock()
}

// Unlock unlocks the player pool mutex.
func Unlock() {
	mut.Unlock()
}

func PendingIp(charid int32) []byte {
	return pendingIps[charid]
}

func AddPendingIp(charid int32, ip []byte) {
	pendingIps[charid] = ip
}

func RemovePendingIp(charid int32) {
	delete(pendingIps, charid)
}

func Add(con *client.Connection) {
	characters[con.Stats().Id()] = con
}

func Remove(con *client.Connection) {
	delete(characters, con.Stats().Id())
}

// Execute calls the given callback on all clients in the player pool.
// See ClientOperationCallback for the callback signature.
func Execute(fn ClientOperationCallback) (err error) {
	for _, client := range characters {
		err = fn(client)
		if err != nil {
			return
		}
	}
	return
}
