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

// MapleFoothold holds information about a maplestory map foothold
type MapleFoothold struct {
	p1, p2     image.Point
	id         int32
	next, prev int32
}

// NewMapleFoothold initializes a foothold with the given lt/rb points and foothold id
func NewMapleFoothold(fhp1, fhp2 image.Point, fhid int32) *MapleFoothold {
	return &MapleFoothold{
		p1: fhp1,
		p2: fhp2,
		id: fhid,
	}
}

func (this *MapleFoothold) Wall() bool {
	return this.p1.X == this.p2.X
}

func (this *MapleFoothold) X1() int         { return this.p1.X }
func (this *MapleFoothold) X2() int         { return this.p2.X }
func (this *MapleFoothold) Y1() int         { return this.p1.Y }
func (this *MapleFoothold) Y2() int         { return this.p2.Y }
func (this *MapleFoothold) Id() int32       { return this.id }
func (this *MapleFoothold) Next() int32     { return this.next }
func (this *MapleFoothold) Prev() int32     { return this.prev }
func (this *MapleFoothold) SetNext(v int32) { this.next = v }
func (this *MapleFoothold) SetPrev(v int32) { this.prev = v }

// A MapleFootholdSorter is an implementation of sort.Interface
// that sorts MapleFoothold objects by Y
type MapleFootholdSorter []*MapleFoothold

func (this MapleFootholdSorter) Len() int      { return len(this) }
func (this MapleFootholdSorter) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this MapleFootholdSorter) Less(i, j int) bool {
	if this[i].Y2() < this[j].Y1() {
		return true
	} else if this[i].Y1() > this[j].Y2() {
		return false
	}
	return false // [i] == [j]
}
