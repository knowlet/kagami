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

package main

import (
	"errors"
	"fmt"
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/maplelib"
	"github.com/ziutek/mymysql/mysql"
)

// CharEquipData is a struct that holds equip data retrieved from the database
// for a certain character that will be used by handleViewAllChar
type CharEquipData struct {
	id   int32
	slot int16
}

// GetCharEquipsFromDB retrieves all of the given character's equips from
// the database and returns them as an array
func GetCharEquipsFromDB(characterId int32) (res []*CharEquipData, err error) {
	db := common.GetDB()
	st, err := db.Prepare("SELECT item_id, slot FROM items " +
		"WHERE character_id = ? AND inv = ? AND slot < 0 " +
		"ORDER BY slot ASC")
	stres, err := st.Run(characterId, consts.EquipInventory)
	rows, err := stres.GetRows()
	if err != nil {
		return
	}

	colitemid := stres.Map("item_id")
	colslot := stres.Map("slot")

	// pre-allocate slice so that it doesn't reallocate it when appending
	res = make([]*CharEquipData, len(rows))

	for i, row := range rows {
		res[i] = &CharEquipData{
			id:   int32(row.Int(colitemid)),
			slot: int16(row.Int(colslot)),
		}
	}

	return
}

// CharData is a struct that holds character data retrieved from the database
// that will be used by handleViewAllChar
type CharData struct {
	id            int32
	name          string
	level         byte
	job           int16
	str           int16
	dex           int16
	intt          int16
	luk           int16
	hp            int16
	maxhp         int16
	mp            int16
	maxmp         int16
	ap            int16
	sp            int16
	exp           int32
	fame          int16
	mapp          int32
	pos           int8
	gender        byte
	skin          byte
	face          int32
	hair          int32
	worldRank     uint32
	worldRankMove uint32
	jobRank       uint32
	jobRankMove   uint32
	equips        []*CharEquipData
}

// EncodeStats serializes the character's stats to the given packet
func (c *CharData) EncodeStats(p maplelib.Packet) {
	huehuehue := make([]byte, 24)
	namelen := len(c.name)

	if namelen > consts.MaxNameSize {
		panic(errors.New(
			fmt.Sprintf("Tried to encode char name %s "+
				"which is bigger than %d characters",
				c.name, consts.MaxNameSize)))
	}

	p.Encode4(uint32(c.id))

	// name is encoded as a constant-size null terminated string
	p.Append([]byte(c.name))

	// fill the remaining space with null termination characters
	for i := namelen; i < consts.MaxNameSize+1; i++ {
		p.Encode1(0x00)
	}

	p.Encode1(c.gender)
	p.Encode1(c.skin)
	p.Encode4(uint32(c.face))
	p.Encode4(uint32(c.hair))
	p.Append(huehuehue)
	p.Encode1(c.level)
	p.Encode2(uint16(c.job))
	p.Encode2(uint16(c.str))
	p.Encode2(uint16(c.dex))
	p.Encode2(uint16(c.intt))
	p.Encode2(uint16(c.luk))
	p.Encode2(uint16(c.hp))
	p.Encode2(uint16(c.maxhp))
	p.Encode2(uint16(c.mp))
	p.Encode2(uint16(c.maxmp))
	p.Encode2(uint16(c.ap))
	p.Encode2(uint16(c.sp))
	p.Encode4(uint32(c.exp))
	p.Encode2(uint16(c.fame))
	p.Encode4(0x00000000) // married flag TODO
	p.Encode4(uint32(c.mapp))
	p.Encode1(byte(c.pos)) // initial spawnpoint
	p.Encode4(0x00000000)
	return
}

