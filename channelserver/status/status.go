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

// Package status contains various information about the current status of the channel server
// such as world config, port and connections that are shared globally within the package
package status

import "sync"
import "github.com/Francesco149/kagami/common/config"

var mut sync.Mutex
var worldId int8 = -1
var chanId int8 = -1
var port int16 = 0
var worldConf *config.WorldConf = nil

// Lock locks the status mutex.
// Must be called before performing any operation on
// the channelserver status
func Lock() {
	mut.Lock()
}

// Unlock unlocks the status mutex.
func Unlock() {
	mut.Unlock()
}

func SetWorldId(wid int8)              { worldId = wid }
func WorldId() int8                    { return worldId }
func SetChanId(cid int8)               { chanId = cid }
func ChanId() int8                     { return chanId }
func SetPort(p int16)                  { port = p }
func Port() int16                      { return port }
func WorldConf() *config.WorldConf     { return worldConf }
func SetWorldConf(c *config.WorldConf) { worldConf = c }
