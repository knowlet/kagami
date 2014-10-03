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

// MapleMonster holds information about a maplestory monster.
type MapleMonster struct {
	*AbstractLoadedMapleLife
	stats           *MapleMonsterStats
	overrideStats   *MapleMonsterStats
	hp              int32
	mp              int32
	curmap          *MapleMap
	venomMultiplier int32
	fake            bool
	dropsDisabled   bool
}

// NewMapleMonster initializes a monster with the given stats.
func NewMapleMonster(id int32, stats *MapleMonsterStats) *MapleMonster {
	res := &MapleMonster{
		AbstractLoadedMapleLife: NewAbstractLoadedMapleLife(id,
			func(this *AbstractMapleMapObject) MapleMapObjectType {
				return MONSTER
			}),
	}

	res.SetStance(5)
	res.stats = stats
	res.hp = stats.Hp()
	res.mp = stats.Mp()

	return res
}

// CloneMapleMonster returns a copy of the given monster.
func CloneMapleMonster(monster *MapleMonster) *MapleMonster {
	return NewMapleMonster(monster.id, monster.stats)
}

func (this *MapleMonster) Stats() *MapleMonsterStats {
	return this.stats
}

func (this *MapleMonster) SetMap(v *MapleMap) {
	this.curmap = v
}
