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
	"fmt"
	"image"
	"math"
	"sync"
	"sync/atomic"
)

// MAX_OID is the maximum allowed object id
const MAX_OID = 20000

// all possible maple object types as an array so that they're iterable
var rangedMapobjectTypes = []MapleMapObjectType{ITEM, MONSTER,
	DOOR, SUMMON, REACTOR}

// MapleMap holds all the data for a single instance of a maplestory map and
// manages the map's objects and current state.
type MapleMap struct {
	objects         map[int32]MapleMapObject
	monsterSpawns   []*SpawnPoint
	spawnedMonsters int64
	portals         map[int32]MaplePortal
	areas           []image.Rectangle
	footholds       *MapleFootholdTree
	mapid           int32
	runningOid      int32
	returnMapId     int32
	monsterRate     float32
	dropsDisabled   bool
	clock           bool
	boat            bool
	docked          bool
	mapName         string
	streetName      string
	//mapEffect *MapleMapEffect
	everlast        bool
	forcedReturnMap int32
	timeLimit       int
	//mapTimer *MapleMapTimer
	dropLife    int
	decHP       int32
	protectItem int32
	town        bool
	mut         sync.Mutex
}

// NewMapleMap initializes a new map
func NewMapleMap(mmapid, mreturnMapId int32, mmonsterRate float32) *MapleMap {
	res := &MapleMap{
		objects:         make(map[int32]MapleMapObject),
		monsterSpawns:   make([]*SpawnPoint, 0),
		spawnedMonsters: 0,
		portals:         make(map[int32]MaplePortal),
		areas:           make([]image.Rectangle, 0),
		footholds:       nil,
		mapid:           mmapid,
		runningOid:      100,
		returnMapId:     mreturnMapId,
		dropsDisabled:   false,
		//mapEffect: nil,
		everlast:        false,
		forcedReturnMap: 999999999,
		//mapTimer: nil,
		dropLife:    180000,
		decHP:       0,
		protectItem: 0,
	}

	if mmonsterRate > 0 {
		// ??????????
		res.monsterRate = mmonsterRate
		greaterThanOne := mmonsterRate > 1.0
		res.monsterRate = float32(math.Abs(1.0 - float64(res.monsterRate)))
		res.monsterRate /= 2.0

		if greaterThanOne {
			res.monsterRate = 1.0 + res.monsterRate
		} else {
			res.monsterRate = 1.0 - res.monsterRate
		}

		// TODO: spawn mob respawn thread
	}

	return res
}

func (this *MapleMap) AddPortal(p MaplePortal) {
	this.portals[p.Id()] = p
}

func (this *MapleMap) SetFootholds(f *MapleFootholdTree) {
	this.footholds = f
}

func (this *MapleMap) AddMapleArea(a image.Rectangle) {
	this.areas = append(this.areas, a)
}

func (this *MapleMap) SpawnMonster(m *MapleMonster) {
	m.SetMap(this)
	// TODO
	atomic.AddInt64(&this.spawnedMonsters, 1)
}

func (this *MapleMap) SpawnReactor(r *MapleReactor) {
	r.SetMap(this)
	// TODO
}

func (this *MapleMap) SetMapName(v string) {
	this.mapName = v
}

func (this *MapleMap) SetStreetName(v string) {
	this.streetName = v
}

func (this *MapleMap) MapName() string {
	return this.mapName
}

func (this *MapleMap) StreetName() string {
	return this.streetName
}

func (this *MapleMap) SetClock(v bool)            { this.clock = v }
func (this *MapleMap) SetEverlast(v bool)         { this.everlast = v }
func (this *MapleMap) SetTown(v bool)             { this.town = v }
func (this *MapleMap) SetHPDec(v int32)           { this.decHP = v }
func (this *MapleMap) SetHPDecProtect(v int32)    { this.protectItem = v }
func (this *MapleMap) SetForcedReturnMap(v int32) { this.forcedReturnMap = v }
func (this *MapleMap) SetBoat(v bool)             { this.boat = v }
func (this *MapleMap) SetTimeLimit(v int)         { this.timeLimit = v }

func (this *MapleMap) AddMapObject(mapobj MapleMapObject) {
	this.mut.Lock() // thread safety
	defer this.mut.Unlock()

	mapobj.SetObjId(this.runningOid)
	this.objects[this.runningOid] = mapobj
	this.incrementRunningOid()
}

func (this *MapleMap) incrementRunningOid() {
	this.runningOid++

	for i := 1; i < MAX_OID; i++ {
		if this.runningOid > MAX_OID {
			this.runningOid = 100
		}

		if this.objects[this.runningOid] != nil {
			this.runningOid++
		} else {
			return
		}
	}

	fmt.Println("Out of OIDs on map", this.mapid)
}
