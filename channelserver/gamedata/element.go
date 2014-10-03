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

import "strings"

// Element is a maplestory element type
type Element int

// Possible values for Element
const (
	NEUTRAL = iota
	FIRE
	ICE
	LIGHTING
	POISON
	HOLY /* SHIT */
	INVALID_ELEMENT
)

// ElementFromChar converts the wz representation of elements to an Element
func ElementFromChar(c string) Element {
	switch strings.ToLower(c) {
	case "f":
		return FIRE
	case "i":
		return ICE
	case "l":
		return LIGHTING
	case "s":
		return POISON
	case "h":
		return HOLY
	default:
		return INVALID_ELEMENT
	}

	return INVALID_ELEMENT
}
