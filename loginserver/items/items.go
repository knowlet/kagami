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

// Package items contains various utilities to manage items
package items

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/loginserver/client"
)

// getItemInventory returns the item's inventory (equip, use...)
func getItemInventory(itemId int32) int8 { return int8(itemId / 1000000) }

// Create adds an item to a character's inventory
func Create(con *client.Connection, id, charid int32, slot int16) (err error) {
	itype := getItemInventory(id)

	// TODO: obtain item info from wz files

	db := common.GetDB()
	st, err := db.Prepare("INSERT INTO items(inv, slot, location, user_id, world_id, item_id, character_id) " +
		"VALUES(?, ?, 'inventory', ?, ?, ?, ?)")
	_, err = st.Run(itype, slot, con.Id(), con.WorldId(), id, charid)
	return
}
