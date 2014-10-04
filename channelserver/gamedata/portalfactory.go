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

import (
	"errors"
	"image"
	"strconv"
)

import (
	"github.com/Francesco149/kagami/common/utils"
	"github.com/Francesco149/maplelib/wz"
)

// PortalFactory is responsible for loading portal data from wz files
// and returning new MaplePortal objects.
type PortalFactory struct {
	nextDoorPortal int32
}

// NewPortalFactory initializes a new portal factory.
func NewPortalFactory() *PortalFactory {
	return &PortalFactory{0x80}
}

// Make loads a generic portal from wz data.
func (f *PortalFactory) Make(portaltype int32, portal wz.MapleData) MaplePortal {
	//DebugPrintln(fmt.Sprintf("PortalFactory.Make(%v, %v)", portaltype, portal))
	var res IMapleGenericPortal = nil
	if portaltype == MAP_PORTAL {
		res = NewMapleMapPortal()
	} else {
		res = NewMapleGenericPortal(portaltype)
	}
	err := f.loadPortal(res, portal)
	if err != nil {
		DebugPrintln("loadPortal returned:", err)
		return nil
	}
	return res
}

func (f *PortalFactory) loadPortal(portal IMapleGenericPortal, portalData wz.MapleData) (err error) {
	pname := wz.GetString(portalData.ChildByPath("pn"))
	ptargetname := wz.GetString(portalData.ChildByPath("tn"))
	ptargetid := wz.GetInt(portalData.ChildByPath("tm"))
	px := wz.GetInt(portalData.ChildByPath("x"))
	py := wz.GetInt(portalData.ChildByPath("y"))
	if utils.AnyNil(pname, ptargetname, ptargetid, px, py) {
		err = errors.New("found nil data")
		return
	}

	portal.SetName(*pname)
	portal.SetTarget(*ptargetname)
	portal.SetTargetMapId(*ptargetid)
	portal.SetPos(image.Pt(int(*px), int(*py)))
	portal.SetScriptName(wz.GetStringD(portalData.ChildByPath("script"), ""))

	if portal.Type() == DOOR_PORTAL {
		portal.SetId(f.nextDoorPortal)
		f.nextDoorPortal++
	} else {
		var val int
		val, err = strconv.Atoi(portalData.Name())
		if err != nil {
			return
		}

		portal.SetId(int32(val))
	}

	return
}
