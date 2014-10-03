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

// IMapleGenericPortal is a generic interface for maplestory portals
type IMapleGenericPortal interface {
	Id() int32     // Id returns the portal id (NOT the object id!)
	SetId(v int32) // Id sets the portal id (NOT the object id!)
	Name() string
	Target() string
	SetStatus(v bool)
	Status() bool
	TargetMapId() int32
	Type() int32
	ScriptName() string
	Pos() image.Point
	SetName(v string)
	SetPos(pos image.Point)
	SetTarget(v string)
	SetTargetMapId(v int32)
	SetScriptName(v string)
	SetState(v bool)
	State() bool
}

// MapleGenericPortal is a generic maplestory portal
// See IMapleGenericPortal for more informations about getters and setters.
type MapleGenericPortal struct {
	name        string
	target      string
	position    image.Point
	targetmap   int32
	typ         int32
	status      bool
	pid         int32
	scriptName  string
	portalState bool
}

// NewMapleGenericPortal initializes a new generic portal of the given type
func NewMapleGenericPortal(portaltype int32) *MapleGenericPortal {
	return &MapleGenericPortal{
		typ:    portaltype,
		status: true,
	}
}

// A MapleMapPortal is a map portal
// See IMapleGenericPortal for more informations about getters and setters.
type MapleMapPortal struct {
	*MapleGenericPortal
}

// NewMapleMapPortal initializes a new map portal
func NewMapleMapPortal() *MapleMapPortal {
	return &MapleMapPortal{NewMapleGenericPortal(MAP_PORTAL)}
}

func (this *MapleGenericPortal) Id() int32              { return this.pid }
func (this *MapleGenericPortal) SetId(v int32)          { this.pid = v }
func (this *MapleGenericPortal) Name() string           { return this.name }
func (this *MapleGenericPortal) Target() string         { return this.target }
func (this *MapleGenericPortal) SetStatus(v bool)       { this.status = v }
func (this *MapleGenericPortal) Status() bool           { return this.status }
func (this *MapleGenericPortal) TargetMapId() int32     { return this.targetmap }
func (this *MapleGenericPortal) Type() int32            { return this.typ }
func (this *MapleGenericPortal) ScriptName() string     { return this.scriptName }
func (this *MapleGenericPortal) Pos() image.Point       { return this.position }
func (this *MapleGenericPortal) SetName(v string)       { this.name = v }
func (this *MapleGenericPortal) SetPos(pos image.Point) { this.position = pos }
func (this *MapleGenericPortal) SetTarget(v string)     { this.target = v }
func (this *MapleGenericPortal) SetTargetMapId(v int32) { this.targetmap = v }
func (this *MapleGenericPortal) SetScriptName(v string) { this.scriptName = v }

func (this *MapleGenericPortal) SetState(v bool) { this.portalState = v }
func (this *MapleGenericPortal) State() bool     { return this.portalState }
