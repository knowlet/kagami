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
	"errors"
	"time"
)

import (
	"github.com/knowlet/kagami/common/consts"
	"github.com/knowlet/kagami/common/utils"
	"github.com/Francesco149/maplelib"
)

var ITEM_MAGIC = []byte{0x80, 0x05}

func koreanTimestamp(millisec int64) int64 {
	time := millisec / 1000 / 60
	return time*600000000 + 116444592000000000
}

func encodeExpiration(p *maplelib.Packet, time int64, show bool) {
	if time != 0 {
		p.Encode4s(utils.MakeItemTimestamp(time))
	} else {
		p.Encode4s(400967355)
	}

	if show {
		p.Encode1(0x01)
	} else {
		p.Encode1(0x00)
	}
}

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
	Encode(p *maplelib.Packet)
}

func encodePetItemInfo(this GenericItem, p *maplelib.Packet) {
	// TODO: get pet from db
	petname := "faggot"
	petlevel := byte(1)
	petcloseness := int16(0)
	petfullness := int8(0)

	p.Encode1(0x01)
	p.Encode4s(this.PetId())
	p.Encode4(0x00000000)
	p.Encode1(0x00)
	p.Append(ITEM_MAGIC)
	p.Append([]byte{0xBB, 0x46, 0xE6, 0x17, 0x02})

	// name is encoded as a constant-size null terminated string
	p.Append([]byte(petname))

	// fill the remaining space with null termination characters
	for i := len(petname); i < consts.MaxNameSize+1; i++ {
		p.Encode1(0x00)
	}

	p.Encode1(petlevel)
	p.Encode2s(petcloseness)
	p.Encode1s(petfullness)
	p.Encode8s(koreanTimestamp(int64(float64(
		time.Now().UnixNano()/1000000) * 1.2)))
	p.Encode4(0x00000000)
}

// Item holds information about a normal item
type Item struct {
	id     int32
	slot   int8
	amount int16
	petid  int32
	owner  string
}

// NewItem initializes and returns a new item object
func NewItem(iid int32, islot int8, iamount int16, ipetid int32) *Item {
	return &Item{
		id:     iid,
		slot:   islot,
		amount: iamount,
		petid:  ipetid,
	}
}

func (this *Item) Type() int8        { return ITEM_ITEM }
func (this *Item) Pos() int8         { return this.slot }
func (this *Item) SetPos(v int8)     { this.slot = v }
func (this *Item) Id() int32         { return this.id }
func (this *Item) Amount() int16     { return this.amount }
func (this *Item) SetAmount(v int16) { this.amount = v }
func (this *Item) Owner() string     { return this.owner }
func (this *Item) SetOwner(v string) { this.owner = v }
func (this *Item) PetId() int32      { return this.petid }

// Clone returns a copy of this item
func (this *Item) Clone() GenericItem {
	res := NewItem(this.Id(), this.Pos(), this.Amount(), this.PetId())
	res.SetOwner(this.Owner())
	return res
}

func (this *Item) Encode(p *maplelib.Packet) {
	p.Encode1s(this.Pos())

	pet := this.PetId() > -1

	if pet {
		p.Encode1(0x03)
	} else {
		p.Encode1s(this.Type())
	}

	p.Encode4s(this.Id())

	if pet {
		encodePetItemInfo(this, p)
		return
	}

	p.Encode2(0x0000)
	p.Append(ITEM_MAGIC)
	encodeExpiration(p, 0, false)

	p.Encode2s(this.Amount())
	p.EncodeString(this.Owner())
	p.Encode2(0x0000)

	if !IsStackable(this) {
		p.Append([]byte{0x02, 0x00, 0x00, 0x00, 0x54, 0x00, 0x00, 0x34})
	}
}

// Equip holds the information about an equip item.
type Equip struct {
	*Item
	level         byte
	slots, locked int8
	job           int16
	str, dex, intt, luk, hp, mp, watk, matk, wdef, mdef, acc,
	avoid, hands, speed, jump int16
	ringid int32
}

