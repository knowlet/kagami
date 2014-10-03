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

// Possible portal types
const (
	MAP_PORTAL  = 2
	DOOR_PORTAL = 6
)

// Possible portal statuses
const (
	OPEN   = true
	CLOSED = false
)

// MaplePortal is a generic interface for portals
type MaplePortal interface {
	Type() int32 // Type returns the portal type (MAP_PORTAL, DOOR_PORTAL)
	Id() int32
	Pos() image.Point
	Name() string
	Target() string
	ScriptName() string
	SetScriptName(v string)
	SetStatus(v bool) // SetStatus sets the portal's status (OPEN/CLOSED)
	Status() bool     // Status returns the portal's status (OPEN/CLOSED)
	TargetMapId() int32
	//Enter(c *client.Connection) error
	SetState(v bool)
	State() bool
}
