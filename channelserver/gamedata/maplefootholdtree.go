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
	"math"
	"sort"
)

// A MapleFootholdTree holds and manages all of the footholds in a map.
type MapleFootholdTree struct {
	nw, ne, sw, se *MapleFootholdTree
	footholds      []*MapleFoothold
	p1             image.Point
	p2             image.Point
	center         image.Point
	depth          int
	maxDepth       int
	maxDropX       int
	minDropX       int
}

// NewMapleFootholdTree initializes a new, empty foothold tree with a
// default maxDepth of 8.
func NewMapleFootholdTree(tp1, tp2 image.Point) *MapleFootholdTree {
	return &MapleFootholdTree{
		p1:       tp1,
		p2:       tp2,
		center:   image.Pt((tp2.X-tp1.X)/2, (tp2.Y-tp1.Y)/2),
		maxDepth: 8,
	}
}

// NewMapleFootholdTreeWithDepth initializes a new, empty foothold tree with the
// given initial depth and a default maxDepth of 8.
func NewMapleFootholdTreeWithDepth(tp1, tp2 image.Point, depth int) *MapleFootholdTree {
	res := NewMapleFootholdTree(tp1, tp2)
	res.depth = depth
	return res
}

// Insert adds a foothold to the tree.
func (this *MapleFootholdTree) Insert(fh *MapleFoothold) {
	if this.depth == 0 {
		if fh.X1() > this.maxDropX {
			this.maxDropX = fh.X1()
		}

		if fh.X1() < this.minDropX {
			this.minDropX = fh.X1()
		}

		if fh.X2() > this.maxDropX {
			this.maxDropX = fh.X2()
		}

		if fh.X2() < this.minDropX {
			this.minDropX = fh.X2()
		}
	}

	if this.depth == this.maxDepth || (fh.X1() >= this.p1.X &&
		fh.X2() <= this.p2.X && fh.Y1() >= this.p1.Y && fh.Y2() >= this.p2.Y) {

		this.footholds = append(this.footholds, fh)
	} else {
		if this.nw == nil {
			this.nw = NewMapleFootholdTreeWithDepth(this.p1, this.center, this.depth+1)
			this.ne = NewMapleFootholdTreeWithDepth(image.Pt(this.center.X, this.p1.Y),
				image.Pt(this.p2.X, this.center.Y), this.depth+1)
			this.sw = NewMapleFootholdTreeWithDepth(image.Pt(this.p1.X, this.center.Y),
				image.Pt(this.center.X, this.p2.Y), this.depth+1)
			this.se = NewMapleFootholdTreeWithDepth(this.center, this.p2, this.depth+1)
		}

		switch {
		case fh.X2() <= this.center.X && fh.Y2() <= this.center.Y:
			this.nw.Insert(fh)
		case fh.X1() > this.center.X && fh.Y2() <= this.center.X:
			this.nw.Insert(fh)
		case fh.X2() <= this.center.X && fh.Y1() > this.center.Y:
			this.sw.Insert(fh)
		default:
			this.se.Insert(fh)
		}
	}
}

func (this *MapleFootholdTree) relevants(p image.Point,
	list []*MapleFoothold) []*MapleFoothold {

	if list == nil {
		return this.relevants(p, make([]*MapleFoothold, 0))
	}

	list = append(list, this.footholds...)
	if this.nw != nil {
		switch {
		case p.X <= this.center.X && p.Y <= this.center.Y:
			this.nw.relevants(p, list)
		case p.X > this.center.X && p.Y <= this.center.Y:
			this.ne.relevants(p, list)
		case p.X <= this.center.X && p.Y > this.center.Y:
			this.sw.relevants(p, list)
		default:
			this.se.relevants(p, list)
		}
	}

	return list
}

func (this *MapleFootholdTree) findWallR(p1, p2 image.Point) *MapleFoothold {
	var res *MapleFoothold

	for _, f := range this.footholds {
		if f.Wall() && f.X1() >= p1.X && f.X1() <= p2.X &&
			f.Y1() >= p1.X && f.Y2() <= p1.Y {

			return f
		}
	}

	if this.nw == nil {
		return nil
	}

	if p1.X <= this.center.X && p1.Y <= this.center.Y {
		res = this.nw.findWallR(p1, p2)
		if res != nil {
			return res
		}
	}

	if (p1.X > this.center.X || p2.X > this.center.X) && p1.Y <= this.center.Y {
		res = this.ne.findWallR(p1, p2)
		if res != nil {
			return res
		}
	}

	if p1.X <= this.center.X && p1.Y > this.center.Y {
		res = this.sw.findWallR(p1, p2)
		if res != nil {
			return res
		}
	}

	if (p1.X > this.center.X || p2.X > this.center.X) && p1.Y > this.center.Y {
		res = this.se.findWallR(p1, p2)
		if res != nil {
			return res
		}
	}

	return nil
}

func (this *MapleFootholdTree) FindWall(p1, p2 image.Point) *MapleFoothold {
	if p1.Y != p2.Y {
		return nil
	}

	return this.findWallR(p1, p2)
}

func (this *MapleFootholdTree) FindBelow(p image.Point) *MapleFoothold {
	relevants := this.relevants(p, nil)
	matches := make([]*MapleFoothold, 0)
	for _, fh := range relevants {
		if fh.X1() <= p.X && fh.X2() >= p.X {
			matches = append(matches, fh)
		}
	}
	sort.Sort(MapleFootholdSorter(relevants))
	for _, fh := range matches {
		if !fh.Wall() && fh.Y1() != fh.Y2() {
			var calcY int
			s1 := math.Abs(float64(fh.Y2() - fh.Y1()))
			s2 := math.Abs(float64(fh.X2() - fh.X1()))
			s4 := math.Abs(float64(p.X - fh.X1()))
			alpha := math.Atan(s2 / s1)
			beta := math.Atan(s1 / s2)
			s5 := math.Cos(alpha) * (s4 / math.Cos(beta))
			if fh.Y2() < fh.Y1() {
				calcY = fh.Y1() - int(s5)
			} else {
				calcY = fh.Y1() + int(s5)
			}

			if calcY >= p.Y {
				return fh
			}
		} else if !fh.Wall() && fh.Y1() >= p.Y {
			return fh
		}
	}

	return nil
}

func (this *MapleFootholdTree) X1() int       { return this.p1.X }
func (this *MapleFootholdTree) X2() int       { return this.p2.X }
func (this *MapleFootholdTree) Y1() int       { return this.p1.Y }
func (this *MapleFootholdTree) Y2() int       { return this.p2.Y }
func (this *MapleFootholdTree) MaxDropX() int { return this.maxDropX }
func (this *MapleFootholdTree) MinDropX() int { return this.minDropX }