// EncodeEquips serializes the character's equips to a packet
func (c *CharData) EncodeEquips(p maplelib.Packet) {
	p.Encode1(c.gender) // yes it repeats gender, skin, face, hair and idk why
	p.Encode1(c.skin)
	p.Encode4(uint32(c.face))
	p.Encode1(0x01)
	p.Encode4(uint32(c.hair))

	// I'm not sure how this all works but it's some logic to encode
	// equips in such a way that the client can determine which ones are
	// covered by cash shop items or something I DUNNO FUCK
	// I'm guessing equipmap[i][0] refers to non-cash items and equipmap[i][1]
	// contains cash items or items that cover other equips
	// If someone understands this, feel free to comment and explain it please

	var equipmap [consts.EquippedSlots][2]int32

	for _, equip := range c.equips {
		slot := -equip.slot

		if slot > 100 {
			slot -= 100
		}

		if equipmap[slot][0] > 0 {
			if equip.slot < -100 {
				equipmap[slot][1] = equipmap[slot][0]
				equipmap[slot][0] = equip.id
			} else {
				equipmap[slot][1] = equip.id
			}
		} else {
			equipmap[slot][0] = equip.id
		}
	}

	// append shown equips
	for i := byte(0); i < consts.EquippedSlots; i++ {
		if equipmap[i][0] > 0 {
			p.Encode1(i)

			if i == consts.EquipWeapon && equipmap[i][1] > 0 {
				// normal weapons
				p.Encode4(uint32(equipmap[i][1]))
			} else {
				p.Encode4(uint32(equipmap[i][0]))
			}
		}
	}

	p.Encode1(0xFF) // -1 as uint8

	// append covered items
	for i := byte(0); i < consts.EquippedSlots; i++ {
		if equipmap[i][1] > 0 && i != consts.EquipWeapon {
			p.Encode1(i)
			p.Encode4(uint32(equipmap[i][1]))
		}
	}

	p.Encode1(0xFF)
	p.Encode4(uint32(equipmap[consts.EquipWeapon][0])) // cash weapon

	ayylmao := make([]byte, 12)
	p.Append(ayylmao)
	return
}

// Encode encodes a charData object into a maplestory packet
func (c *CharData) Encode(p maplelib.Packet) {
	c.EncodeStats(p)
	c.EncodeEquips(p)

	// ranks
	p.Encode1(0x01)
	p.Encode4(c.worldRank)
	p.Encode4(c.worldRankMove)
	p.Encode4(c.jobRank)
	p.Encode4(c.jobRankMove)
	return
}

// GetCharDataFromDBRow populates a charData structure with the character data in the
// given mysql row, which must belong to the given mysql result
func GetCharDataFromDBRow(row mysql.Row, res mysql.Result) (data *CharData, err error) {
	// column indices
	colid := res.Map("id")
	colname := res.Map("name")
	colgender := res.Map("gender")
	colskin := res.Map("skin")
	colface := res.Map("face")
	colhair := res.Map("hair")
	collevel := res.Map("level")
	coljob := res.Map("job")
	colstr := res.Map("str")
	coldex := res.Map("dex")
	colint := res.Map("int")
	colluk := res.Map("luk")
	colchp := res.Map("chp")
	colmhp := res.Map("mhp")
	colcmp := res.Map("cmp")
	colmmp := res.Map("mmp")
	colap := res.Map("ap")
	colsp := res.Map("sp")
	colexp := res.Map("exp")
	colfame := res.Map("fame")
	colmap := res.Map("map")
	colpos := res.Map("pos")
	colworldcpos := res.Map("world_cpos")
	colworldopos := res.Map("world_opos")
	coljobcpos := res.Map("job_cpos")
	coljobopos := res.Map("job_opos")

	// TODO: ignore ranks for gm job

	// reusable stuff
	charid := int32(row.Int(colid))
	charworldrank := uint32(row.Int(colworldcpos))
	charjobrank := uint32(row.Int(coljobcpos))

	charequips, err := GetCharEquipsFromDB(charid)
	if err != nil {
		return
	}

	data = &CharData{
		id:            charid,
		name:          row.Str(colname),
		level:         byte(row.Int(collevel)),
		job:           int16(row.Int(coljob)),
		str:           int16(row.Int(colstr)),
		dex:           int16(row.Int(coldex)),
		intt:          int16(row.Int(colint)),
		luk:           int16(row.Int(colluk)),
		hp:            int16(row.Int(colchp)),
		maxhp:         int16(row.Int(colmhp)),
		mp:            int16(row.Int(colcmp)),
		maxmp:         int16(row.Int(colmmp)),
		ap:            int16(row.Int(colap)),
		sp:            int16(row.Int(colsp)),
		exp:           int32(row.Int(colexp)),
		fame:          int16(row.Int(colfame)),
		mapp:          int32(row.Int(colmap)),
		pos:           int8(row.Int(colpos)),
		gender:        byte(row.Int(colgender)),
		skin:          byte(row.Int(colskin)),
		face:          int32(row.Int(colface)),
		hair:          int32(row.Int(colhair)),
		worldRank:     charworldrank,
		worldRankMove: charworldrank - uint32(row.Int(colworldopos)),
		jobRank:       charjobrank,
		jobRankMove:   charjobrank - uint32(row.Int(coljobopos)),
		equips:        charequips,
	}

	return
}
