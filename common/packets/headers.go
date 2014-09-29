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

// Package packets contains all packet headers GMS v62 and packet builders
// that do not require server-specific packages
package packets

// Send packet headers
const (
	// common
	OPing = 0x0011

	// login server
	OLoginStatus        = 0x0000
	OServerStatus       = 0x0003
	OPinOperation       = 0x0006
	OAllCharlist        = 0x0008
	OServerList         = 0x000A
	OCharList           = 0x000B
	OServerIP           = 0x000C
	OCharNameResponse   = 0x000D
	OAddNewCharEntry    = 0x000E
	ODeleteCharResponse = 0x000F
	ORelogResponse      = 0x0016
	OGenderDone         = 0x0004
	OPinAssigned        = 0x0007

	// channel server
	OConnectData   = 0x005C // warp to map
	OServerMessage = 0x0041
	OChangeChannel = 0x0010
)

// Recv packet headers
const (
	// common
	IPong = 0x0018

	// login server
	ILoginPassword       = 0x0001
	IAfterLogin          = 0x0009
	IServerListRequest   = 0x000B
	IServerListRerequest = 0x0004
	IServerStatusRequest = 0x0006
	IViewAllChar         = 0x000D
	IRelog               = 0x001C
	ICharlistRequest     = 0x0005
	ICharSelect          = 0x0013
	ICheckCharName       = 0x0015
	ICreateChar          = 0x0016
	IDeleteChar          = 0x0017
	IPickAllChar         = 0x000E // unused, idk what this does
	ISetGender           = 0x0008
	IRegisterPin         = 0x000A
	IGuestLogin          = 0x0002
	IUnknownPlsIgnore1   = 0x001A // this gets spammed while on login screen, apparently it means client error
	IUnknownPlsIgnore2   = 0x000F
	IPlayerUpdateIgnore  = 0x00C0

	// channel server
	ILoadCharacter = 0x0014
)
