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

// MapleMapObject is a generic interface for object instances
// inside a certain map.
type MapleMapObject interface {
	ObjId() int32
	SetObjId(v int32)
	Type() MapleMapObjectType
	Pos() image.Point
	SetPos(v image.Point)
	//SendSpawnData(c *client.Connection)
	//SendDestroyData(c *client.Connection)
}
