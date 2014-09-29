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
	IOAuth                = 0x1001
	IOMessageToChannel    = 0x1002
	IOWorldConnect        = 0x1003
	IOLoginChannelConnect = 0x1004
	IORemoveChannel       = 0x1005
	IOChannelConnect      = 0x1006
	IORegisterChannel     = 0x1007
)

// Inter-server sync headers
const (
	IOSyncWorldCharacterCreated       = 0x8801
	IOSyncWorldCharacterDeleted       = 0x8802
	IOSyncChannelCharacterCreated     = 0x8803
	IOSyncChannelCharacterDeleted     = 0x8804
	IOSyncChannelNewPlayer            = 0x8805
	IOSyncWorldPerformChangeChannel   = 0x8806
	IOSyncChannelPerformChangeChannel = 0x8807
	IOSyncWorldLoadCharacter          = 0x8808
	IOSyncChannelUpdatePlayer         = 0x8809
)
