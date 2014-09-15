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
        "fmt"
        "net"
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

// A IOError is returned when an I/O error occurs while reading data from the socket
type IOError struct {
        bytesRead int 
        err error 
}
func (e IOError) Error() string {
        return fmt.Sprintf("Could only read %d bytes. Err = %v.", e.bytesRead, e.err)   
}

// An InvalidPacketError is returned when a received packet has invalid size specified in the 
// encrypted header
type InvalidPacketError int
func (e InvalidPacketError) Error() string {
        return fmt.Sprintf("Recieved invalid packet of size %d.", int(e))        
}

// NewTcpServer creates and returns a TCP listener on the given port
func NewTcpServer(port string) (sock net.Listener, err error) {
        return net.Listen("tcp", port)
}

// An EncryptedConnection represent an individual client connected to our socket 
// that will send and receive MapleStory-encrypted packets
type EncryptedConnection struct {
        con net.Conn
        send maplelib.Crypt
        recv maplelib.Crypt
}

// sendHandshake sends the handshake packet with the encryption keys to the client
func (c *EncryptedConnection) sendHandshake(isTestServer bool) {
        hs := packets.Handshake(consts.MapleVersion, c.SendCrypt().IV()[:4], 
                c.RecvCrypt().IV()[:4], isTestServer)
        io.Copy(c.Conn(), bytes.NewReader(hs))
}

// NewEncryptedConnection creates an encrypted connection around the given 
// connection and initializes the encryption by performing the handshake
func NewEncryptedConnection(con net.Conn, isTestServer bool) (c EncryptedConnection) {
        var ivrecv, ivsend [4]byte
        
        c.con = con
        
        // randomly generate initialization vectors
        binary.LittleEndian.PutUint32(ivrecv[:], rand.Uint32())
        binary.LittleEndian.PutUint32(ivsend[:], rand.Uint32())
        
        // init encryption
        c.send = maplelib.NewCrypt(ivsend, consts.MapleVersion)
        c.recv = maplelib.NewCrypt(ivrecv, consts.MapleVersion)
        
        c.sendHandshake(isTestServer)
        return
}

// Conn returns the connection associated with this encrypted connection
func (c *EncryptedConnection) Conn() net.Conn {
        return c.con        
}

// SendCrypt returns the send encryption key
func (c *EncryptedConnection) SendCrypt() *maplelib.Crypt {
        return &c.send
}

// RecvCrypt returns the recv decryption key
func (c *EncryptedConnection) RecvCrypt() *maplelib.Crypt {
        return &c.recv
}

// RecvPacket listens for the next encrypted packet, decrypts it and returns it. 
// NOTE: the returned packet DOES NOT include the 4-byte encrypted header
func (c *EncryptedConnection) RecvPacket() (packet maplelib.Packet, err error) {
        var plen int = 0
        
        r := bufio.NewReader(c.con)
        
        // read encrypted header (4 bytes)
        p := make([]byte, consts.EncryptedHeaderSize)
        n, err := r.Read(p)
        if n != consts.EncryptedHeaderSize || err != nil {
                packet, err = nil, IOError{n, err}
                return
        }
        
        // retrieve decrypted packet length
        plen = maplelib.GetPacketLength(p)
        if plen < 2 {
                packet, err = nil, InvalidPacketError(plen)
                return
        }
        
        // read and decrypt data
        data := make([]byte, plen)
        r.Read(data)
        c.RecvCrypt().Decrypt(data)
        c.RecvCrypt().Shuffle()
        
        packet, err = maplelib.Packet(data), nil
        return
}

// SendPacket encrypts and sends the given packet. NOTE: the packet must have 
// a 4 byte placeholder at the beginning for the encrypted header
func (c *EncryptedConnection) SendPacket(p maplelib.Packet) {
        byteslice := []byte(p)
        c.SendCrypt().Encrypt(byteslice[:])
        c.SendCrypt().Shuffle()
        io.Copy(c.Conn(), bytes.NewReader(p))
}

