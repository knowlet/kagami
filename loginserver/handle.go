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
	"time"
)

import (
	"github.com/knowlet/kagami/common"
	"github.com/knowlet/kagami/common/consts"
	"github.com/knowlet/kagami/common/interserver"
	"github.com/knowlet/kagami/common/packets"
	"github.com/knowlet/kagami/common/utils"
	"github.com/knowlet/kagami/loginserver/client"
	"github.com/knowlet/kagami/loginserver/items"
	"github.com/knowlet/kagami/loginserver/validators"
	"github.com/knowlet/kagami/loginserver/worlds"
	"github.com/knowlet/maplelib"
)

// Handle handles loginserver packets
func Handle(con *client.Connection, p maplelib.Packet) (handled bool, err error) {
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	switch header {
	case packets.IUnknownPlsIgnore1:
		return true, nil

	case packets.IUnknownPlsIgnore2:
		return true, nil

	// this shouldn't be received by the login server but sometimes it happens
	case packets.IPlayerUpdateIgnore:
		return true, nil

	case packets.ILoginPassword:
		return handleLoginPassword(con, it)

	case packets.IAfterLogin:
		return handleAfterLogin(con, it)

	case packets.IServerListRequest, packets.IServerListRerequest:
		return handleServerListRequest(con)

	case packets.IServerStatusRequest:
		return handleServerStatusRequest(con, it)

	case packets.IViewAllChar:
		return handleViewAllChar(con)

	case packets.IRelog:
		return handleRelog(con)

	case packets.ICharlistRequest:
		return handleCharlistRequest(con, it)

	case packets.ICharSelect:
		return handleCharSelect(con, it)

	case packets.ICheckCharName:
		return handleCheckCharName(con, it)

	case packets.ICreateChar:
		return handleCreateChar(con, it)

	case packets.IDeleteChar:
		return handleDeleteChar(con, it)

	case packets.ISetGender:
		return true, nil // not gonna use account-based gender

	case packets.IRegisterPin:
		return handleRegisterPin(con, it)

	case packets.IGuestLogin:
		return true, nil // we're gonna ignore this for now
	}

	return false, nil // forward packet to next handler
}

// TODO: split these handlers into multiple files?

