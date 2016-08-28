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

import (
	"fmt"
	"time"
)

import (
	"github.com/knowlet/kagami/common/utils"
	"github.com/Francesco149/maplelib/wz"
)

var mapWz, mobWz, itemWz, characterWz, stringWz, reactorWz wz.MapleDataProvider
var cashStringData, consumeStringData, eqpStringData,
	etcStringData, insStringData, petStringData,
	mobStringData, npcStringData, mapStringData wz.MapleData

func GetMapWz() wz.MapleDataProvider       { return mapWz }
func GetMobWz() wz.MapleDataProvider       { return mobWz }
func GetItemWz() wz.MapleDataProvider      { return itemWz }
func GetCharacterWz() wz.MapleDataProvider { return characterWz }
func GetStringWz() wz.MapleDataProvider    { return stringWz }
func GetReactorWz() wz.MapleDataProvider   { return reactorWz }

func GetCashStringImg() wz.MapleData    { return cashStringData }
func GetConsumeStringImg() wz.MapleData { return consumeStringData }
func GetEqpStringImg() wz.MapleData     { return eqpStringData }
func GetEtcStringImg() wz.MapleData     { return etcStringData }
func GetInsStringImg() wz.MapleData     { return insStringData }
func GetPetStringImg() wz.MapleData     { return petStringData }
func GetMobStringImg() wz.MapleData     { return mobStringData }
func GetNpcStringImg() wz.MapleData     { return npcStringData }
func GetMapStringImg() wz.MapleData     { return mapStringData }

// InitProviders initializes all of the wz data providers.
func InitProviders() (err error) {
	fmt.Println("Loading wz files...")
	starttime := time.Now()
	var wzLoad = map[*wz.MapleDataProvider]string{
		&mapWz:       "wz/Map.wz",
		&mobWz:       "wz/Mob.wz",
		&itemWz:      "wz/Item.wz",
		&characterWz: "wz/Character.wz",
		&stringWz:    "wz/String.wz",
		&reactorWz:   "wz/Reactor.wz",
	}

	for pprovider, wzpath := range wzLoad {
		fmt.Println(wzpath)
		*pprovider, err = wz.NewMapleDataProvider(wzpath)
		if err != nil {
			return
		}
	}
	fmt.Printf("Done! (%v)\n\n", time.Since(starttime))

	starttime = time.Now()
	fmt.Println("Loading img directories...")
	var imgLoad = map[*wz.MapleData]*utils.Pair{
		&cashStringData:    &utils.Pair{stringWz, "Cash.img"},
		&consumeStringData: &utils.Pair{stringWz, "Consume.img"},
		&eqpStringData:     &utils.Pair{stringWz, "Eqp.img"},
		&etcStringData:     &utils.Pair{stringWz, "Etc.img"},
		&insStringData:     &utils.Pair{stringWz, "Ins.img"},
		&petStringData:     &utils.Pair{stringWz, "Pet.img"},
		&mobStringData:     &utils.Pair{stringWz, "Mob.img"},
		&npcStringData:     &utils.Pair{stringWz, "Npc.img"},
		&mapStringData:     &utils.Pair{stringWz, "Map.img"},
	}

	for pdata, pair := range imgLoad {
		fmt.Println(pair.First.(wz.MapleDataProvider).Root().Name(), "->", pair.Second.(string))
		*pdata, err = pair.First.(wz.MapleDataProvider).Get(pair.Second.(string))
		if err != nil {
			return
		}
	}
	fmt.Printf("Done! (%v)\n\n", time.Since(starttime))

	return
}
