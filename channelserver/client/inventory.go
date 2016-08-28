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

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

import (
	"github.com/knowlet/kagami/channelserver/gamedata"
	"github.com/knowlet/kagami/common"
	"github.com/Francesco149/maplelib"
)

// InventoryType defines which inventory tab this inventory refers to
type InventoryType int8

// Possible values for InventoryType
const (
	INVENTORY_EQUIPPED = -1
	INVENTORY_INVALID  = 0
	INVENTORY_EQUIP    = 1
	INVENTORY_USE      = 2
	INVENTORY_SETUP    = 3
	INVENTORY_ETC      = 4
	INVENTORY_CASH     = 5
)

// Bitmask returns the inventory type encoded as a bitmask
func (this InventoryType) Bitmask() uint16 {
	return uint16(2) << (uint16(this) % 32)
	// Shifting by -1 (EQUIPPED) will result in << 31.
	// This mimicks java's negative shifting behaviour.
	// Basically, if you shift by a negative value it will
	// take the first 5 bits of the shift value (mod 32).
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

// NewInventory initializes a new inventory data structure.
func NewInventory(invtype InventoryType, maxSlots int8) *Inventory {
	return &Inventory{
		inv:      make(map[int8]gamedata.GenericItem),
		capacity: maxSlots,
		typ:      invtype,
	}
}

type itemSorter []gamedata.GenericItem

func (this itemSorter) Len() int      { return len(this) }
func (this itemSorter) Swap(i, j int) { this[i], this[j] = this[j], this[i] }
func (this itemSorter) Less(i, j int) bool {
	return math.Abs(float64(this[i].Pos())) < math.Abs(float64(this[j].Pos()))
}

// Encode encodes an entire inventory to a maple packet
func (this *Inventory) Encode(p *maplelib.Packet) {
	if this.Type() == INVENTORY_EQUIPPED {
		ordered := make([]gamedata.GenericItem, 0)
		for _, item := range this.inv {
			if item == nil {
				continue
			}
			ordered = append(ordered, item)
		}

		if len(ordered) > 0 {
			sort.Sort(itemSorter(ordered))

			for _, item := range ordered {
				item.Encode(p)
			}
		}
	}

	for _, item := range this.inv {
		item.Encode(p)
	}
}

func (this *Inventory) Capacity() int8 { return this.capacity }

func (this *Inventory) LoadFromDB(charid int32) (err error) {
	q := "SELECT * FROM items WHERE location = 'inventory' AND character_id = ? AND inv = ?"

	t := this.Type()
	if t == INVENTORY_EQUIPPED {
		t = INVENTORY_EQUIP
		q += " AND slot < 0"
	} else {
		q += " AND slot > 0"
	}

	db := common.GetDB()
	st, err := db.Prepare(q)
	if err != nil {
		fmt.Println("Unexpected invalid query in inventory.LoadFromDB")
		return
	}

	res, err := st.Run(charid, int8(t))
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	colitemid := res.Map("item_id")
	colslot := res.Map("slot")
	colamount := res.Map("amount")

	if rows == nil || len(rows) == 0 {
		return
	}

	for _, row := range rows {
		if row == nil {
			return
		}

		var it gamedata.GenericItem

		if t == INVENTORY_EQUIP {
			it = gamedata.NewEquip(int32(row.Int(colitemid)), int8(row.Int(colslot)),
				-1) // todo: get ring id from db
			// TODO: set equip data n shit
		} else {
			it = gamedata.NewItem(int32(row.Int(colitemid)), int8(row.Int(colslot)),
				int16(row.Int(colamount)), -1) // todo: get pet id from db
		}

		if err = this.AddWithPosition(it); err != nil {
			return
		}
	}

	return
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

// Get returns the item located at the given slot
func (this *Inventory) Get(slot int8) gamedata.GenericItem {
	return this.inv[slot]
}

// Remove removes the given amount of the item located at the given slot
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

// Full returns true if all the slots in the inventory are currently being used
func (this *Inventory) Full() bool {
	return len(this.inv) >= int(this.capacity)
}

// WillBeFull returns true if all of the slots will be occupied after then given
// amount of additional slots is occupied.
func (this *Inventory) WillBeFull(addedAmount int) bool {
	return len(this.inv)+addedAmount >= int(this.capacity)
}

// NextFreeSlot tries to find the first free slot in the inventory.
// If no free slot is found, NextFreeSlot returns -1.
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

// Type returns the inventory's type.
// See InventoryType.
func (this *Inventory) Type() InventoryType {
	return this.typ
}

// Map returns a map by slot of the contents of the inventory.
func (this *Inventory) Map() map[int8]gamedata.GenericItem { return this.inv }
