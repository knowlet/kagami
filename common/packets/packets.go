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

package packets

import "github.com/Francesco149/maplelib"

// Handshake returns a handshake packet that must be sent UNENCRYPTED to newly connected clients
// The initialization vectors ivsend and ivrecv are 4 bytes, any extra data will be ignored
func Handshake(mapleVersion uint16, ivsend []byte, 
        ivrecv []byte, testserver bool) (p maplelib.Packet) {
        
        testbyte := byte(8)
        if testserver {
                testbyte = 5
        }
        
        p = maplelib.NewPacket()
        p.Encode2(OHandshake) // header
        p.Encode2(mapleVersion) // game version
        p.Encode2(0x0000) // dunno maybe version is a dword
        p.Append(ivrecv[:4])
        p.Append(ivsend[:4])
        p.Encode1(testbyte) // 5 = test server, else 8
        return
}

// Ping returns a ping packet with a placeholder for the encrypted header
func Ping() (p maplelib.Packet) {
        p = maplelib.NewPacket()
        p.Encode4(0x00000000) // placeholder for the encrypted header
        p.Encode2(OPing) // header      
        return
}