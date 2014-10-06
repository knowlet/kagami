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

// Package client contains various data structures and utilities related to individual
// maplestory clients that are currently connected to the channel server
package client

import (
	"errors"
	"fmt"
	"net"
)

import (
	"github.com/Francesco149/kagami/channelserver/gamedata"
	"github.com/Francesco149/kagami/channelserver/status"
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/kagami/common/utils"
	"github.com/Francesco149/maplelib"
)

// TODO: generic interface for movable entities

// A client.Connection is a MapleStory in-game client connected to the channel server.
// It's a wrapper around EncryptedConnection specialized for in-game MapleStory clients.
// It caches various data from the database such as gm level, look and so on.
type Connection struct {
	*common.EncryptedConnection       // underlying encrypted connection
	connected                   bool  // true if the player has successfully connected
	admin                       bool  // true if the user is an admin
	gmchat                      bool  // true if the user's gm chat is enabled
	disconnecting               bool  // true if the user is disconnecting
	worldid                     int8  // numeric world id
	userid                      int32 // account id
	lastmap                     int32 // last map id
	gmLevel                     int32 // gm level
	uptime                      int64 // total online time in seconds
	buddylistSize               byte
	curmap                      *gamedata.MapleMap
	meso                        int32
	stats                       *common.CharStats
}

// NewConnection initializes and returns an encrypted connection to a MapleStory client
func NewConnection(basecon net.Conn, testserver bool) *Connection {
	return &Connection{
		EncryptedConnection: common.NewEncryptedConnection(basecon, testserver, false), // base class
		connected:           false,
		admin:               false,
		gmchat:              false,
		disconnecting:       false,
		worldid:             -1,
		userid:              -1,
		lastmap:             -1,
		gmLevel:             0,
		uptime:              0,
		buddylistSize:       0,
		curmap:              nil,
		meso:                -1,
		stats:               nil,
	}
}

func (c *Connection) String() string {
	return fmt.Sprintf(
		`%v:{
	connected: %v
	admin: %v
	gmchat: %v
	disconnecting: %v
	worldid: %v
	userid: %v
	lastmap: %v
	gmLevel: %v
	uptime: %v
	buddylistSize: %v
	map: %v
	meso: %v
	stats: %v
}`,
		c.Conn().RemoteAddr(),
		c.Connected(),
		c.Admin(),
		c.GmChat(),
		c.Disconnecting(),
		c.WorldId(),
		c.UserId(),
		c.LastMap(),
		c.GmLevel(),
		c.Uptime(),
		c.BuddylistSize(),
		utils.Indent(c.Map().String(), 1),
		c.Meso(),
		utils.Indent(c.Stats().String(), 1),
	)
}

func (c *Connection) Connected() bool                     { return c.connected }
func (c *Connection) SetConnected(connected bool)         { c.connected = connected }
func (c *Connection) Admin() bool                         { return c.admin }
func (c *Connection) SetAdmin(admin bool)                 { c.admin = admin }
func (c *Connection) GmChat() bool                        { return c.gmchat }
func (c *Connection) SetGmChat(gmchat bool)               { c.gmchat = gmchat }
func (c *Connection) Disconnecting() bool                 { return c.disconnecting }
func (c *Connection) SetDisconnecting(disconnecting bool) { c.disconnecting = disconnecting }
func (c *Connection) WorldId() int8                       { return c.worldid }
func (c *Connection) SetWorldId(worldid int8)             { c.worldid = worldid }
func (c *Connection) UserId() int32                       { return c.userid }
func (c *Connection) SetUserId(userid int32)              { c.userid = userid }
func (c *Connection) Meso() int32                         { return c.meso }
func (c *Connection) SetMeso(v int32)                     { c.meso = v }
func (c *Connection) Stats() *common.CharStats            { return c.stats }
func (c *Connection) Map() *gamedata.MapleMap             { return c.curmap }
func (c *Connection) LastMap() int32                      { return c.lastmap }
func (c *Connection) SetLastMap(lastmap int32)            { c.lastmap = lastmap }
func (c *Connection) GmLevel() int32                      { return c.gmLevel }
func (c *Connection) SetGmLevel(gmLevel int32)            { c.gmLevel = gmLevel }
func (c *Connection) Uptime() int64                       { return c.uptime }
func (c *Connection) SetUptime(uptime int64)              { c.uptime = uptime }
func (c *Connection) BuddylistSize() byte                 { return c.buddylistSize }
func (c *Connection) SetBuddylistSize(buddylistSize byte) { c.buddylistSize = buddylistSize }
func (c *Connection) Alive() bool                         { return c.Stats().Hp() > 0 }

