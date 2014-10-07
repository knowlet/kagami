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

// ItemType represents the type of an item stored in the player's inventory.
type ItemType byte

// Possible values for ItemType.
const (
	ITEM_ITEM  = 2
	ITEM_EQUIP = 1
)

// GenericItem is a generic interface for items stored in the inventory.
type GenericItem interface {
	Type() int8
	Pos() int8
	SetPos(v int8)
	Id() int32
	Amount() int16
	SetAmount(v int16)
	Owner() string
	SetOwner(v string)
	PetId() int32
	Clone() GenericItem
}
