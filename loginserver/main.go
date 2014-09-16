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
        "fmt"
        "time"
        "math/rand"
        "net"
)

import "github.com/Francesco149/kagami/common"

// clientLoop sends the handshake and handles packets for a single client in a loop
func clientLoop(basecon net.Conn) {
        defer basecon.Close()
        con := common.NewEncryptedConnection(basecon, false)
        
        for {
                inpacket, err := con.RecvPacket()
                if err != nil {
                        fmt.Println(err)
                        break
                }
                
                handled, err := common.Handle(con, inpacket)
                if err != nil {
                        fmt.Println(err)
                        break        
                }
                
                // TODO: handle loginserver packets here
                
                if !handled {
                        fmt.Println("Unhandled packet")
                        break
                }
        }
        
        fmt.Println("Dropping: ", con.Conn().RemoteAddr())
}

func main() {
        const loginport = 8484 // TODO: config file
        
        rand.Seed(time.Now().UnixNano())
        
        fmt.Println("Kagami Pre-Alpha")
        fmt.Println("Initializing LoginServer...")
        
        sock, err := common.NewTcpServer(fmt.Sprintf(":%d", loginport))
        if err != nil { 
                fmt.Println("Failed to create socket: ", err)
                return 
        }
        
        fmt.Println("Listening on port", loginport)
        
        for {
                con, err := sock.Accept()
                if err != nil {
                        fmt.Println("Failed to accept connection: ", err)
                        return        
                }
                
                fmt.Println("Accepted: ", con.RemoteAddr())
                go clientLoop(con)
        }
}