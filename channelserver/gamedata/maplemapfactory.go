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
	"fmt"
	"image"
	"math/rand"
	"strconv"
	"sync"
)

import (
	"github.com/Francesco149/kagami/common/utils"
	"github.com/Francesco149/maplelib/wz"
)

const debug = false

// A MapleMapFactory is responsible for extracting and parsing data from wz files
// to initialize MapleMap objects for the desired maps.
// It also caches loaded maps to reuse them.
type MapleMapFactory struct {
	source   wz.MapleDataProvider
	nameData wz.MapleData
	maps     map[int32]*MapleMap
	mut      sync.Mutex
}

// NewMapleMapFactory initializes a new map factory from the given Map.wz and String.wz
// data sources
func NewMapleMapFactory(mapSource, stringSource wz.MapleDataProvider) (*MapleMapFactory, error) {
	tmp, err := stringSource.Get("Map.img")
	if err != nil {
		return nil, err
	}
	return &MapleMapFactory{
		source:   mapSource,
		nameData: tmp,
		maps:     make(map[int32]*MapleMap),
	}, nil
}

// DebugPrintln is a wrapper around fmt.Println that is only called when debug is enabled iternally
func DebugPrintln(a ...interface{}) {
	if debug {
		fmt.Println(a...)
	}
}

