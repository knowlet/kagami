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
	"image"
	"sync/atomic"
	"time"
)

// SpawnPoint holds information about a monster spawnpoint
type SpawnPoint struct {
	monster           *MapleMonster
	pos               image.Point
	nextPossibleSpawn int64
	mobTime           int32
	spawnedMonsters   int64
	immobile          bool
}

// NewSpawnPoint initializes a new monster spawnpoint at the given
// position and delay.
func NewSpawnPoint(smonster *MapleMonster,
	spos image.Point, smobTime int32) *SpawnPoint {

	return &SpawnPoint{
		monster:           smonster,
		pos:               spos,
		mobTime:           smobTime,
		immobile:          !smonster.Stats().Mobile(),
		nextPossibleSpawn: time.Now().UnixNano() / 1000000,
	}
}

func (s *SpawnPoint) SpawnReady() bool {
	if s.mobTime < 0 {
		return false
	}

	if (s.mobTime != 0 || s.immobile) && atomic.LoadInt64(&s.spawnedMonsters) > 0 ||
		atomic.LoadInt64(&s.spawnedMonsters) > 2 {

		return false
	}

	return s.nextPossibleSpawn <= time.Now().UnixNano()/1000000
}

// SpawnMonster forces the monster to spawn.
func (s *SpawnPoint) SpawnMonster(mapleMap *MapleMap) *MapleMonster {
	mob := CloneMapleMonster(s.monster)
	mob.SetPos(s.pos)
	atomic.AddInt64(&s.spawnedMonsters, 1)
	// TODO: set OnKilled callback
	mapleMap.SpawnMonster(mob)
	if s.mobTime == 0 {
		s.nextPossibleSpawn = time.Now().UnixNano()/1000000 + 5000
	}
	return mob
}