// handleLoginPassword handles a login packet
func handleLoginPassword(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	// TODO split this func into smaller funcs so that it's more readable

	var successful bool = false
	var online, banned, admin, gmlevel int = 0, 0, 0, 0
	var banreason, deletepassword uint = 0, 0
	var bantime, creation int64
	var userid int32
	handled = false

	user, err := it.DecodeString()
	pass, err := it.DecodeString()
	if err != nil {
		return
	}

	ip := utils.RemoteAddrToIp(con.Conn().RemoteAddr().String())
	// we don't need the extra data

	if len(user) > consts.MaxNameSize || len(user) < consts.MinNameSize {
		err = errors.New("Invalid username size")
		return
	}

	if len(pass) > consts.MaxPasswordSize || len(user) < consts.MinPasswordSize {
		err = errors.New("Invalid password size")
		return
	}

	// look for the account in the database
	db := common.GetDB()
	st, err := db.Prepare("SELECT * FROM accounts WHERE username = ?")
	res, err := st.Run(user)
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	// column indices
	colpassword := res.Map("password")
	colsalt := res.Map("salt")
	coluserid := res.Map("id")
	colonline := res.Map("online")
	colbanned := res.Map("banned")
	colbanreason := res.Map("ban_reason")
	colbanexpire := res.Map("ban_expire")
	colcreation := res.Map("creation_date")
	coladmin := res.Map("admin")
	colgmlevel := res.Map("gm_level")
	coldeletepassword := res.Map("char_delete_password")

	handled = true

	// {autoregister begin:
	// account not found, see if we can autoregister else send login failed
	if len(rows) == 0 {
		if consts.AutoRegister {
			st, err = db.Prepare("INSERT INTO accounts(username, password, char_delete_password, creation_date) " +
				"VALUES(?, ?, 11111111, NOW())")
			_, err = st.Run(user, pass)
			// auto registrations won't hash the password right away to save server load
			// it will be hashed the first time they log in
			if err != nil {
				handled = false
				return
			}

			// get the data we just inserted
			st, err = db.Prepare("SELECT * FROM accounts WHERE username = ?")
			res, err = st.Run(user)
			rows, err = res.GetRows()
			if err != nil {
				handled = false
				return
			}
			if len(rows) == 0 {
				handled = false
				err = errors.New("Could not find account in database after auto-registration")
				return
			}

			// store account info obtained from the database
			online = rows[0].Int(colonline)
			banned = rows[0].Int(colbanned)
			banreason = rows[0].Uint(colbanreason)
			bantime = rows[0].Localtime(colbanexpire).Unix()
			creation = rows[0].Localtime(colcreation).Unix()
			admin = rows[0].Int(coladmin)
			gmlevel = rows[0].Int(colgmlevel)
			deletepassword = rows[0].Uint(coldeletepassword)
			userid = int32(rows[0].Int(coluserid))

			successful = true
		} else {
			err = con.SendPacket(packets.LoginFailed(packets.LoginNotRegistered))
		}
		// autoregister end}

		// {regular login begin
	} else {
		// check ip ban
		st, err = db.Prepare("SELECT id FROM ip_bans WHERE ip = ?")
		res, err = st.Run(ip)
		ipbanrows, iperr := res.GetRows()
		err = iperr
		if err != nil {
			handled = false
			return
		}

		if len(ipbanrows) != 0 {
			// the user is ip banned
			// I don't think this date matters
			ipbantime := time.Date(7100, time.January, 1, 0, 0, 0, 0, time.Local)
			err = con.SendPacket(packets.LoginBanned(utils.UnixToTempBanTimestamp(
				ipbantime.Unix()), packets.BanDeleted))
		} else {
			// store account info obtained from the database
			dbpassword := rows[0].Str(colpassword)
			dbsalt := rows[0].Str(colsalt)

			userid = int32(rows[0].Int(coluserid))
			online = rows[0].Int(colonline)
			banned = rows[0].Int(colbanned)
			banreason = rows[0].Uint(colbanreason)
			bantime = rows[0].Localtime(colbanexpire).Unix()
			creation = rows[0].Localtime(colcreation).Unix()
			admin = rows[0].Int(coladmin)
			gmlevel = rows[0].Int(colgmlevel)
			deletepassword = rows[0].Uint(coldeletepassword)

			switch {
			// unhashed password, hash and accept login if correct
			case len(dbsalt) == 0: // empty string = NULL
				if pass != dbpassword {
					// the unhashed password is invalid
					err = con.SendPacket(packets.LoginFailed(packets.LoginIncorrectPassword))
				} else {
					// the unhashed password is valid, hash it
					newsalt := utils.MakeSalt()
					hashedpass := utils.HashPassword(pass, newsalt)

					st, err = db.Prepare("UPDATE accounts SET password = ?, salt = ? WHERE id = ?")
					_, err = st.Run(hashedpass, newsalt, userid)
					if err != nil {
						handled = false
						return
					}
					successful = true
				}

			// regularly hashed password that matches the account's password
			case utils.HashPassword(pass, dbsalt) == dbpassword:
				successful = true

			// invalid password
			default:
				err = con.SendPacket(packets.LoginFailed(packets.LoginIncorrectPassword))
			}
		}
	}
	// regular login end}

	// correct info but the account is already logged in
	if successful && online > 0 {
		err = con.SendPacket(packets.LoginFailed(packets.LoginAlreadyLoggedIn))
		successful = false
	}

	// correct info but the account is banned
	if successful && banned > 0 {
		err = con.SendPacket(packets.LoginBanned(utils.UnixToTempBanTimestamp(bantime), byte(banreason)))
		successful = false
	}

	// unsuccessful login
	if !successful {
		con.RegisterInvalidLogin() // increase failed login counter

		// drop the user for too many failed attempts
		if consts.MaxLoginFails != 0 && con.InvalidLogins() > consts.MaxLoginFails {
			handled = false
			err = errors.New("Too many failed log-in attempts.")
		}
		return
	}

	st, err = db.Prepare("UPDATE accounts SET last_login = NOW() WHERE id = ?")
	_, err = st.Run(userid)
	if err != nil {
		handled = false
		return
	}

	con.SetPlayerStatus(client.LoggedIn)
	con.SetId(userid)

	// TODO: check silence

	con.SetAccountCreationTime(creation)
	con.SetCharDeletePassword(uint32(deletepassword))
	con.SetAdmin(admin > 0)
	con.SetGmLevel(int32(gmlevel))

	// confirm successful login
	err = con.SendPacket(packets.AuthSuccessRequestPin(user))
	fmt.Println(ip, "logged in")
	fmt.Println(con)

	handled = err == nil
	return
}

