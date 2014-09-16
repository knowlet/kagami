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

import "fmt"

import (
	"github.com/Francesco149/kagami/common/consts"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

var db mysql.Conn = nil

func GetDB() mysql.Conn {
	if db == nil {
		fmt.Println("Connecting to database on", consts.MySQLHost)
		db = mysql.New("tcp", "", consts.MySQLHost, consts.MySQLUser,
			consts.MySQLPassword, consts.MySQLDB)
		err := db.Connect()
		if err != nil {
			panic(err)
		}
		fmt.Println("Connected!")
	}

	return db
}
