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

// Package gamedata contains various utilities to parse and manage game data
// such as wz map data, skills, monsters and so on
// This is a nearly 1:1 port of OdinMS' wz xml parsing, so credits to OdinMS.
package gamedata

import "image"

// AbstractAnimatedMapleMapObject is a generic maplestory entity that
// has sprite animations and can have a stance.
type AbstractAnimatedMapleMapObject struct {
	*AbstractMapleMapObject
	stance int32
}

// NewAbstractAnimatedMapleMapObject initializes a generic maplestory animated entity
// with the given position and object id.
// typeCallback is a function that returns the object's type.
// See MapleMapObjectTypeCallback for more info.
func NewAbstractAnimatedMapleMapObject(pos image.Point, oid int32,
	typeCallback MapleMapObjectTypeCallback) *AbstractAnimatedMapleMapObject {

	return &AbstractAnimatedMapleMapObject{
		AbstractMapleMapObject: NewAbstractMapleMapObject(pos, oid, typeCallback),
	}
}

func (o *AbstractAnimatedMapleMapObject) Stance() int32     { return o.stance }
func (o *AbstractAnimatedMapleMapObject) FacingLeft() bool  { return o.stance%2 == 1 }
func (o *AbstractAnimatedMapleMapObject) SetStance(v int32) { o.stance = v }