// handleAfterLogin handles an after-login (pin) packet
func handleAfterLogin(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	var pepperoni, pizza byte
	handled = false

	pepperoni, err = it.Decode1()
	pizza, err = it.Decode1()
	if err != nil {
		return
	}

	if con.PlayerStatus() != client.LoggedIn {
		err = errors.New(fmt.Sprintf(
			"Unexpected after login packet with player status = %s",
			con.PlayerStatusString()))
		return
	}

	if pepperoni > 0 && pizza > 0 {
		err = con.SendPacket(packets.PinAccepted()) // pins are for faggots
	} else {
		err = errors.New("Invalid pin packet when pins are unimplemented")
	}

	handled = err == nil
	return
}

// handleServerListRequest handles a server list request packet by sending the world and channel list
func handleServerListRequest(con *client.Connection) (handled bool, err error) {
	handled = false
	if con.PlayerStatus() != client.LoggedIn {
		err = errors.New(fmt.Sprintf(
			"Tried to request worlds with invalid player status %s",
			con.PlayerStatus()))
		return
	}

	worlds.Lock()
	defer worlds.Unlock()
	err = worlds.Show(con)
	handled = err == nil
	return
}

// handleServerStatusRequest handles a world selection request by sending the world load
func handleServerStatusRequest(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	if con.PlayerStatus() != client.LoggedIn {
		err = errors.New(fmt.Sprintf(
			"Tried to select world with invalid player status %s",
			con.PlayerStatus()))
		return
	}

	worldId, err := it.Decode1s()
	if err != nil {
		return
	}

	worlds.Lock()
	defer worlds.Unlock()
	world := worlds.Get(worldId)

	if world == nil {
		err = errors.New("Selected an invalid world")
		return
	}

	con.SetWorldId(worldId)
	fmt.Println(con.Conn().RemoteAddr(), "selected world", worldId)

	servstatus := uint16(packets.ServerNormal)

	switch {
	case world.PlayerLoad() == world.Conf().MaxPlayerLoad():
		servstatus = packets.ServerFull

	case world.PlayerLoad() >= int32(float64(world.Conf().MaxPlayerLoad())*0.9):
		servstatus = packets.ServerHigh
	}

	err = con.SendPacket(packets.ServerStatus(servstatus))
	handled = err == nil
	return
}

// sendWorldAllChars returns a packet that sends the characters list for one world to the client
// when "show all chars" has been requested
func sendWorldAllChars(worldId int8, charlist []*common.CharData) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(packets.OAllCharlist)
	p.Encode1(0x00)
	p.Encode1s(worldId)
	p.Encode1(byte(len(charlist)))

	// encode all characters
	for _, char := range charlist {
		char.Encode(&p)
	}

	return
}

// handleViewAllChar sends the character list when a client click "view all chars"
func handleViewAllChar(con *client.Connection) (handled bool, err error) {
	handled = false
	if con.PlayerStatus() != client.LoggedIn {
		err = errors.New(fmt.Sprintf(
			"Tried to get all charlists with invalid player status %s",
			con.PlayerStatus()))
		return
	}

	// get user's chars from the database
	db := common.GetDB()
	st, err := db.Prepare("SELECT * FROM characters WHERE user_id = ?")
	res, err := st.Run(con.Id())
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	colworldid := res.Map("world_id")
	charcount := uint32(0)
	charmap := make(map[int8][]*common.CharData) // char list of each world mapped by world id

	worlds.Lock()
	defer worlds.Unlock()

	// loop chars in rows and append to the map
	// TODO: check if order counts
	for _, row := range rows {
		// get world id and make sure that it's online
		worldId := int8(row.Int(colworldid))
		w := worlds.Get(worldId)
		if w == nil || !w.Connected() {
			// ignore char as the world it's on is offline
			continue
		}

		// append character to the map
		var cdata *common.CharData
		cdata, err = common.GetCharDataFromDBRow(row, res)
		if err != nil {
			return
		}
		charmap[worldId] = append(charmap[worldId], cdata)
		charcount++ // increase valid char count
	}

	// this probabilly indicates the last character slot that will be visible
	// <= 3 chars = 3
	// <= 6 chars = 6
	// <= 9 chars = 9
	// and so on
	unk := charcount + (3 - charcount%3)
	err = con.SendPacket(packets.SendAllCharsBegin(uint32(len(charmap)), unk))
	if err != nil {
		return
	}

	// iterate the valid characters map and send them to the user
	for worldId, charList := range charmap {
		err = con.SendPacket(sendWorldAllChars(worldId, charList))
		if err != nil {
			return
		}
	}

	handled = true
	return
}

