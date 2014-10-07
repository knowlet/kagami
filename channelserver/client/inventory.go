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

package client

import "errors"
import "github.com/Francesco149/kagami/channelserver/gamedata"

// InventoryType defines which inventory tab this inventory refers to
type InventoryType int8

// Possible values for InventoryType
const (
	INVENTORY_EQUIPPED = -1
	INVENTORY_INVALID  = 0
	INVENTORY_EQUIP    = iota
	INVENTORY_USE
	INVENTORY_SETUP
	INVENTORY_ETC
	INVENTORY_CASH
)

// Bitmask returns the inventory type encoded as a bitmask
func (this InventoryType) Bitmask() uint16 {
	return uint16(2) << (uint16(this) % 32) // shifting by -1 (EQUIPPED) will result in << 31
}

// InventoryTypeByWzName returns the proper inventory type for the given name
func InventoryTypeByWzName(name string) InventoryType {
	switch name {
	case "Install":
		return INVENTORY_SETUP
	case "Consume":
		return INVENTORY_USE
	case "Etc":
		return INVENTORY_ETC
	case "Cash":
		return INVENTORY_CASH
	case "Pet":
		return INVENTORY_CASH
	default:
		return INVENTORY_INVALID
	}
	return INVENTORY_INVALID
}

// Inventory holds all of the items for a single inventory tab.
type Inventory struct {
	inv      map[int8]gamedata.GenericItem
	capacity int8
	typ      InventoryType
}

// Look for the given item in the inventory. If the item is not found, nil will be returned.
func (this *Inventory) ById(itemId int32) gamedata.GenericItem {
	for _, item := range this.inv {
		if item.Id() == itemId {
			return item
		}
	}
	return nil
}

// Add tries to find a free slot and puts the item in it. Returns the used slot.
// If no slots are available, Add will return -1.
func (this *Inventory) Add(i gamedata.GenericItem) int8 {
	slot := this.NextFreeSlot()
	if slot < 0 {
		return -1
	}
	this.inv[slot] = i
	i.SetPos(slot)
	return slot
}

// AddWithPosition adds an item that has a pre-defined position (such as items retrieved from the db)
func (this *Inventory) AddWithPosition(i gamedata.GenericItem) (err error) {
	if i.Pos() < 0 && this.typ != INVENTORY_EQUIPPED {
		err = errors.New("Tried to insert negative-position item in non-equip inventory.")
		return
	}
	this.inv[i.Pos()] = i
	return
}

// Move tries to move an item from slot A to slot B, stacking items if possible
func (this *Inventory) Move(slotA, slotB int8, maxStack int16) (err error) {
	itemA := this.inv[slotA]
	itemB := this.inv[slotB]

	switch {
	// empty source
	case itemA == nil:
		err = errors.New("Tried to move empty item slot")

	// move item to another empty slot
	case itemB == nil:
		itemA.SetPos(slotB)
		this.inv[slotB] = itemA
		delete(this.inv, slotA)

	// swapping identical items - see if they're stackable
	case itemB.Id() == itemA.Id() && gamedata.IsStackable(itemA):
		switch {
		// equips can't stack
		case this.Type() == INVENTORY_EQUIP:
			this.swap(itemA, itemB)

		// the stack overflows, so just stack as much as possible
		case itemA.Amount()+itemB.Amount() > maxStack:
			remainder := itemA.Amount() + itemB.Amount() - maxStack
			itemA.SetAmount(remainder)
			itemB.SetAmount(maxStack)

		// merge 2 stacks
		default:
			itemB.SetAmount(itemA.Amount() + itemB.Amount())
		}

	// swap two different items
	default:
		this.swap(itemA, itemB)
	}

	return
}

func (this *Inventory) swap(a, b gamedata.GenericItem) {
	this.inv[a.Pos()], this.inv[b.Pos()] = this.inv[b.Pos()], this.inv[a.Pos()]
	tmp := a.Pos()
	a.SetPos(b.Pos())
	b.SetPos(tmp)
}

func (this *Inventory) Get(slot int8) gamedata.GenericItem {
	return this.inv[slot]
}

func (this *Inventory) Remove(slot int8, amount int16) {
	item := this.inv[slot]

	if item == nil {
		return
	}

	item.SetAmount(item.Amount() - amount)

	switch {
	case item.Amount() < 0:
		item.SetAmount(0)
	case item.Amount() == 0:
		delete(this.inv, slot)
	}
}

func (this *Inventory) Full() bool {
	return len(this.inv) >= int(this.capacity)
}

func (this *Inventory) WillBeFull(addedAmount int) bool {
	return len(this.inv)+addedAmount >= int(this.capacity)
}

func (this *Inventory) NextFreeSlot() int8 {
	if this.Full() {
		return -1
	}

	for i := int8(1); i <= this.capacity; i++ {
		if this.inv[i] == nil {
			return i
		}
	}

	return -1
}

func (this *Inventory) Type() InventoryType {
	return this.typ
}

func (this *Inventory) Map() map[int8]gamedata.GenericItem { return this.inv }
