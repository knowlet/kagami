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

import "errors"

import (
	"github.com/Francesco149/kagami/common"
	"github.com/Francesco149/kagami/common/interserver"
	"github.com/Francesco149/maplelib"
)

// Handle handles inter-server loginserver packets
func HandleInter(con *common.InterserverConnection, p maplelib.Packet) (handled bool, err error) {
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	// check auth
	if !con.Authenticated() {
		if header != interserver.IOAuth {
			return false, errors.New("Tried to send packets without being authenticated")
		}
		err = con.CheckAuth(it)
		if err != nil {
			return
		}
	}

	// TODO

	switch header {

	}

	return false, nil // forward packet to next handler
}
