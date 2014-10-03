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

package gamedata

// This is a nearly 1:1 port of OdinMS' wz xml parsing, so credits to OdinMS.

import "image"

// MapleReactor holds information about a reactor entity and its state.
type MapleReactor struct {
	*AbstractMapleMapObject
	rid         int32
	stats       *MapleReactorStats
	state       int8
	delay       int
	curmap      *MapleMap
	name        string
	timerActive bool
	alive       bool
}

// NewMapleReactor initializes a reactor object with the given stats.
func NewMapleReactor(rstats *MapleReactorStats, reactid int32) *MapleReactor {
	return &MapleReactor{
		AbstractMapleMapObject: NewAbstractMapleMapObject(image.Pt(0, 0), 0,
			func(this *AbstractMapleMapObject) MapleMapObjectType {
				return REACTOR
			}),
		stats: rstats,
		rid:   reactid,
		alive: true,
	}
}

func (this *MapleReactor) SetTimerActive(v bool) {
	this.timerActive = v
}

func (this *MapleReactor) TimerActive() bool {
	return this.timerActive
}

func (this *MapleReactor) ReactorId() int32 {
	return this.rid
}

func (this *MapleReactor) State() int8 {
	return this.state
}

func (this *MapleReactor) SetState(v int8) {
	this.state = v
}

func (this *MapleReactor) SetDelay(v int) {
	this.delay = v
}

func (this *MapleReactor) Delay() int {
	return this.delay
}

func (this *MapleReactor) ReactorType() int32 {
	return this.stats.Type(this.state)
}

func (this *MapleReactor) Map() *MapleMap {
	return this.curmap
}

func (this *MapleReactor) SetMap(v *MapleMap) {
	this.curmap = v
}

func (this *MapleReactor) ReactItem() [2]int32 {
	return this.stats.ReactItem(this.state)
}

func (this *MapleReactor) Alive() bool {
	return this.alive
}

func (this *MapleReactor) Area() image.Rectangle {
	ltx := this.Pos().X + this.stats.TL().X
	lty := this.Pos().Y + this.stats.TL().Y
	rbx := this.Pos().X + this.stats.BR().X
	rby := this.Pos().Y + this.stats.BR().Y
	return image.Rect(ltx, lty, rbx, rby)
}

func (this *MapleReactor) Name() string {
	return this.name
}

func (this *MapleReactor) SetName(v string) {
	this.name = v
}
