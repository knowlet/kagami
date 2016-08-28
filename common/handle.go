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
	"github.com/knowlet/kagami/common/packets"
	"github.com/Francesco149/maplelib"
)

// Handle handles packets that are common to all three servers
func Handle(con Connection, p maplelib.Packet) (handled bool, err error) {
	it := p.Begin()
	header, err := it.Decode2()
	if err != nil {
		return false, err
	}

	switch {
	case
		header == packets.IPong,
		// this can only happen when the connection is a client in inter-server connections
		con.IsClient() && header == packets.OPing:
		return handlePong(con)
	}

	return false, nil // forward packet to next handler
}

func handlePong(con Connection) (handled bool, err error) {
	err = con.OnPong()
	handled = (err == nil)
	return
}