// handleRelog handles a relog request
func handleRelog(con *client.Connection) (handled bool, err error) {
	handled = false
	err = con.SendPacket(packets.RelogResponse())
	handled = err == nil
	return
}

// sendWorldChars returns a packet that sends the characters list for one world to the client
// after the user selects a channel
func sendWorldChars(charlist []*common.CharData, maxchars uint32) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(packets.OCharList)
	p.Encode1(0x00)

	// encode all characters
	p.Encode1(byte(len(charlist)))
	for _, char := range charlist {
		char.Encode(&p)
	}

	p.Encode4(maxchars)
	return
}

// handleCharlistRequest handles a character list request for a certain channel
func handleCharlistRequest(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	if con.PlayerStatus() != client.LoggedIn {
		err = errors.New(fmt.Sprintf(
			"Tried to get charlist with invalid player status %s",
			con.PlayerStatus()))
		return
	}

	// get world / channel from recv
	clientWorldId, err := it.Decode1s()
	channelId, err := it.Decode1s()
	if err != nil {
		return
	}

	// hack / error checking
	if clientWorldId != con.WorldId() {
		err = errors.New(fmt.Sprintf("Selected a channel on a different "+
			"world than the currently selected one (got world %d, expected %d)",
			clientWorldId, con.WorldId()))
		return
	}

	worlds.Lock()
	defer worlds.Unlock()

	// check if world is valid
	w := worlds.Get(con.WorldId())
	if w == nil {
		err = errors.New("Tried to select a channel before selecting a world")
		return
	}

	// check if channel is online / exists
	ch := w.Channel(channelId)
	if ch == nil {
		err = errors.New(fmt.Sprintf("Tried to select channel %d "+
			"on world %d, but the channel is offline "+
			"or does not exist", channelId, con.WorldId()))
		return
	}

	// we can now safely assume that the user correctly selected this channel
	con.SetChannel(channelId)
	fmt.Println(con.Conn().RemoteAddr(), "selected channel",
		channelId, "on world", con.WorldId())

	// get the user's characters on this world
	db := common.GetDB()
	st, err := db.Prepare("SELECT * FROM characters WHERE user_id = ? AND world_id = ?")
	res, err := st.Run(con.Id(), con.WorldId())
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	chars := make([]*common.CharData, len(rows))

	for i, row := range rows {
		// append character to the array
		var cdata *common.CharData
		cdata, err = common.GetCharDataFromDBRow(row, res)
		if err != nil {
			return
		}
		chars[i] = cdata
	}

	// get max character slots
	st, err = db.Prepare("SELECT char_slots FROM storage WHERE user_id = ? AND world_id = ?")
	res, err = st.Run(con.Id(), con.WorldId())
	rows, err = res.GetRows()
	colcharslots := res.Map("char_slots")
	if err != nil {
		return
	}

	var maxslots uint32
	if len(rows) > 0 {
		maxslots = uint32(rows[0].Int(colcharslots))
	} else {
		maxslots = uint32(consts.InitialCharSlots)
	}

	// send character list
	err = con.SendPacket(sendWorldChars(chars, maxslots))
	handled = err == nil
	return
}

