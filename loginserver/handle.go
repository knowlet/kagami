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
)

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/kagami/common/packets"
	"github.com/Francesco149/maplelib"
)

// Handle handles loginserver packets
func Handle(con common.Connection, p maplelib.Packet) (handled bool, err error) {
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
func handleLoginPassword(con common.Connection, it maplelib.PacketIterator) (handled bool, err error) {
	var successful bool = false
	var online int = 0
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

	handled = true

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

			successful = true
		} else {
			err = con.SendPacket(packets.LoginFailed(packets.LoginNotRegistered))
		}
	} else {
		// check ban
		// TODO

		// check password
		dbpassword := rows[0].Str(res.Map("password"))
		dbsalt := rows[0].Str(res.Map("salt"))
		userid := rows[0].Int(res.Map("id"))
		online = rows[0].Int(res.Map("online"))

		switch {
		// unhashed password, hash and accept login if correct
		case len(dbsalt) == 0:
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

	// todo: check ban

	if successful && online > 0 {
		err = con.SendPacket(packets.LoginFailed(packets.LoginAlreadyLoggedIn))
		successful = false
	}

	if !successful {
		// TODO: increase failed login counter and disconnect if they are too many
		return
	}

	// TODO: set player status as loggedin
	// TODO: store useful data such as creation time, deletion password etc for the player

	err = con.SendPacket(packets.AuthSuccessRequestPin(user))
	fmt.Println(ip, "logged in")
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
