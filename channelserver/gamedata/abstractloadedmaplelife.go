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
	"strconv"
	"strings"
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/maplelib/wz"
)

// NOTE: this doesn't need thread safety because it is only called when loading
// maps and maps are thread-safe already

var initialized = false
var data, stringDataWZ wz.MapleDataProvider = nil, nil
var mobStringData, npcStringData wz.MapleData = nil, nil
var monsterStats = make(map[int32]*MapleMonsterStats)

// IAbstractLoadedMapleLife is a generic interface to access the data of
// Life entities that have been parsed from wz / wz xml files or mcdb.
type IAbstractLoadedMapleLife interface {
	F() int32
	SetF(v int32)
	Hidden() bool
	SetHide(v bool)
	Fh() int32     // Fh returns the entity's foothold
	SetFh(v int32) // SetFh sets the entity's foothold
	StartFh() int32
	SetStartFh(v int32)
	Cy() int32
	SetCy(v int32)
	Rx0() int32
	SetRx0(v int32)
	Rx1() int32
	SetRx1(v int32)
	Id() int32

	// ObjId returns the id of this particular
	// instance of the entity in this particular map
	ObjId() int32

	// Type sets the object's id. See ObjId for more info.
	SetObjId(v int32)

	// Type returns this entity's type.
	// See MapleMapObjectType for a list of all possible types.
	Type() MapleMapObjectType

	Pos() image.Point
	SetPos(v image.Point)
}

// AbstractLoadedMapleLife is a generic maplestory life entity such as a monster or an npc.
// See IAbstractLoadedMapleLife for more info on the getters.
type AbstractLoadedMapleLife struct {
	*AbstractAnimatedMapleMapObject
	id      int32
	f       int32
	hide    bool
	fh      int32
	startFh int32
	cy      int32
	rx0     int32
	rx1     int32
}

// NewAbstractLoadedMapleLife initializes a generic life entity and returns it
// lifeid is the wz id of the entity (NOT the object id!)
// typeCallback is a function that returns the object's type.
// See MapleMapObjectTypeCallback for more info.
func NewAbstractLoadedMapleLife(lifeid int32,
	typeCallback MapleMapObjectTypeCallback) *AbstractLoadedMapleLife {

	return &AbstractLoadedMapleLife{
		AbstractAnimatedMapleMapObject: NewAbstractAnimatedMapleMapObject(image.Pt(0, 0), 0, typeCallback),
		id: lifeid,
	}
}

// CloneAbstractLoadedMapleLife returns a copy of the given life entity
func CloneAbstractLoadedMapleLife(life *AbstractLoadedMapleLife) *AbstractLoadedMapleLife {
	res := NewAbstractLoadedMapleLife(life.Id(), life.GetTypeFunc())
	res.f = life.f
	res.hide = life.hide
	res.fh = life.fh
	res.startFh = life.startFh
	res.cy = life.cy
	res.rx0 = life.rx0
	res.rx1 = life.rx1
	return res
}

func (l *AbstractLoadedMapleLife) F() int32           { return l.f }
func (l *AbstractLoadedMapleLife) SetF(v int32)       { l.f = v }
func (l *AbstractLoadedMapleLife) Hidden() bool       { return l.hide }
func (l *AbstractLoadedMapleLife) SetHide(v bool)     { l.hide = v }
func (l *AbstractLoadedMapleLife) Fh() int32          { return l.fh }
func (l *AbstractLoadedMapleLife) SetFh(v int32)      { l.fh = v }
func (l *AbstractLoadedMapleLife) StartFh() int32     { return l.startFh }
func (l *AbstractLoadedMapleLife) SetStartFh(v int32) { l.startFh = v }
func (l *AbstractLoadedMapleLife) Cy() int32          { return l.cy }
func (l *AbstractLoadedMapleLife) SetCy(v int32)      { l.cy = v }
func (l *AbstractLoadedMapleLife) Rx0() int32         { return l.rx0 }
func (l *AbstractLoadedMapleLife) SetRx0(v int32)     { l.rx0 = v }
func (l *AbstractLoadedMapleLife) Rx1() int32         { return l.rx1 }
func (l *AbstractLoadedMapleLife) SetRx1(v int32)     { l.rx1 = v }
func (l *AbstractLoadedMapleLife) Id() int32          { return l.id }

func initIfNotInitialized() (err error) {
	if initialized {
		return
	}

	data, err = wz.NewMapleDataProvider("wz/Mob.wz")
	if err != nil {
		return
	}

	stringDataWZ, err = wz.NewMapleDataProvider("wz/String.wz")
	if err != nil {
		return
	}

	mobStringData, err = stringDataWZ.Get("Mob.img")
	if err != nil {
		return
	}

	npcStringData, err = stringDataWZ.Get("Npc.img")
	if err != nil {
		return
	}

	if mobStringData == nil || npcStringData == nil {
		err = errors.New("mobStringData or npcStringData are nil")
	}

	initialized = true

	return
}

// MakeMapleLife looks up the given life entity in wz files and returns
// a generic interface to the entity's data.
// If the entity is not found or invalid / not supported, nil will be returned.
func MakeMapleLife(id int32, lifetype string) IAbstractLoadedMapleLife {
	err := initIfNotInitialized()
	if err != nil {
		return nil
	}

	switch strings.ToLower(lifetype) {
	case "n": // NPC
		npc := getNPC(id)
		if npc == nil {
			return nil
		}
		return npc

	case "m": // Monster
		mob := getMonster(id)
		if mob == nil {
			return nil
		}
		return mob

	default: // Invalid
		return nil
	}

	return nil
}

