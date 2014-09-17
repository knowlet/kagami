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
	"strings"
	"time"
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/kagami/loginserver/client"
	"github.com/Francesco149/maplelib"
)

// Handle handles loginserver packets
func Handle(con *client.Connection, p maplelib.Packet) (handled bool, err error) {
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	switch header {
	case packets.ILoginPassword:
		return handleLoginPassword(con, it)

	case packets.IUnknownPlsIgnore:
		return true, nil

	case packets.IRegisterPin:
		return handleRegisterPin(con, it)
	}

	return false, nil // forward packet to next handler
}

// handleLoginPassword handles a login packet
func handleLoginPassword(con *client.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	// TODO split this func into smaller funcs so that it's more readable

	var successful bool = false
	var online, banned, admin, gmlevel int = 0, 0, 0, 0
	var banreason, deletepassword uint = 0, 0
	var bantime, creation int64
	handled = false

	user, err := it.DecodeString()
	pass, err := it.DecodeString()
	if err != nil {
		return
	}

	ip := strings.Split(con.Conn().RemoteAddr().String(), ":")[0]
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
			st, err = db.Prepare("INSERT INTO accounts(username, password, char_delete_password) VALUES(?, ?, 11111111)")
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
			ipbantime := time.Date(7100, time.January, 1, 0, 0, 0, 0, time.Local)
			err = con.SendPacket(packets.LoginBanned(common.UnixToTempBanTimestamp(ipbantime.Unix()), packets.BanDeleted))
		} else {
			// store account info obtained from the database
			dbpassword := rows[0].Str(colpassword)
			dbsalt := rows[0].Str(colsalt)
			userid := rows[0].Int(coluserid)

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
					newsalt := common.MakeSalt()
					hashedpass := common.HashPassword(pass, newsalt)

					st, err = db.Prepare("UPDATE accounts SET password = ?, salt = ? WHERE id = ?")
					_, err = st.Run(hashedpass, newsalt, userid)
					if err != nil {
						handled = false
						return
					}
					successful = true
				}

			// regularly hashed password that matches the account's password
			case common.HashPassword(pass, dbsalt) == dbpassword:
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
		err = con.SendPacket(packets.LoginBanned(common.UnixToTempBanTimestamp(bantime), byte(banreason)))
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

	con.SetPlayerStatus(client.LoggedIn)

	// TODO: check silence

	con.SetAccountCreationTime(creation)
	con.SetCharDeletePassword(uint32(deletepassword))
	con.SetAdmin(admin > 0)
	con.SetGmLevel(int32(gmlevel))

	// confirm successful login
	err = con.SendPacket(packets.AuthSuccessRequestPin(user))
	fmt.Println(ip, "logged in")
	fmt.Println(con)
	return
}

// handleLoginPassword handles a register pin packet
// pins are unused for now so the pin will be ignored
func handleRegisterPin(con common.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	handled = false
	status, err := it.Decode1()
	if err != nil {
		return
	}

	handled = true

	switch status {
	case 0x00:
		err = con.SendPacket(packets.PinAssigned())

	default:
		handled = false
		err = errors.New(fmt.Sprintf("%d is not a valid register pin status", status))
	}

	return
}
