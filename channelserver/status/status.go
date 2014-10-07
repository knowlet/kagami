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

import (
	"github.com/Francesco149/kagami/channelserver/gamedata"
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/config"
)

var Get = make(chan *Status, 1)

type Status struct {
	worldId, chanId int8
	port            int16
	worldConf       *config.WorldConf
	worldConn       *common.InterserverClient
	loginConn       *common.InterserverClient
	mapFactory      *gamedata.MapleMapFactory
}

func (this *Status) WorldId() int8                         { return this.worldId }
func (this *Status) ChanId() int8                          { return this.chanId }
func (this *Status) Port() int16                           { return this.port }
func (this *Status) WorldConf() *config.WorldConf          { return this.worldConf }
func (this *Status) WorldConn() *common.InterserverClient  { return this.worldConn }
func (this *Status) LoginConn() *common.InterserverClient  { return this.loginConn }
func (this *Status) MapFactory() *gamedata.MapleMapFactory { return this.mapFactory }

func (this *Status) SetWorldId(v int8)                         { this.worldId = v }
func (this *Status) SetChanId(v int8)                          { this.chanId = v }
func (this *Status) SetPort(v int16)                           { this.port = v }
func (this *Status) SetWorldConf(v *config.WorldConf)          { this.worldConf = v }
func (this *Status) SetWorldConn(v *common.InterserverClient)  { this.worldConn = v }
func (this *Status) SetLoginConn(v *common.InterserverClient)  { this.loginConn = v }
func (this *Status) SetMapFactory(v *gamedata.MapleMapFactory) { this.mapFactory = v }

func Init() {
	Get <- &Status{
		worldId:    -1,
		chanId:     -1,
		port:       -1,
		worldConf:  nil,
		worldConn:  nil,
		loginConn:  nil,
		mapFactory: nil,
	}
}
