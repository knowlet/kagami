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

// An ElementalEffectiveness object represents how effective a certain element is.
type ElementalEffectiveness int

// Possible values for ElementalEffectiveness
const (
	NORMAL = iota
	IMMUNE
	STRONG
	WEAK
	INVALID_EFFECTIVENESS
)

// ElementalEffectivenessFromNumber converts the wz representation of elemental
// effectiveness to an ElementalEffectiveness object.
func ElementalEffectivenessFromNumber(num int32) ElementalEffectiveness {
	switch num {
	case 1:
		return IMMUNE
	case 2:
		return STRONG
	case 3:
		return WEAK
	default:
		return INVALID_EFFECTIVENESS
	}

	return INVALID_EFFECTIVENESS
}