// handleCharSelect handles a char selection packet by initiating a server transfer
func handleCharSelect(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	charId, err := it.Decode4s()
	if err != nil {
		return
	}

	// select char when not logged in
	if con.PlayerStatus() != client.LoggedIn {
		err = errors.New(fmt.Sprintf(
			"Tried to select character with invalid player status %s",
			con.PlayerStatus()))
		return
	}

	// select other people's chars
	if !validators.OwnsCharacter(con, charId) {
		err = errors.New(fmt.Sprintf(
			"Tried to select character %d which he doesn't own",
			charId))
		return
	}

	worlds.Lock()
	defer worlds.Unlock()

	w := worlds.Get(con.WorldId())

	if w == nil || !w.Connected() {
		err = errors.New("Selected a character in an invalid world")
		return
	}

	// TODO: match user's subnet and connect to 127.0.0.1 if they are on the same subnet

	chanIp := make([]byte, 4) // maple doesn't support ipv6 :(
	port := int16(-1)

	ch := w.Channel(con.Channel())

	if ch == nil {
		// TODO: find out the channel closed error packet header
		fmt.Println(con.Conn().RemoteAddr(), "tried to connect to an offline channel")
		handled = true
		return
	}

	port = ch.Port()
	// TODO: resolve this to the external ip address to actually make it work online
	// FIXME
	chanIp = ch.Ip()
	if len(chanIp) != 4 {
		err = errors.New("Ipv6 not supported")
		return
	}

	charip := utils.RemoteAddrToBytes(con.Conn().RemoteAddr().String())
	w.WorldCon().SendPacket(interserver.MessageToChannel(con.Channel(), interserver.PlayerJoiningChannel(charId, charip)))
	err = con.SendPacket(packets.ConnectIp(chanIp, port, charId))
	handled = err == nil
	return
}

// handleCheckCharName handles a char name check request packet
func handleCheckCharName(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	used := true

	name, err := it.DecodeString()
	if err != nil {
		return
	}

	namelen := len(name)
	if namelen < consts.MinNameSize || namelen > consts.MaxNameSize {
		err = errors.New(fmt.Sprintf("Name %s has invalid length of %d", name, namelen))
		return
	}

	switch {
	case !validators.ValidName(name), validators.NameTaken(name):
		used = true
	default:
		used = false
	}

	err = con.SendPacket(packets.CharNameResponse(name, used))
	handled = err == nil
	return
}

// sendChar returns a packet that sends the information for a newly created character
func sendChar(char *common.CharData) (p maplelib.Packet) {
	p = packets.NewEncryptedPacket(packets.OAddNewCharEntry)
	p.Encode1(0x00) // idk what this byte does
	char.Encode(&p)
	return
}