// Get looks up the given map in wz files and initializes a new map object.
// respawns, npcs and reactors determine whether those objects will be loaded or not.
// Thread safe.
func (f *MapleMapFactory) Get(mapid int32, respawns, npcs, reactors bool) *MapleMap {
	DebugPrintln(fmt.Sprintf("MapleMapFactory.Get(%v, %v, %v, %v)",
		mapid, respawns, npcs, reactors))

	f.mut.Lock()
	defer f.mut.Unlock()

	// see if the map was already previously loaded
	res := f.maps[mapid]
	if res != nil {
		return res
	}

	// img file
	mapPath := GetMapPath(mapid)
	mapData, err := f.source.Get(mapPath)
	if err != nil {
		DebugPrintln("nil mapPath")
		return nil
	}

	DebugPrintln("mapPath = ", mapPath)

	// spawn rate
	monsterRate := float32(0)
	if respawns {
		pmobrate := wz.GetFloat(mapData.ChildByPath("info/mobRate"))
		if pmobrate != nil {
			monsterRate = *pmobrate
			DebugPrintln("found info/mobRate =", monsterRate)
		}
	}
	DebugPrintln("monsterRate =", monsterRate)

	// return map
	preturnMap := wz.GetInt(mapData.ChildByPath("info/returnMap"))
	if preturnMap == nil {
		DebugPrintln("no info/returnMap")
		return nil
	}
	DebugPrintln("*preturnMap =", *preturnMap)

	// initialize the map
	res = NewMapleMap(mapid, *preturnMap, monsterRate)

	// portals
	portalFactory := NewPortalFactory()
	portalData := mapData.ChildByPath("portal")
	if portalData == nil {
		DebugPrintln("no portal")
		return nil
	}

	DebugPrintln("\nPortals:")
	for _, portal := range portalData.Children() {
		pportaltype := wz.GetInt(portal.ChildByPath("pt"))
		if pportaltype == nil {
			DebugPrintln("ignored portal with no pt")
			continue
		}

		newportal := portalFactory.Make(*pportaltype, portal)
		if newportal == nil {
			DebugPrintln("portalFactory.Make returned nil")
			return nil
		}
		tmpstr := ""

		if len(newportal.Target()) > 0 {
			tmpstr = fmt.Sprint("Target: ", newportal.Target(),
				" TargetMapId: ", newportal.TargetMapId())
		}

		DebugPrintln(newportal.Name(), "=> Type:", *pportaltype, "Id:", newportal.Id(), tmpstr)
		res.AddPortal(newportal)
	}

	// footholds and bounds
	footholds := make([]*MapleFoothold, 0)
	lbound := image.Pt(0, 0)
	ubound := image.Pt(0, 0)

	footholdData := mapData.ChildByPath("foothold")
	if footholdData == nil {
		DebugPrintln("no foothold")
		return nil
	}

	DebugPrintln("\nFootholds:")
	for _, footRoot := range footholdData.Children() {
		for _, footCat := range footRoot.Children() {
			for _, footHold := range footCat.Children() {
				px1 := wz.GetInt(footHold.ChildByPath("x1"))
				py1 := wz.GetInt(footHold.ChildByPath("y1"))
				px2 := wz.GetInt(footHold.ChildByPath("x2"))
				py2 := wz.GetInt(footHold.ChildByPath("y2"))
				pprev := wz.GetInt(footHold.ChildByPath("prev"))
				pnext := wz.GetInt(footHold.ChildByPath("next"))

				if utils.AnyNil(px1, py1, px2, py2, pprev, pnext) {
					DebugPrintln("x1, y1, x2, y2, prev or next are nil")
					DebugPrintln("values: ", px1, py1, px2, py2, pprev, pnext)
					return nil
				}

				fhid, err := strconv.Atoi(footHold.Name())
				if err != nil {
					DebugPrintln("\tmath.Atoi failed on", footHold.Name())
					return nil
				}
				fh := NewMapleFoothold(image.Pt(int(*px1), int(*py1)),
					image.Pt(int(*px2), int(*py2)), int32(fhid))
				fh.SetPrev(*pprev)
				fh.SetNext(*pnext)

				if fh.X1() < lbound.X {
					lbound.X = fh.X1()
				}

				if fh.X2() > ubound.X {
					ubound.X = fh.X2()
				}

				if fh.Y1() < lbound.Y {
					lbound.Y = fh.Y1()
				}

				if fh.Y2() > ubound.Y {
					ubound.Y = fh.Y2()
				}

				DebugPrintln(fmt.Sprintf("[%d", *px1), fmt.Sprintf("%d]", *py1),
					fmt.Sprintf("[%d", *px2), fmt.Sprintf("%d]", *py2),
					"prev:", *pprev, "next:", *pnext)
				footholds = append(footholds, fh)
			}
		}
	}

	DebugPrintln("\nlbound =", lbound)
	DebugPrintln("ubound =", ubound)

	// sort footholds in a foothold tree
	ftree := NewMapleFootholdTree(lbound, ubound)
	for fhindex, fh := range footholds {
		if fh == nil {
			DebugPrintln("nil foothold at index", fhindex)
			continue
		}
		ftree.Insert(fh)
	}
	res.SetFootholds(ftree)

	// areas (stuff like PQ platforms)
	DebugPrintln("\nAreas:")
	areaData := mapData.ChildByPath("area")
	if areaData != nil {
		for _, area := range areaData.Children() {
			px1 := wz.GetInt(area.ChildByPath("x1"))
			py1 := wz.GetInt(area.ChildByPath("y1"))
			px2 := wz.GetInt(area.ChildByPath("x2"))
			py2 := wz.GetInt(area.ChildByPath("y2"))

			if utils.AnyNil(px1, py1, px2, py2) {
				DebugPrintln("x1, y1, x2 or y2 are nil")
				DebugPrintln("values: ", px1, py1, px2, py2)
				return nil
			}

			DebugPrintln(fmt.Sprintf("[%d", *px1), fmt.Sprintf("%d]", *py1),
				fmt.Sprintf("[%d", *px2), fmt.Sprintf("%d]", *py2))
			mapArea := image.Rect(int(*px1), int(*py1), int(*px2), int(*py2))
			res.AddMapleArea(mapArea)
		}
	}

	// life entities
	DebugPrintln("\nLife:")
	lifeData := mapData.ChildByPath("life")
	if lifeData == nil {
		DebugPrintln("no life")
		return nil
	}

	for _, life := range lifeData.Children() {
		pid := wz.GetIntConvert(life.ChildByPath("id"))
		plifetype := wz.GetString(life.ChildByPath("type"))

		if utils.AnyNil(pid, plifetype) {
			DebugPrintln("pid or plifetype are nil")
			DebugPrintln("values: ", pid, plifetype)
			return nil
		}

		lifetypename := "NPC"
		if *plifetype == "m" {
			lifetypename = "Monster"
		}

		if npcs || *plifetype != "n" {
			loadedlife := loadLife(life, *pid, *plifetype)
			if loadedlife == nil {
				DebugPrintln("loadedlife is nil")
				return nil
			}

			mapleMonster, ok := loadedlife.(*MapleMonster)
			if ok {
				mobTime := wz.GetIntD(life.ChildByPath("mobTime"), 0)

				if mapleMonster.Stats().Boss() {
					// wut
					mobTime += int32(float64(mobTime) / 10.0 *
						(2.5 + 10.0*rand.Float64()))
					DebugPrintln("randomized mobTime to", mobTime)
				}

				// doesn't respawn so spawn it once immediately
				if mobTime == -1 && respawns {
					res.SpawnMonster(mapleMonster)
				}
			} else {
				//DebugPrintln("not a *MapleMonster")
				res.AddMapObject(loadedlife)
			}

			DebugPrintln(lifetypename, loadedlife.Id(), "ObjId:", loadedlife.ObjId())
		}
	}

	// reactor entities
	DebugPrintln("\nReactor:")
	reactorData := mapData.ChildByPath("reactor")
	if reactorData != nil {
		for _, reactor := range reactorData.Children() {
			pid := wz.GetIntConvert(reactor.ChildByPath("id"))
			if pid == nil {
				DebugPrintln("pid is nil, ignoring")
				continue
			}

			newreactor := loadReactor(reactor, *pid)
			res.SpawnReactor(newreactor)

			DebugPrintln(newreactor.Name(), newreactor.ReactorId(), "ObjId:", newreactor.ObjId())
		}
	}

	// strings data
	mapStringPath := GetMapStringPath(mapid)
	DebugPrintln("\nmapStringPath =", mapStringPath)

	mapStringsData := f.nameData.ChildByPath(mapStringPath)
	if mapStringsData == nil {
		DebugPrintln("no strings for this map")
		res.SetMapName("")
		res.SetStreetName("")
	} else {
		res.SetMapName(wz.GetStringD(mapStringsData.ChildByPath("mapName"), ""))
		res.SetStreetName(wz.GetStringD(mapStringsData.ChildByPath("streetName"), ""))
		DebugPrintln(res.MapName(), "/", res.StreetName())
	}

	// various map data
	res.SetClock(mapData.ChildByPath("clock") != nil)
	res.SetEverlast(mapData.ChildByPath("everlast") != nil)
	res.SetTown(mapData.ChildByPath("town") != nil)
	res.SetHPDec(wz.GetIntConvertD(mapData.ChildByPath("decHP"), 0))
	res.SetHPDecProtect(wz.GetIntConvertD(mapData.ChildByPath("protectItem"), 0))
	res.SetForcedReturnMap(wz.GetIntD(mapData.ChildByPath("info/forcedReturn"), 999999999))
	res.SetBoat(mapData.ChildByPath("shipObj") != nil)
	res.SetTimeLimit(int(wz.GetIntConvertD(mapData.ChildByPath("info/timeLimit"), -1)))

	// cache the map and return it
	f.maps[mapid] = res
	return res
}