// LoadFromDB retrieves the given character id's data and assigns it to this connection
func (con *Connection) LoadFromDB(charid int32) (err error) {
	// get char data from db
	db := common.GetDB()
	st, err := db.Prepare("SELECT c.*, a.gm_level, a.admin FROM `characters` c " +
		"INNER JOIN `accounts` a ON c.user_id = a.id " +
		"WHERE c.character_id = ?")
	if err != nil {
		fmt.Println("Unexpected invalid query in handleLoadCharacter")
		return
	}
	res, err := st.Run(charid)
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	if len(rows) < 1 {
		err = errors.New("Character not found.")
		return
	}

	row := rows[0]

	cstats := common.GetCharStatsFromDBRow(row, res)

	coluserid := res.Map("user_id")
	colgmlevel := res.Map("gm_level")
	coladmin := res.Map("admin")
	colworldid := res.Map("world_id")
	colmeso := res.Map("meso")
	colbuddysize := res.Map("buddylist_size")

	/*
		colequipslots := res.Map("equip_slots")
		coluseslots := res.Map("use_slots")
		coletcslots := res.Map("etc_slots")
		colcashslots := res.Map("cash_slots")
	*/

	con.SetUserId(int32(row.Int(coluserid)))
	con.SetGmLevel(int32(row.Int(colgmlevel)))
	con.SetAdmin(row.Int(coladmin) > 0)
	con.SetWorldId(int8(row.Int(colworldid)))
	con.SetStats(cstats)
	con.SetMeso(int32(row.Int(colmeso)))
	con.SetBuddylistSize(byte(row.Int(colbuddysize)))

	// TODO: get max inventory slots and init inventories

	// TODO: do not reset uptime if the player is just xfering

	con.SetUptime(0)
	con.SetGmChat(con.GmChat() && con.GmLevel() > 0)

	// TODO: get book cover (wtf is a book cover)
	// TODO: init keymaps
	// TODO: init hpmp

	return
}

func (c *Connection) SetStats(v *common.CharStats) {
	c.stats = v
	c.SetMapId(v.MapId())
}

func (c *Connection) SetMapId(mapid int32) error {
	c.Stats().SetMapId(mapid)
	fmt.Println("loading map", c.Stats().MapId())

	st := <-status.Get
	defer func() { status.Get <- st }()
	c.curmap = st.MapFactory().Get(mapid, true, true, true)

	if c.curmap == nil {
		return errors.New("failed to load map")
	}
	fmt.Println("done")
	return nil
}

// Enter warps the client through this portal if possible
func (this *Connection) Enter(p gamedata.IMapleGenericPortal) (err error) {
	// TODO: check distance from portal and D/C if hacking

	changedMap := false
	if len(p.ScriptName()) > 0 {
		// TODO: handle 4th job portal script
		return
	}

	if p.TargetMapId() != 999999999 {
		oldmap := this.Stats().MapId()
		err = this.SetMapId(p.TargetMapId())
		if err != nil {
			this.SetMapId(oldmap)
			// TODO: send some error
			return this.SendPacket(packets.EnableActions())
		}

		newportal := this.Map().Portal(p.Target())
		if newportal == nil {
			newportal = this.Map().PortalById(0)
		}

		err = this.WarpToMap(this.Map(), newportal)
		changedMap = true
	}

	if !changedMap {
		err = this.SendPacket(packets.EnableActions())
	}

	return
}

// WarpToMap sends a map warp packet for the given map and portal.
// NOTE: this must be called after calling SetMapId
func (c *Connection) WarpToMap(newmap *gamedata.MapleMap,
	newportal gamedata.MaplePortal) error {

	pid := newportal.Id()

	switch newmap.Id() {
	case 100000200, 211000100, 220000300: // dunno why you have to change portal id here
		pid -= 2
	}

	st := <-status.Get
	defer func() { status.Get <- st }()
	return c.SendPacket(packets.WarpToMap(newmap.Id(), pid, 50, st.ChanId())) // todo: real hp

	// TODO: update party, player pool and everything
}

// SetDBOnline updates the player's online status in the database
func (c *Connection) SetDBOnline(online bool) (err error) {
	db := common.GetDB()
	st, err := db.Prepare("UPDATE `accounts` a INNER JOIN `characters` c ON a.id = c.user_id " +
		"SET a.online = ?, c.online = ? WHERE c.character_id = ?")
	_, err = st.Run(online, online, c.Stats().Id())
	return
}

func (c *Connection) EncodeQuestInfo(p *maplelib.Packet) {
	// TODO
	p.Encode2(0) // 0 started quests
	p.Encode2(0) // 0 completed
}

// SaveStats saves all of the player's stats to the database
func (c *Connection) SaveStats() (err error) {
	db := common.GetDB()
	st, err := db.Prepare(
		"UPDATE characters SET " +
			"level = ?, " +
			"job = ?, " +
			"str = ?, " +
			"dex = ?, " +
			"int = ?, " +
			"luk = ?, " +
			"chp = ?, " +
			"mhp = ?, " +
			"cmp = ?, " +
			"mmp = ?, " +
			"ap = ?, " +
			"sp = ?, " +
			"exp = ?, " +
			"fame = ?, " +
			"map = ?, " +
			"pos = ?, " +
			"gender = ?, " +
			"skin = ?, " +
			"face = ?, " +
			"hair = ? " +
			"WHERE character_id = ?")
	if err != nil {
		return
	}

	_, err = st.Run(
		c.Stats().Level(),
		c.Stats().Job(),
		c.Stats().Str(),
		c.Stats().Dex(),
		c.Stats().Int(),
		c.Stats().Luk(),
		c.Stats().Hp(),
		c.Stats().MaxHp(),
		c.Stats().Mp(),
		c.Stats().MaxMp(),
		c.Stats().Ap(),
		c.Stats().Sp(),
		c.Stats().Exp(),
		c.Stats().Fame(),
		c.Stats().MapId(),
		c.Stats().Pos(),
		c.Stats().Gender(),
		c.Stats().Skin(),
		c.Stats().Face(),
		c.Stats().Hair(),
		c.Stats().Id(),
	)
	return
}

// Saves saves all of the player's information to the database
func (c *Connection) Save() (err error) {
	fmt.Println("Saving", c.Stats().Name(), "'s data")
	err = c.SaveStats()
	// TODO: save inventory
	// TODO: save storage
	// TODO: save monster book
	// TODO: save mounts
	// TODO: save pets
	// TODO: save quests
	// TODO: save skills
	// TODO: save variables
	return
}
