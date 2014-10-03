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

// MapleMonsterStats holds all the stats for a monster.
type MapleMonsterStats struct {
	exp            int32
	hp, mp         int32
	level          int32
	removeAfter    int32
	boss           bool
	undead         bool
	ffaLoot        bool
	name           string
	animationTimes map[string]*int32
	resistance     map[Element]ElementalEffectiveness
	revives        []int32
	tagColor       byte
	tagBgColor     byte
	skills         [][2]int32
	firstAttack    bool
	buffToGive     int32
}

// NewMapleMonsterStats initializes a zeroed monster stats object.
func NewMapleMonsterStats() *MapleMonsterStats {
	return &MapleMonsterStats{
		animationTimes: make(map[string]*int32),
		resistance:     make(map[Element]ElementalEffectiveness),
		revives:        make([]int32, 0),
		skills:         make([][2]int32, 0),
	}
}

func (s *MapleMonsterStats) SetExp(v int32)         { s.exp = v }
func (s *MapleMonsterStats) SetHp(v int32)          { s.hp = v }
func (s *MapleMonsterStats) SetMp(v int32)          { s.mp = v }
func (s *MapleMonsterStats) SetLevel(v int32)       { s.level = v }
func (s *MapleMonsterStats) SetRemoveAfter(v int32) { s.removeAfter = v }
func (s *MapleMonsterStats) SetBoss(v bool)         { s.boss = v }
func (s *MapleMonsterStats) SetUndead(v bool)       { s.undead = v }
func (s *MapleMonsterStats) SetFfaLoot(v bool)      { s.ffaLoot = v }
func (s *MapleMonsterStats) SetName(v string)       { s.name = v }

func (s *MapleMonsterStats) SetAnimationTime(name string, v int32) {
	s.animationTimes[name] = &v
}

func (s *MapleMonsterStats) SetEffectiveness(e Element, ee ElementalEffectiveness) {
	s.resistance[e] = ee
}

func (s *MapleMonsterStats) SetRevives(v []int32)   { s.revives = v }
func (s *MapleMonsterStats) SetTagColor(v byte)     { s.tagColor = v }
func (s *MapleMonsterStats) SetTagBgColor(v byte)   { s.tagBgColor = v }
func (s *MapleMonsterStats) SetSkills(v [][2]int32) { s.skills = v }
func (s *MapleMonsterStats) SetFirstAttack(v bool)  { s.firstAttack = v }
func (s *MapleMonsterStats) SetBuffToGive(v int32)  { s.buffToGive = v }

func (s *MapleMonsterStats) Hp() int32    { return s.hp }
func (s *MapleMonsterStats) Mp() int32    { return s.mp }
func (s *MapleMonsterStats) Boss() bool   { return s.boss }
func (s *MapleMonsterStats) Name() string { return s.name }

func (s *MapleMonsterStats) Mobile() bool {
	return s.animationTimes["move"] != nil || s.animationTimes["fly"] != nil
}
