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

func IsThrowingStar(i GenericItem) bool {
	return i.Id() >= 2070000 && i.Id() < 2080000
}

func IsBullet(i GenericItem) bool {
	return i.Id()/10000 == 233
}

func IsStackable(i GenericItem) bool {
	return !IsThrowingStar(i) && !IsBullet(i)
}
