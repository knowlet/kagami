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

import "net"

// Listen creates and returns a TCP listener on the given port
func Listen(ipport string) (sock net.Listener, err error) {
	return net.Listen("tcp", ipport)
}

// Listen connects to a TCP server on the given port
func Dial(ipport string) (con net.Conn, err error) {
        return net.Dial("tcp", ipport)
}
