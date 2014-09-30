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

// Package intersever contains inter-server packet headers and builders
package interserver

// Inter-server headers
const (
	IOAuth                    = 0x1001
	IOWorldConnect            = 0x1003
	IOLoginChannelConnect     = 0x1004
	IORemoveChannel           = 0x1005
	IOChannelConnect          = 0x1006
	IORegisterChannel         = 0x1007
	IOSyncPlayerJoinedChannel = 0x1008
	IOSyncPlayerLeftChannel   = 0x1009
	IOSyncChannelPopulation   = 0x1010
)