// LoadedMapCount returns how many maps are currently cached in memory
func (f *MapleMapFactory) LoadedMapCount() int {
	return len(f.maps)
}

// IsMapLoaded checks if the given map has already been cached
func (f *MapleMapFactory) IsMapLoaded(mapid int32) bool {
	return f.maps[mapid] != nil
}

// makeLife initializes a life entity from the given data
func makeLife(id, f int32, hide bool, fh, cy, rx0, rx1, x, y int32,
	lifetype string) *AbstractLoadedMapleLife {

	loadedlife := MakeMapleLife(id, lifetype)
	loadedlife.SetCy(cy)
	loadedlife.SetF(f)
	loadedlife.SetFh(fh)
	loadedlife.SetRx0(rx0)
	loadedlife.SetRx1(rx1)
	loadedlife.SetPos(image.Pt(int(x), int(y)))
	loadedlife.SetHide(hide)
	aloadedlife, ok := loadedlife.(*AbstractLoadedMapleLife)
	if !ok {
		panic(errors.New("wtf"))
	}
	return aloadedlife
}

// loadLife loads a life entity from wz data
func loadLife(life wz.MapleData, id int32, lifetype string) IAbstractLoadedMapleLife {
	loadedlife := MakeMapleLife(id, lifetype)
	if loadedlife == nil {
		DebugPrintln("loadedlife is nil")
		return nil
	}

	pcy := wz.GetInt(life.ChildByPath("cy"))
	if pcy == nil {
		return nil
	}
	loadedlife.SetCy(*pcy)

	pf := wz.GetInt(life.ChildByPath("f"))
	if pf != nil {
		loadedlife.SetF(*pf)
	}

	pfh := wz.GetInt(life.ChildByPath("fh"))
	prx0 := wz.GetInt(life.ChildByPath("rx0"))
	prx1 := wz.GetInt(life.ChildByPath("rx1"))
	px := wz.GetInt(life.ChildByPath("x"))
	py := wz.GetInt(life.ChildByPath("y"))
	if utils.AnyNil(pfh, prx0, prx1, px, py) {
		return nil
	}

	loadedlife.SetFh(*pfh)
	loadedlife.SetRx0(*prx0)
	loadedlife.SetRx1(*prx1)
	loadedlife.SetPos(image.Pt(int(*px), int(*py)))
	loadedlife.SetHide(wz.GetIntD(life.ChildByPath("fh"), 0) > 0)

	return loadedlife
}

