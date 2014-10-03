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

// MapleMapObjectTypeCallback is a function that must return the passed
// AbstractMapleMapObject's type. See MapleMapObjectType for all possible values.
type MapleMapObjectTypeCallback func(*AbstractMapleMapObject) MapleMapObjectType

// AbstractMapleMapObject is a generic instance of an entity in a certain map of the server
type AbstractMapleMapObject struct {
	position image.Point
	objid    int32
	funcType MapleMapObjectTypeCallback
}

// NewAbstractMapleMapObject initializes a new generic entity with the given position
// and object id.
// typeCallback is a function that returns the object's type.
// See MapleMapObjectTypeCallback for more info.
func NewAbstractMapleMapObject(pos image.Point, oid int32,
	typeCallback MapleMapObjectTypeCallback) *AbstractMapleMapObject {

	return &AbstractMapleMapObject{
		position: pos,
		objid:    oid,
		funcType: typeCallback,
	}
}

// GetTypeFunc returns the object's MapleMapObjectTypeCallback.
// It is meant for internal usage when copying and embedding this struct.
func (this *AbstractMapleMapObject) GetTypeFunc() MapleMapObjectTypeCallback {
	return this.funcType
}

// Type is a wrapper that calls the object's MapleMapObjectTypeCallback.
func (this *AbstractMapleMapObject) Type() MapleMapObjectType { return this.funcType(this) }
func (this *AbstractMapleMapObject) Pos() image.Point         { return this.position }
func (this *AbstractMapleMapObject) SetPos(v image.Point)     { this.position = v }

// ObjId returns the id of this particular
// instance of the entity in this particular map
func (this *AbstractMapleMapObject) ObjId() int32 { return this.objid }

// Type sets the object's id. See ObjId for more info.
func (this *AbstractMapleMapObject) SetObjId(v int32) { this.objid = v }
