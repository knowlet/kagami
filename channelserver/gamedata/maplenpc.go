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

// MapleNPC holds information about a NPC.
type MapleNPC struct {
	*AbstractLoadedMapleLife
	name   string
	custom bool
}

// NewMapleNPC initializes a new NPC with the given id and name.
func NewMapleNPC(id int32, name string) *MapleNPC {
	return &MapleNPC{
		AbstractLoadedMapleLife: NewAbstractLoadedMapleLife(id,
			func(this *AbstractMapleMapObject) MapleMapObjectType {
				return NPC
			}),
		custom: false,
	}
}

func (this *MapleNPC) Name() string     { return this.name }
func (this *MapleNPC) Custom() bool     { return this.custom }
func (this *MapleNPC) SetCustom(v bool) { this.custom = v }