// handleCreateChar handles a character creation packet
func handleCreateChar(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false

	// create char when not logged in
	if con.PlayerStatus() != client.LoggedIn {
		err = errors.New(fmt.Sprintf(
			"Tried to create character with invalid player status %s",
			con.PlayerStatus()))
		return
	}

	name, err := it.DecodeString()
	face, err := it.Decode4s()
	hair, err := it.Decode4s()
	haircolor, err := it.Decode4s()

	tmp1, err := it.Decode4s()
	skincolor := int8(tmp1)

	top, err := it.Decode4s()
	bottom, err := it.Decode4s()
	shoes, err := it.Decode4s()
	weapon, err := it.Decode4s()
	gender, err := it.Decode1s()

	// rolled stats, sum must be 25 and all must be > 4
	tmp2, err := it.Decode1s()
	str := int16(tmp2)

	tmp2, err = it.Decode1s()
	dex := int16(tmp2)

	tmp2, err = it.Decode1s()
	intt := int16(tmp2)

	tmp2, err = it.Decode1s()
	luk := int16(tmp2)

	if err != nil {
		return
	}

	// TODO: see if it's possible to roll stats server-side by modding the client or something

	// name length check
	namelen := len(name)
	if namelen < consts.MinNameSize || namelen > consts.MaxNameSize {
		err = errors.New(fmt.Sprintf("Name %s has invalid length of %d", name, namelen))
		return
	}

	// forbidden name check
	if !validators.ValidName(name) {
		err = errors.New(fmt.Sprintf("Name %s is forbidden", name))
		return
	}

	// stat roll check
	if !validators.ValidRoll(str, dex, intt, luk) {
		err = errors.New(fmt.Sprintf("Invalid dice roll of %d, %d, %d, %d", str, dex, intt, luk))
		return
	}

	// equips / look check
	if !validators.ValidNewCharacter(face, hair, haircolor, skincolor,
		top, bottom, shoes, weapon, gender) {
		err = errors.New("Invalid equips/look")
		return
	}

	// all data has been validated, the character can be safely created
	db := common.GetDB()
	st, err := db.Prepare("INSERT INTO characters(name, user_id, world_id, " +
		"face, hair, skin, gender, str, dex, `int`, luk) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	res, err := st.Run(name, con.Id(), con.WorldId(), face, hair+haircolor, skincolor, gender, str, dex, intt, luk)
	if err != nil {
		return
	}

	charid := int32(res.InsertId())

	// create equips
	err = items.Create(con, top, charid, -consts.EquipTop)
	err = items.Create(con, bottom, charid, -consts.EquipBottom)
	err = items.Create(con, shoes, charid, -consts.EquipShoe)
	err = items.Create(con, weapon, charid, -consts.EquipWeapon)
	err = items.Create(con, consts.BeginnersGuidebook, charid, 1)
	if err != nil {
		return
	}

	// get the newly created character's data
	st, err = db.Prepare("SELECT * FROM characters WHERE character_id = ?")
	res, err = st.Run(charid)
	rows, err := res.GetRows()
	if err != nil {
		return
	}

	if len(rows) < 1 {
		err = errors.New(fmt.Sprintf("Char id %d not found in database after creating it", charid))
		return
	}

	thechar, err := common.GetCharDataFromDBRow(rows[0], res)
	if err != nil {
		return
	}

	// send the new character's data
	err = con.SendPacket(sendChar(thechar))
	if err != nil {
		return
	}

	worlds.Lock()
	defer worlds.Unlock()

	// sync new character with the world server
	w := worlds.Get(con.WorldId())
	if w == nil {
		err = errors.New(fmt.Sprintf("The user is somehow connected to an offline world %d", con.WorldId()))
		return
	}

	handled = err == nil
	return
}

// handleDeleteChar handles a char deletion request
func handleDeleteChar(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	if con.PlayerStatus() != client.LoggedIn {
		err = errors.New(fmt.Sprintf(
			"Tried to create character with invalid player status %s",
			con.PlayerStatus()))
		return
	}

	bdaycode, err := it.Decode4()
	charid, err := it.Decode4s()
	if err != nil {
		err = con.SendPacket(packets.DeleteCharResponse(charid, packets.DeleteFail))
		handled = err == nil
		return
	}

	// trying to delete someone else's char
	if !validators.OwnsCharacter(con, charid) {
		err = con.SendPacket(packets.DeleteCharResponse(charid, packets.DeleteFail))
		handled = err == nil
		return
	}

	// DeleteOk = 0 // ok
	// DeleteFail = 1 // failed to delete character
	// DeleteInvalidCode = 12 // invalid birthday
	status := byte(packets.DeleteOk)

	db := common.GetDB()

	/*
	   // get character's world
	   st, err := db.Prepare("SELECT world_id FROM characters WHERE id = ?")
	   res, err := st.Run(charid)
	   rows, err := res.GetRows()
	   if err != nil {
	           return
	   }
	   if len(rows) == 0 {
	           err = errors.New(fmt.Sprintf("Tried to delete a character that doesn't exist (id=%d).", charid))
	           return
	   }

	   colworldid := res.Map("world_id")
	   worldid := byte(rows[0].Int(colworldid))
	*/

	// check birthday code
	if bdaycode != con.CharDeletePassword() {
		status = packets.DeleteInvalidCode
	} else {
		// TODO: remove character from guild
		// TODO: delete pets
		st, sterr := db.Prepare("DELETE FROM characters WHERE character_id = ?")
		err = sterr
		_, err = st.Run(charid)
		if err != nil {
			return
		}
	}

	// char delete response
	err = con.SendPacket(packets.DeleteCharResponse(charid, status))
	if err != nil {
		return
	}

	worlds.Lock()
	defer worlds.Unlock()

	// sync deleted character with the world server
	w := worlds.Get(con.WorldId())
	if w == nil {
		err = errors.New(fmt.Sprintf("The user is somehow connected to an offline world %d", con.WorldId()))
		return
	}

	handled = err == nil
	return
}

// handleRegisterPin handles a register pin packet
// pins are unused for now so the pin will be ignored
func handleRegisterPin(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	status, err := it.Decode1()
	if err != nil {
		return
	}

	switch status {
	case 0x00:
		err = con.SendPacket(packets.PinAssigned())

	default:
		err = errors.New(fmt.Sprintf("%d is not a valid register pin status", status))
	}

	handled = err == nil
	return
}