// loadReactor loads a reactor entoty from wz data
func loadReactor(reactor wz.MapleData, id int32) *MapleReactor {
	loadedreactor := NewMapleReactor(NewMapleReactorStats(id), id)

	px := wz.GetInt(reactor.ChildByPath("x"))
	py := wz.GetInt(reactor.ChildByPath("y"))
	pdelay := wz.GetInt(reactor.ChildByPath("reactorTime"))
	pname := wz.GetString(reactor.ChildByPath("name"))
	if utils.AnyNil(px, py, pdelay, pname) {
		return nil
	}

	loadedreactor.SetPos(image.Pt(int(*px), int(*py)))
	loadedreactor.SetDelay(int(*pdelay * 1000)) // milliseconds?
	loadedreactor.SetState(0)
	loadedreactor.SetName(*pname)

	return loadedreactor
}

// GetMapPath gets the path of the given map's img file.
func GetMapPath(mapid int32) string {
	idstring := fmt.Sprintf("%09d", mapid)
	areacode := idstring[0:1]
	return fmt.Sprint("Map/Map", areacode, "/", idstring, ".img")
}

// GetMapStringPath returns the path of the given map's strings.
func GetMapStringPath(mapid int32) string {
	continent := ""

	switch {
	case mapid < 100000000:
		continent = "maple"
	case mapid >= 100000000 && mapid < 200000000:
		continent = "victoria"
	case mapid >= 200000000 && mapid < 300000000:
		continent = "ossyria"
	case mapid >= 540000000 && mapid < 541010110:
		continent = "singapore"
	case mapid >= 600000000 && mapid < 620000000:
		continent = "MasteriaGL"
	case mapid >= 670000000 && mapid < 682000000:
		continent = "weddingGL"
	case mapid >= 682000000 && mapid < 683000000:
		continent = "HalloweenGL"
	case mapid >= 800000000 && mapid < 900000000:
		continent = "jp"
	default:
		continent = "etc"
	}

	continent += fmt.Sprintf("/%d", mapid)
	return continent
}
