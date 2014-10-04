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
	"github.com/Francesco149/maplelib"
)

// TODO: generic interface for movable entities

// A client.Connection is a MapleStory in-game client connected to the channel server.
// It's a wrapper around EncryptedConnection specialized for in-game MapleStory clients.
// It caches various data from the database such as gm level, look and so on.
type Connection struct {
	*common.EncryptedConnection        // underlying encrypted connection
	connected                   bool   // true if the player has successfully connected
	name                        string // character name
	admin                       bool   // true if the user is an admin
	gmchat                      bool   // true if the user's gm chat is enabled
	disconnecting               bool   // true if the user is disconnecting
	worldid                     int8   // numeric world id
	mapPos                      int8   // portal id at which the player will spawn
	gender                      byte   // the character's gender
	charid                      int32  // character id
	userid                      int32  // account id
	face                        int32  // face id
	hair                        int32  // hair id
	skin                        int8   // skin id
	mapid                       int32  // map id
	lastmap                     int32  // last map id
	gmLevel                     int32  // gm level
	uptime                      int64  // total online time in seconds
	buddylistSize               byte
	curmap                      *gamedata.MapleMap
}

// NewConnection initializes and returns an encrypted connection to a MapleStory client
func NewConnection(basecon net.Conn, testserver bool) *Connection {
	return &Connection{
		EncryptedConnection: common.NewEncryptedConnection(basecon, testserver, false), // base class
		admin:               false,
		gmchat:              false,
		disconnecting:       false,
		worldid:             -1,
		mapPos:              -1,
		gender:              0,
		charid:              -1,
		userid:              -1,
		face:                -1,
		hair:                -1,
		mapid:               -1,
		lastmap:             -1,
		gmLevel:             0,
		uptime:              0,
	}
}

func (c *Connection) String() string {
	return fmt.Sprintf(""+
		"\n%v:{\n"+
		"\tadmin: %v\n"+
		"\tgmchat: %v\n"+
		"\tdisconnecting: %v\n"+
		"\tworldid: %v\n"+
		"\tmapPos: %v\n"+
		"\tgender: %v\n"+
		"\tcharid: %v\n"+
		"\tuserid: %v\n"+
		"\tface: %v\n"+
		"\thair: %v\n"+
		"\tmapid: %v\n"+
		"\tlastmap: %v\n"+
		"\tgmLevel: %v\n"+
		"\tuptime: %v\n"+
		"}\n",
		c.Conn().RemoteAddr(), c.Admin(), c.GmChat(),
		c.Disconnecting(), c.WorldId(), c.MapPos(), c.Gender(),
		c.CharId(), c.UserId(), c.Face(), c.Hair(), c.MapId(),
		c.LastMap(), c.GmLevel(), c.Uptime())
}

func (c *Connection) Connected() bool                     { return c.connected }
func (c *Connection) SetConnected(connected bool)         { c.connected = connected }
func (c *Connection) Name() string                        { return c.name }
func (c *Connection) SetName(name string)                 { c.name = name }
func (c *Connection) Admin() bool                         { return c.admin }
func (c *Connection) SetAdmin(admin bool)                 { c.admin = admin }
func (c *Connection) GmChat() bool                        { return c.gmchat }
func (c *Connection) SetGmChat(gmchat bool)               { c.gmchat = gmchat }
func (c *Connection) Disconnecting() bool                 { return c.disconnecting }
func (c *Connection) SetDisconnecting(disconnecting bool) { c.disconnecting = disconnecting }
func (c *Connection) WorldId() int8                       { return c.worldid }
func (c *Connection) SetWorldId(worldid int8)             { c.worldid = worldid }
func (c *Connection) MapPos() int8                        { return c.mapPos }
func (c *Connection) SetMapPos(mapPos int8)               { c.mapPos = mapPos }
func (c *Connection) Gender() byte                        { return c.gender }
func (c *Connection) SetGender(gender byte)               { c.gender = gender }
func (c *Connection) CharId() int32                       { return c.charid }
func (c *Connection) SetCharId(charid int32)              { c.charid = charid }
func (c *Connection) UserId() int32                       { return c.userid }
func (c *Connection) SetUserId(userid int32)              { c.userid = userid }
func (c *Connection) Face() int32                         { return c.face }
func (c *Connection) SetFace(face int32)                  { c.face = face }
func (c *Connection) Hair() int32                         { return c.hair }
func (c *Connection) SetHair(hair int32)                  { c.hair = hair }
func (c *Connection) Skin() int8                          { return c.skin }
func (c *Connection) SetSkin(skin int8)                   { c.skin = skin }
func (c *Connection) MapId() int32                        { return c.mapid }

func (c *Connection) SetMapId(mapid int32) error {
	c.mapid = mapid
	fmt.Println("loading map", c.mapid)
	c.curmap = status.MapFactory().Get(mapid, true, true, true)
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
		oldmap := this.MapId()
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

	status.Lock()
	defer status.Unlock()
	return c.SendPacket(packets.WarpToMap(newmap.Id(), pid, 50, status.ChanId())) // todo: real hp

	// TODO: update party, player pool and everything
}

func (c *Connection) Map() *gamedata.MapleMap             { return c.curmap }
func (c *Connection) LastMap() int32                      { return c.lastmap }
func (c *Connection) SetLastMap(lastmap int32)            { c.lastmap = lastmap }
func (c *Connection) GmLevel() int32                      { return c.gmLevel }
func (c *Connection) SetGmLevel(gmLevel int32)            { c.gmLevel = gmLevel }
func (c *Connection) Uptime() int64                       { return c.uptime }
func (c *Connection) SetUptime(uptime int64)              { c.uptime = uptime }
func (c *Connection) BuddylistSize() byte                 { return c.buddylistSize }
func (c *Connection) SetBuddylistSize(buddylistSize byte) { c.buddylistSize = buddylistSize }

// SetDBOnline updates the player's online status in the database
func (c *Connection) SetDBOnline(online bool) (err error) {
	db := common.GetDB()
	st, err := db.Prepare("UPDATE `accounts` a INNER JOIN `characters` c ON a.id = c.user_id " +
		"SET a.online = ?, c.online = ? WHERE c.character_id = ?")
	_, err = st.Run(online, online, c.charid)
	return
}

func (c *Connection) EncodeStats(p *maplelib.Packet) (err error) {
	// TODO
	db := common.GetDB()
	st, err := db.Prepare("SELECT * FROM characters WHERE character_id = ?")
	res, err := st.Run(c.CharId())
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	if len(rows) < 1 {
		err = errors.New(fmt.Sprintf("Char id %d not found in database when encoding stats", c.CharId()))
		return
	}

	thechar, err := common.GetCharDataFromDBRow(rows[0], res)
	if err != nil {
		return
	}

	thechar.EncodeStats(p)
	return
}

func (c *Connection) EncodeQuestInfo(p *maplelib.Packet) {
	// TODO
	p.Encode2(0) // 0 started quests
	p.Encode2(0) // 0 completed
}

// SaveStats saves all of the player's stats to the database
func (c *Connection) SaveStats() (err error) {
	// TODO
	return nil
}

// Saves saves all of the player's information to the database
func (c *Connection) Save() (err error) {
	fmt.Println("Saving", c.Name(), "'s data")
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
