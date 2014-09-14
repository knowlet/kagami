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
        "net"
        "time"
        "math/rand"
        "io"
        "bytes"
        "encoding/binary"
)

import (
        "github.com/Francesco149/kagami/common/packets"
)

func clientLoop(con net.Conn) {
        defer con.Close()
        
        ivrecv := make([]byte, 4)
        ivsend := make([]byte, 4)
         
        binary.LittleEndian.PutUint32(ivrecv, rand.Uint32())
        binary.LittleEndian.PutUint32(ivsend, rand.Uint32())
        hello := packets.Handshake(62, ivsend, ivrecv, false)
        
        fmt.Printf("Sending hello packet: %v\n", hello)
        io.Copy(con, bytes.NewReader(hello))
        
        for {
                time.Sleep(100 * time.Millisecond)        
        }
}

func main() {
        rand.Seed(time.Now().UnixNano())
        
        sock, err := net.Listen("tcp", ":8484")
        if err != nil { 
                fmt.Println("Failed to create socket: ", err)
                return 
        }
        
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