// NewEquip initializes and returns an Equip object
func NewEquip(id int32, slot int8, ring int32) *Equip {
	return &Equip{
		Item:  NewItem(id, slot, 1, -1),
		slots: 0, level: 0, locked: 0,
		job: 0, str: 0, dex: 0, intt: 0, luk: 0, hp: 0,
		mp: 0, watk: 0, matk: 0, wdef: 0, mdef: 0, acc: 0,
		avoid: 0, hands: 0, speed: 0, jump: 0,
		ringid: ring,
	}
}

func (this *Equip) Encode(p *maplelib.Packet) {
	ring := this.RingId() > -1
	pos := this.Pos()
	masking := false
	equipped := pos < 0
	pet := this.PetId() > -1

	if equipped {
		pos *= -1

		if (pos > 100 || pos == -128) || ring {
			masking = true
			p.Encode1(0x00)
			p.Encode1s(pos - 100)
		} else {
			p.Encode1s(pos)
		}
	} else {
		p.Encode1s(this.Pos())
	}

	if this.PetId() > -1 {
		p.Encode1(0x03)
	} else {
		p.Encode1s(this.Type())
	}

	p.Encode4s(this.Id())

	if ring {
		p.Encode1(0x01)
		p.Encode4s(this.RingId())
		p.Encode4(0x00000000)
	}

	if pet {
		encodePetItemInfo(this, p)
		return
	}

	switch {
	case masking && !ring:
		p.Append([]byte{0x01, 0x41, 0xB4, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x20, 0x6F})
		encodeExpiration(p, 0, false)
	case ring:
		p.Encode8s(koreanTimestamp(int64(float64(
			time.Now().UnixNano()/1000000) * 1.2)))
	default:
		p.Encode2(0x0000)
		p.Append(ITEM_MAGIC)
		encodeExpiration(p, 0, false)
	}

	p.Encode1s(this.slots)
	p.Encode1(this.level)
	p.Encode2s(this.str)
	p.Encode2s(this.dex)
	p.Encode2s(this.intt)
	p.Encode2s(this.luk)
	p.Encode2s(this.hp)
	p.Encode2s(this.mp)
	p.Encode2s(this.watk)
	p.Encode2s(this.matk)
	p.Encode2s(this.wdef)
	p.Encode2s(this.mdef)
	p.Encode2s(this.acc)
	p.Encode2s(this.avoid)
	p.Encode2s(this.hands)
	p.Encode2s(this.speed)
	p.Encode2s(this.jump)
	p.EncodeString(this.Owner())
	p.Encode1s(this.locked)

	p.Encode1(0x00)

	if !masking && !ring {
		p.Encode8(0x0000000000000000)
	}
}

func (this *Equip) Type() int8        { return ITEM_EQUIP }
func (this *Equip) SetRingId(v int32) { this.ringid = v }
func (this *Equip) RingId() int32     { return this.ringid }
func (this *Equip) SetAmount(v int16) {
	if v != 1 {
		panic(errors.New("cannot set equip amount"))
	}
	this.amount = v
}

// Clone returns a copy of this equip
func (this *Equip) Clone() GenericItem {
	res := NewEquip(this.Item.Id(), this.Item.Pos(), this.RingId())
	res.slots = this.slots
	res.level = this.level
	res.locked = this.locked
	res.job = this.job
	res.str = this.str
	res.dex = this.dex
	res.intt = this.intt
	res.luk = this.luk
	res.hp = this.hp
	res.mp = this.mp
	res.watk = this.watk
	res.matk = this.matk
	res.wdef = this.wdef
	res.mdef = this.mdef
	res.acc = this.acc
	res.avoid = this.avoid
	res.hands = this.hands
	res.speed = this.speed
	res.jump = this.jump
	return res
}
