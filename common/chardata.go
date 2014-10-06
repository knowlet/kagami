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

package common

import (
	"errors"
	"fmt"
)

import (
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
	db := GetDB()
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

// CharStats holds a character's stats and appearance (not including equips)
type CharStats struct {
	id     int32
	name   string
	level  byte
	job    int16
	str    int16
	dex    int16
	intt   int16
	luk    int16
	hp     int16
	maxhp  int16
	mp     int16
	maxmp  int16
	ap     int16
	sp     int16
	exp    int32
	fame   int16
	mapp   int32
	pos    int8
	gender int8
	skin   int8
	face   int32
	hair   int32
}

func (this *CharStats) Id() int32    { return this.id }
func (this *CharStats) Name() string { return this.name }
func (this *CharStats) Level() byte  { return this.level }
func (this *CharStats) Job() int16   { return this.job }
func (this *CharStats) Str() int16   { return this.str }
func (this *CharStats) Dex() int16   { return this.dex }
func (this *CharStats) Int() int16   { return this.intt }
func (this *CharStats) Luk() int16   { return this.luk }
func (this *CharStats) Hp() int16    { return this.hp }
func (this *CharStats) MaxHp() int16 { return this.maxhp }
func (this *CharStats) Mp() int16    { return this.mp }
func (this *CharStats) MaxMp() int16 { return this.maxmp }
func (this *CharStats) Ap() int16    { return this.ap }
func (this *CharStats) Sp() int16    { return this.sp }
func (this *CharStats) Exp() int32   { return this.exp }
func (this *CharStats) Fame() int16  { return this.fame }
func (this *CharStats) MapId() int32 { return this.mapp }
func (this *CharStats) Pos() int8    { return this.pos }
func (this *CharStats) Gender() int8 { return this.gender }
func (this *CharStats) Skin() int8   { return this.skin }
func (this *CharStats) Face() int32  { return this.face }
func (this *CharStats) Hair() int32  { return this.hair }

func (this *CharStats) SetMapId(v int32) { this.mapp = v }

func (this *CharStats) String() string {
	return fmt.Sprintf(
		`{
	id: %v
	name: %v
	level: %v
	job: %v
	str: %v
	dex: %v
	int: %v
	luk: %v
	hp: %v
	maxhp: %v
	mp: %v
	maxmp: %v
	ap: %v
	sp: %v
	exp: %v
	fame: %v
	mapid: %v
	pos: %v
	gender: %v
	skin: %v
	face: %v
	hair: %v
}`,
		this.Id(),
		this.Name(),
		this.Level(),
		this.Job(),
		this.Str(),
		this.Dex(),
		this.Int(),
		this.Luk(),
		this.Hp(),
		this.MaxHp(),
		this.Mp(),
		this.MaxMp(),
		this.Ap(),
		this.Sp(),
		this.Exp(),
		this.Fame(),
		this.MapId(),
		this.Pos(),
		this.Gender(),
		this.Skin(),
		this.Face(),
		this.Hair(),
	)
}

// GetCharStatsFromDBRow retrieves the character stats from the given mysql query result
func GetCharStatsFromDBRow(row mysql.Row, res mysql.Result) *CharStats {
	colid := res.Map("character_id")
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

	return &CharStats{
		id:     int32(row.Int(colid)),
		name:   row.Str(colname),
		level:  byte(row.Int(collevel)),
		job:    int16(row.Int(coljob)),
		str:    int16(row.Int(colstr)),
		dex:    int16(row.Int(coldex)),
		intt:   int16(row.Int(colint)),
		luk:    int16(row.Int(colluk)),
		hp:     int16(row.Int(colchp)),
		maxhp:  int16(row.Int(colmhp)),
		mp:     int16(row.Int(colcmp)),
		maxmp:  int16(row.Int(colmmp)),
		ap:     int16(row.Int(colap)),
		sp:     int16(row.Int(colsp)),
		exp:    int32(row.Int(colexp)),
		fame:   int16(row.Int(colfame)),
		mapp:   int32(row.Int(colmap)),
		pos:    int8(row.Int(colpos)),
		gender: int8(row.Int(colgender)),
		skin:   int8(row.Int(colskin)),
		face:   int32(row.Int(colface)),
		hair:   int32(row.Int(colhair)),
	}
}

// Encode serializes the character's stats to the given packet
func (c *CharStats) Encode(p *maplelib.Packet) {
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

	p.Encode1s(c.gender)
	p.Encode1s(c.skin)
	p.Encode4s(c.face)
	p.Encode4s(c.hair)
	p.Append(huehuehue)
	p.Encode1(c.level)
	p.Encode2s(c.job)
	p.Encode2s(c.str)
	p.Encode2s(c.dex)
	p.Encode2s(c.intt)
	p.Encode2s(c.luk)
	p.Encode2s(c.hp)
	p.Encode2s(c.maxhp)
	p.Encode2s(c.mp)
	p.Encode2s(c.maxmp)
	p.Encode2s(c.ap)
	p.Encode2s(c.sp)
	p.Encode4s(c.exp)
	p.Encode2s(c.fame)
	p.Encode4(0x00000000) // married flag TODO
	p.Encode4s(c.mapp)
	p.Encode1s(c.pos) // initial spawnpoint
	p.Encode4(0x00000000)
	return
}

// CharData is a struct that holds character data retrieved from the database
// that will be used by handleViewAllChar and a couple other packet handlers.
type CharData struct {
	*CharStats
	worldRank     uint32
	worldRankMove uint32
	jobRank       uint32
	jobRankMove   uint32
	equips        []*CharEquipData
}

func (c *CharData) Stats() *CharStats { return c.CharStats }

// EncodeEquips serializes the character's equips to a packet
func (c *CharData) EncodeEquips(p *maplelib.Packet) {
	p.Encode1s(c.Stats().Gender()) // yes it repeats gender, skin, face, hair and idk why
	p.Encode1s(c.Stats().Skin())
	p.Encode4s(c.Stats().Face())
	p.Encode1(0x00)
	p.Encode4s(c.Stats().Hair())

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
				p.Encode4s(equipmap[i][1])
			} else {
				p.Encode4s(equipmap[i][0])
			}
		}
	}

	p.Encode1s(-1)
	// append covered items
	for i := byte(0); i < consts.EquippedSlots; i++ {
		if equipmap[i][1] > 0 && i != consts.EquipWeapon {
			p.Encode1(i)
			p.Encode4s(equipmap[i][1])
		}
	}

	p.Encode1s(-1)
	p.Encode4s(equipmap[consts.EquipWeapon][0]) // cash weapon

	ayylmao := make([]byte, 12)
	p.Append(ayylmao)
	return
}

// Encode encodes a charData object into a maplestory packet
func (c *CharData) Encode(p *maplelib.Packet) {
	c.CharStats.Encode(p)
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
	colworldcpos := res.Map("world_cpos")
	colworldopos := res.Map("world_opos")
	coljobcpos := res.Map("job_cpos")
	coljobopos := res.Map("job_opos")

	// TODO: ignore ranks for gm job

	charworldrank := uint32(row.Int(colworldcpos))
	charjobrank := uint32(row.Int(coljobcpos))

	cstats := GetCharStatsFromDBRow(row, res)
	cequips, err := GetCharEquipsFromDB(cstats.Id())
	if err != nil {
		return
	}

	data = &CharData{
		CharStats:     cstats,
		worldRank:     charworldrank,
		worldRankMove: charworldrank - uint32(row.Int(colworldopos)),
		jobRank:       charjobrank,
		jobRankMove:   charjobrank - uint32(row.Int(coljobopos)),
		equips:        cequips,
	}

	return
}
