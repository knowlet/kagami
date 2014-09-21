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

package channels

// Channel holds information about one channel and its population
type Channel struct {
        id int8 
        port int16 
        population int32 
        con *Connection
}

func NewChannel(ccon *Connection, cid int8, cport int16) *Channel {
        return &Channel {
                id: cid, 
                port: cport, 
                population: 0, 
                con: ccon, 
        }
}

func (c *Channel) IncPopulation() { c.population++ }
func (c *Channel) DecPopulation() { c.population-- }
func (c *Channel) Population() int32 { return c.population }
func (c *Channel) SetConn(con *Connection) { c.con = con }
func (c *Channel) Conn() *Connection { return c.con }
func (c *Channel) SetPort(port int16) { c.port = port }
func (c *Channel) Port() int16 { return c.port }