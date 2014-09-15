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
        "bufio"
)

import (
        "github.com/Francesco149/maplelib"
        "github.com/Francesco149/kagami/common/packets"
        "github.com/Francesco149/kagami/common/consts"
)

// TODO: make a generic server package for packet handling

// Returned when an I/O error occurs while reading data from the socket
type IOError struct {
        BytesRead int 
        Err error 
}
func (e IOError) Error() string {
        return fmt.Sprintf("Could only read %d bytes. Err = %v.", e.BytesRead, e.Err)   
}

// Returned when a received packet has invalid size specified in the 
// encrypted header
type InvalidPacketError int
func (e InvalidPacketError) Error() string {
        return fmt.Sprintf("Recieved invalid packet of size %d.", int(e))        
}

// Sends a packet through the given connection
func SendPacket(con net.Conn, packet maplelib.Packet) {
        io.Copy(con, bytes.NewReader(packet))
}

func RecvPacket(con net.Conn, c *maplelib.Crypt) (packet maplelib.Packet, err error) {
        var plen int = 0
        
        r := bufio.NewReader(con)
        
        // encrypted header (4 bytes)
        p := make([]byte, consts.EncryptedHeaderSize)
        
        n, err := r.Read(p)
        
        if n != consts.EncryptedHeaderSize || err != nil {
                packet, err = nil, IOError{n, err}
                return
        }
        
        fmt.Printf("Received encrypted header % X\n", p)
        plen = maplelib.GetPacketLength(p)
        fmt.Printf("Packet length is %d\n", plen)
        
        if plen < 2 {
                packet, err = nil, InvalidPacketError(plen)
                return
        }
        
        // data
        data := make([]byte, plen)
        r.Read(data)
        c.Decrypt(data)
        c.Shuffle()
        
        packet, err = maplelib.Packet(data), nil
        return
}

// Sends the handshake and handles packets for a single client
func clientLoop(con net.Conn) {
        var ivrecv, ivsend [4]byte
        
        defer con.Close()
         
        binary.LittleEndian.PutUint32(ivrecv[:], rand.Uint32())
        binary.LittleEndian.PutUint32(ivsend[:], rand.Uint32())
        hs := packets.Handshake(62, ivsend, ivrecv, false)
        
        fmt.Printf("Sending handshake: %v\n", hs)
        SendPacket(con, hs)
        
        send := maplelib.NewCrypt(ivsend, consts.MapleVersion)
        recv := maplelib.NewCrypt(ivrecv, consts.MapleVersion)
        
        fmt.Println("ivsend:", send)
        fmt.Println("ivrecv:", recv)
        
        for {
                inpacket, err := RecvPacket(con, &recv)
                if err != nil {
                        fmt.Println(err)
                        break
                }
                
                fmt.Println("Decrypted packet:", inpacket)
                time.Sleep(100 * time.Millisecond)        
        }
        
        fmt.Println("Dropping: ", con.RemoteAddr())
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