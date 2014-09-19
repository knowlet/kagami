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

// Package validators contains various utilities to validate data
package validators

import "fmt"

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/consts"
	"github.com/Francesco149/kagami/loginserver/client"
)

// OwnsCharacter checks if the user owns the given character id
func OwnsCharacter(con *client.Connection, charId int32) bool {
	db := common.GetDB()
	st, err := db.Prepare("SELECT 1 FROM characters WHERE id = ? AND user_id = ? LIMIT 1")
	res, err := st.Run(charId, con.Id())
	rows, err := res.GetRows()
	if err != nil {
		fmt.Println("ownsCharacter:", err)
		return false
	}

	return len(rows) > 0
}

// ValidName checks if the given name is not forbidden
func ValidName(name string) bool {
	// TODO: check list of forbidden names and curse words
	return true
}

// NameTaken checks if the given character name is already taken
func NameTaken(name string) bool {
	db := common.GetDB()
	st, err := db.Prepare("SELECT 1 FROM characters WHERE name = ? LIMIT 1")
	res, err := st.Run(name)
	rows, err := res.GetRows()
	if err != nil {
		fmt.Println("nameTaken: ", err)
		return false
	}
	return len(rows) > 0
}

// ValidRoll checks if a given stat roll is valid (sum must be 25, none of the stats must be < 4)
func ValidRoll(str, dex, intt, luk int16) bool {
	switch {
	case
		str+dex+intt+luk != 25,
		str < 4,
		dex < 4,
		intt < 4,
		luk < 4:
		return false
	}

	return true
}

// ValidNewCharacter checks if the equips and look of a new character are valid
// TODO: generic item check function that checks items and class in wz files
func ValidNewCharacter(face, hair, haircolor int32, skincolor int8,
	top, bottom, shoes, weapon int32, gender int8) bool {
	switch {
	case
		!consts.BeginnerFaces[gender][face],
		!consts.BeginnerHairstyles[gender][hair],
		!consts.BeginnerTops[gender][top],
		!consts.BeginnerBottoms[gender][bottom],
		skincolor < consts.BeginnerMinSkinColor || skincolor > consts.BeginnerMaxSkinColor,
		!consts.BeginnerWeapons[weapon],
		!consts.BeginnerHairColors[haircolor]:
		return false
	}

	return true
}