func getMonster(id int32) *MapleMonster {
	// see if the monster had already been previously loaded
	stats := monsterStats[id]
	if stats != nil {
		// return the cached mob
		return NewMapleMonster(id, stats)
	}

	// nope, the monster needs to be loaded

	// the id must be zero-padded to 7 digits
	monsterData, err := data.Get(fmt.Sprintf("%07d", id) + ".img")
	if err != nil {
		DebugPrintln(err)
		return nil
	}

	// various mob data
	monsterInfoData := monsterData.ChildByPath("info")
	stats = NewMapleMonsterStats()
	php := wz.GetIntConvert(monsterInfoData.ChildByPath("maxHP"))
	mp := wz.GetIntConvertD(monsterInfoData.ChildByPath("maxMP"), 0)
	exp := wz.GetIntConvertD(monsterInfoData.ChildByPath("exp"), 0)
	plevel := wz.GetIntConvert(monsterInfoData.ChildByPath("level"))
	removeafter := wz.GetIntConvertD(monsterInfoData.ChildByPath("removeAfter"), 0)
	boss := wz.GetIntConvertD(monsterInfoData.ChildByPath("boss"), 0)
	ffaloot := wz.GetIntConvertD(monsterInfoData.ChildByPath("publicReward"), 0)
	undead := wz.GetIntConvertD(monsterInfoData.ChildByPath("undead"), 0)

	// check all pointers
	if common.AnyNil(php, plevel) {
		return nil
	}

	// we can now safely assign all the retrieved data
	stats.SetHp(*php)
	stats.SetMp(mp)
	stats.SetExp(exp)
	stats.SetLevel(*plevel)
	stats.SetRemoveAfter(removeafter)
	stats.SetBoss(boss > 0)
	stats.SetFfaLoot(ffaloot > 0)
	stats.SetUndead(undead > 0)

	// name & buff
	stats.SetName(wz.GetStringD(mobStringData.ChildByPath(fmt.Sprintf("%d/name", id)),
		"LOLI 404 CHEST NOT FOUND"))
	stats.SetBuffToGive(wz.GetIntConvertD(monsterInfoData.ChildByPath("buff"), -1))

	// first attack info
	firstAttackData := monsterInfoData.ChildByPath("firstAttack")
	firstAttack := int32(0)

	if firstAttackData != nil {
		if firstAttackData.Type() == wz.FLOAT {
			pfirstAttack := wz.GetFloat(firstAttackData)
			if pfirstAttack == nil {
				return nil
			}

			firstAttack = int32(*pfirstAttack)
		} else {
			pfirstAttack := wz.GetInt(firstAttackData)
			if pfirstAttack == nil {
				return nil
			}

			firstAttack = *pfirstAttack
		}
	}

	stats.SetFirstAttack(firstAttack > 0)

	// 8810018 is HT
	if stats.Boss() || id == 8810018 {
		// if the mob is a boss, retrieve the boss hp bar info
		hpTagColor := monsterInfoData.ChildByPath("hpTagColor")
		hpTagBgColor := monsterInfoData.ChildByPath("hpTagBgcolor")

		// these values default to zero, so if they're missing or invalid
		// the boss will still work but without hp bars
		stats.SetTagColor(byte(wz.GetIntConvertD(hpTagColor, 0)))
		stats.SetTagBgColor(byte(wz.GetIntConvertD(hpTagBgColor, 0)))
	}

	// all data that isn't info is animations
	for _, idata := range monsterData.Children() {
		if idata.Name() == "info" {
			continue
		}

		delay := int32(0)
		for _, pic := range idata.Children() {
			delay += wz.GetIntConvertD(pic.ChildByPath("delay"), 0)
		}
		stats.SetAnimationTime(idata.Name(), delay)
	}

	// revive info, not sure what this does yet
	reviveInfo := monsterInfoData.ChildByPath("revive")
	if reviveInfo != nil {
		revives := make([]int32, 0)
		for _, revivedata := range reviveInfo.Children() {
			pval := wz.GetInt(revivedata)
			if pval == nil {
				continue
			}

			revives = append(revives, *pval)
		}
		stats.SetRevives(revives)
	}

	// elemental weakness / strength
	decodeElementalString(stats,
		wz.GetStringD(monsterInfoData.ChildByPath("elemAttr"), ""))

	// monster skills
	monsterSkillData := monsterInfoData.ChildByPath("skill")
	if monsterSkillData != nil {
		i := 0
		skills := make([][2]int32, 0)
		for monsterSkillData.ChildByPath(fmt.Sprintf("%d", i)) != nil {
			skills = append(skills, [2]int32{
				wz.GetIntD(monsterSkillData.ChildByPath(fmt.Sprintf("%d/skill", i)), 0),
				wz.GetIntD(monsterSkillData.ChildByPath(fmt.Sprintf("%d/level", i)), 0),
			})
			i++
		}
		stats.SetSkills(skills)
	}

	// cache the monster for future usages and return it
	monsterStats[id] = stats
	return NewMapleMonster(id, stats)
}

func getNPC(id int32) *MapleNPC {
	// NPC's only need strings data
	return NewMapleNPC(id, wz.GetStringD(
		npcStringData.ChildByPath(fmt.Sprintf("%d/name", id)),
		"LOLI 404 CHEST NOT FOUND"))
}

func decodeElementalString(dststats *MapleMonsterStats, elemAttr string) (err error) {
	// the format of elemental weakness strings is "ABABABAB"
	// where A = element and B = weakness / strength
	for i := 0; i < len(elemAttr); i += 2 {
		e := ElementFromChar(elemAttr[i : i+1]) // first char
		var number int
		number, err = strconv.Atoi(elemAttr[i+1 : i+2]) // 2nd char
		if err != nil {
			return
		}
		ee := ElementalEffectivenessFromNumber(int32(number))
		dststats.SetEffectiveness(e, ee)
	}
	return
}
