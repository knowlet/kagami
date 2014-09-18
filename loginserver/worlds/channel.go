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

package worlds

// A Channel holds information about a single channel, such as the port and the population
type Channel struct {
	port       uint32
	population int32
}

func (c *Channel) Port() uint32            { return c.port }
func (c *Channel) SetPort(port uint32)     { c.port = port }
func (c *Channel) Population() int32       { return c.population }
func (c *Channel) SetPopulation(pop int32) { c.population = pop }
