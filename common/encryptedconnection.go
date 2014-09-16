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

// Package common contains all of the code that can be reused across login, world and channel server
package common

import (
        "fmt"
        "net"
        "math/rand"
        "encoding/binary"
        "time"
        "errors"
)

import (
        "github.com/Francesco149/maplelib"
        "github.com/Francesco149/kagami/common/packets"
        "github.com/Francesco149/kagami/common/consts"
)

// A IOError is returned when an I/O error occurs while reading/writing data from the socket
type IOError struct {
        bytesRead int 
        err error 
}
func (e IOError) Error() string {
        return fmt.Sprintf("Could only read/write %d bytes. Err = %v.", e.bytesRead, e.err)   
}

// An InvalidPacketError is returned when a received packet has invalid size specified in the 
// encrypted header
type InvalidPacketError int
func (e InvalidPacketError) Error() string {
        return fmt.Sprintf("Recieved invalid packet of size %d.", int(e))        
}

// An EncryptedConnection represent an individual client connected to our socket 
// that will send and receive MapleStory-encrypted packets
type EncryptedConnection struct {
        con net.Conn
        send maplelib.Crypt
        recv maplelib.Crypt
        pinged bool
        lastping int64
        lastactive int64
}

// Checks that the given error is not nil and is a timeout
func isTimeout(err error) bool {
        // x.(T) asserts that x is not nil and the value stored in x is of type T
        if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
                return true        
        }
        
        return false
}

func (c *EncryptedConnection) renewSendTimeout() {
        c.Conn().SetWriteDeadline(time.Now().Add(consts.ClientTimeout * time.Second))        
}

// sendHandshake sends the handshake packet with the encryption keys to the client
func (c *EncryptedConnection) sendHandshake(isTestServer bool) error {
        hs := packets.Handshake(consts.MapleVersion, c.SendCrypt().IV()[:4], 
                c.RecvCrypt().IV()[:4], isTestServer)
        
        c.renewSendTimeout()
        n, err := c.Conn().Write(hs)
        if isTimeout(err) {
                return errors.New("Write timeout")        
        }
        if err != nil {
                return IOError{n, err}                
        }
        
        return nil
}

// NewEncryptedConnection creates an encrypted connection around the given 
// connection and initializes the encryption by performing the handshake
func NewEncryptedConnection(con net.Conn, isTestServer bool) (c *EncryptedConnection) {
        var ivrecv, ivsend [4]byte
        
        c = &EncryptedConnection{}
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

// Ping pings the client
func (c *EncryptedConnection) Ping() error {
        if c.pinged {
                return nil        
        }
        
        fmt.Println("Pinging", c.Conn().RemoteAddr())
        c.lastping = time.Now().Unix()
        c.pinged = true
        return c.SendPacket(packets.Ping())
}

// OnPong resets the ping status and timeout time
func (c *EncryptedConnection) OnPong() error {
        if !c.pinged { // fake pong
                return errors.New(fmt.Sprintf("%v attempted to fake a pong", c.Conn().RemoteAddr()))
        }
        
        fmt.Println("Got pong from", c.Conn().RemoteAddr())
        c.pinged = false
        return nil
}

// tryRead attempts to read a packet from the connection and sends a ping if the client goes idle
func (c *EncryptedConnection) tryRead(p []byte) (err error) {
        // basically, set read timeout to a fraction the client timeout and try reading
        // for a short time multiple times while checking if it's time to ping
        
        loops := consts.ClientTimeout / consts.ClientIdle
        
        for i := 0; i < loops; i++ {
                // the client has been inactive long enough so we're gonna ping it
                if time.Now().Unix() - c.lastactive > consts.ClientIdle {
                        err = c.Ping()     
                        if err != nil {
                                return        
                        }
                }
                
                // this will make read time out
                c.Conn().SetReadDeadline(time.Now().Add(consts.ClientIdle * time.Second))
                
                // read data
                n, err := c.Conn().Read(p)
                if isTimeout(err) {
                        continue
                }
                if n != cap(p) || err != nil {
                        return IOError{n, err}
                }
                
                break // no errors
        }
        
        if isTimeout(err) {
                err = errors.New("Read timeout")
        }
        
        return
}

// RecvPacket listens for the next encrypted packet, decrypts it and returns it. 
// NOTE: the returned packet DOES NOT include the 4-byte encrypted header
func (c *EncryptedConnection) RecvPacket() (packet maplelib.Packet, err error) {
        var plen int = 0
        
        // read encrypted header
        header := make([]byte, consts.EncryptedHeaderSize)
        err = c.tryRead(header)
        if err != nil {
                return
        }
        
        // retrieve decrypted packet length
        plen = maplelib.GetPacketLength(header)
        if plen < 2 {
                packet, err = nil, InvalidPacketError(plen)
                return
        }
        
        // read packet data
        data := make([]byte, plen)
        err = c.tryRead(data)
        if err != nil {
                return
        }
        
        c.RecvCrypt().Decrypt(data)
        c.RecvCrypt().Shuffle()
        
        c.lastactive = time.Now().Unix() // reset idle timer
        
        packet, err = maplelib.Packet(data), nil
        return
}

// SendPacket encrypts and sends the given packet. NOTE: the packet must have 
// a 4 byte placeholder at the beginning for the encrypted header
func (c *EncryptedConnection) SendPacket(p maplelib.Packet) error {
        byteslice := []byte(p)
        c.SendCrypt().Encrypt(byteslice[:])
        
        c.renewSendTimeout()
        n, err :=  c.Conn().Write(p)
        if err != nil {
                return IOError{n, err}                
        }
        
        c.SendCrypt().Shuffle()
        
        return nil
}

