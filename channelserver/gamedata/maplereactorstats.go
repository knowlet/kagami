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
)

import (
	"github.com/Francesco149/kagami/common/utils"
	"github.com/Francesco149/maplelib/wz"
)

var reactorData wz.MapleDataProvider = nil
var reactorStats = make(map[int32]*MapleReactorStats)

// stateData contains information about a MapleReactor's status
type stateData struct {
	typ       int32
	reactItem [2]int32
	nextState int8
}

// MapleReactorStats holds information about a reactor and its status.
type MapleReactorStats struct {
	tl, br    image.Point
	stateInfo map[int8]*stateData
}

// NewMapleReactorStats retrieves data for the given reactor from wz files
// and caches it for future usage.
func NewMapleReactorStats(reactId int32) *MapleReactorStats {
	// initialize wz source if it isn't already
	if reactorData == nil {
		var err error
		reactorData, err = wz.NewMapleDataProvider("wz/Reactor.wz")
		if err != nil {
			return nil
		}
	}

	// see if the map is already cached
	stats := reactorStats[reactId]
	if stats != nil {
		return stats
	}

	// img file
	infoId := reactId
	reactorImg, err := reactorData.Get(fmt.Sprintf("%09d", infoId) + ".img")
	if err != nil {
		return nil
	}

	// from what I can see, link redirects the reactor to another one
	link := reactorImg.ChildByPath("info/link")
	if link != nil {
		pinfoId := wz.GetIntConvert(link)
		if pinfoId != nil {
			infoId = *pinfoId
			stats = reactorStats[infoId]
		}
	}

	// see if the linked reactor was already cached
	if stats != nil {
		return stats
	}

	// img file of the linked reactor
	reactorImg, err = reactorData.Get(fmt.Sprintf("%09d", infoId) + ".img")
	if err != nil {
		return nil
	}

	// events
	reactorInfoData := reactorImg.ChildByPath("0/event/0")
	stats = &MapleReactorStats{}

	if reactorInfoData == nil {
		// static reactor
		stats.AddState(0, 999, [2]int32{-1, -1}, 0) // TODO: check if -1 is okay as a nil value
	} else {
		areaSet := false
		i := 0
		for reactorInfoData != nil {
			var reactItem = [2]int32{-1, -1}
			ptype := wz.GetIntConvert(reactorInfoData.ChildByPath("type"))
			if ptype == nil {
				return nil
			}

			if *ptype == 100 {
				pfirst := wz.GetIntConvert(reactorInfoData.ChildByPath("0"))
				psecond := wz.GetIntConvert(reactorInfoData.ChildByPath("1"))
				if utils.AnyNil(pfirst, psecond) {
					return nil
				}

				reactItem = [2]int32{*pfirst, *psecond}

				if !areaSet {
					plt := wz.GetPoint(reactorInfoData.ChildByPath("lt"))
					prb := wz.GetPoint(reactorInfoData.ChildByPath("rb"))
					if utils.AnyNil(plt, prb) {
						return nil
					}

					stats.SetTL(*plt)
					stats.SetBR(*prb)
					areaSet = true
				}
			}

			pnextState := wz.GetIntConvert(reactorInfoData.ChildByPath("state"))
			if pnextState == nil {
				return nil
			}

			stats.AddState(int8(i), *ptype, reactItem, int8(*pnextState))
			i++
			reactorInfoData = reactorImg.ChildByPath(fmt.Sprintf("%d/event/0", i))
		}
	}

	// cache and return the reactor
	reactorStats[infoId] = stats
	if reactId != infoId {
		// cache and return the linked reactor
		reactorStats[reactId] = stats
	}

	return stats
}

func (this *MapleReactorStats) SetTL(v image.Point) { this.tl = v }
func (this *MapleReactorStats) SetBR(v image.Point) { this.br = v }
func (this *MapleReactorStats) TL() image.Point     { return this.tl }
func (this *MapleReactorStats) BR() image.Point     { return this.br }

func (this *MapleReactorStats) AddState(state int8, styp int32,
	sreactItem [2]int32, snextState int8) {

	this.stateInfo[state] = &stateData{typ: styp, reactItem: sreactItem,
		nextState: snextState}
}

func (this *MapleReactorStats) NextState(state int8) int8 {
	if this.stateInfo[state] != nil {
		return this.stateInfo[state].nextState
	} else {
		return -1
	}
}

func (this *MapleReactorStats) Type(state int8) int32 {
	if this.stateInfo[state] != nil {
		return this.stateInfo[state].typ
	} else {
		return -1
	}
}

func (this *MapleReactorStats) ReactItem(state int8) [2]int32 {
	if this.stateInfo[state] != nil {
		return this.stateInfo[state].reactItem
	} else {
		return [2]int32{-1, -1}
	}
}
