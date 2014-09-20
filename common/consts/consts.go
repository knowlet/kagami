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

// Package consts contains various constants used everywhere in kagami
package consts

// These are default settings used for testing. When the server will support
// config files, these will be used as a fallback for when the config files are
// missing.
// TODO: support conf files

const MySQLUser = "kagami"         // MySQLUser is the MySQL username
const MySQLPassword = "testing"    // MySQLPassword is the MySQL password
const MySQLHost = "127.0.0.1:3306" // MySQLHost contains the ip:port of the MySQL database
const MySQLDB = "my_kagami"        // MySQLDB contains the name of the used MySQL database

const LoginPort = 8484            // Loginport is the port the Login Server will listen on
const LoginInterserverPort = 8485 // LoginInterserverPort is the port the Login Server will listen on for inter-server connections
const LoginIp = "127.0.0.1"       // LoginIp is the ip of the loginserver for inter-server connections

const InterServerPassword = "topfuckingkek" // The internal password that will be used to do inter-server communication

const MapleVersion = 62       // MapleVersion represents the required game client version
const EncryptedHeaderSize = 4 // EncryptedHeaderSize is the size in bytes of encrypted headers
const ClientTimeout = 30      // ClientTimeout is the number of seconds a client has to reply to a ping before it times out
const ClientIdle = 1          // ClientIdle is the number of seconds with no packet activity after which a client is considered idle

const MinNameSize = 4      // MinNameSize is the minimum length of a character/user name
const MaxNameSize = 12     // MaxNameSize is the maximum length of a character/user name
const MinPasswordSize = 4  // MinPasswordSize is the minimum length of a password
const MaxPasswordSize = 12 // MaxPasswordSize is the maximum length of a password

const InitialCharSlots = 3 // InitialCharSlots is how many character slots a new user has by default
const AutoRegister = false // AutoRegister defines whether it's possible to automatically register by attempting to log into a non existing account
const SaltLength = 10      // SaltLength is the length of password salts
const MaxLoginFails = 10   // MaxLoginFails is the amount of failed logins it takes to get disconnected, 0 = disabled

const InventoryTypes = 5 // InventoryTypes is the number of different inventories

// Inventory slots
const (
	EquipInventory = 1
	UseInventory   = 2
	SetupInventory = 3
	EtcInventory   = 4
	CashInventory  = 5
)

// this should be how many different equip slots there are but I'm not sure why
// it's 51 instead of 50
const EquippedSlots = 51

// Equip slots
const (
	EquipHelm            = 1
	EquipFace            = 2
	EquipEye             = 3
	EquipEarring         = 4
	EquipTop             = 5
	EquipBottom          = 6
	EquipShoe            = 7
	EquipGlove           = 8
	EquipCape            = 9
	EquipShield          = 10
	EquipWeapon          = 11
	EquipRing1           = 12
	EquipRing2           = 13
	EquipPetEquip1       = 14
	EquipRing3           = 15
	EquipRing4           = 16
	EquipPendant         = 17
	EquipMount           = 18
	EquipSaddle          = 19
	EquipPetCollar       = 20
	EquipPetLabelRing1   = 21
	EquipPetItemPouch1   = 22
	EquipPetMesoMagnet1  = 23
	EquipPetAutoHp       = 24
	EquipPetAutoMp       = 25
	EquipPetWingBoots1   = 26
	EquipPetBinoculars1  = 27
	EquipPetMagicScales1 = 28
	EquipPetQuoteRing1   = 29
	EquipPetEquip2       = 30
	EquipPetLabelRing2   = 31
	EquipPetQuoteRing2   = 32
	EquipPetItemPouch2   = 33
	EquipPetMesoMagnet2  = 34
	EquipPetWingBoots2   = 35
	EquipPetBinoculars2  = 36
	EquipPetMagicScales2 = 37
	EquipPetEquip3       = 38
	EquipPetLabelRing3   = 39
	EquipPetQuoteRing3   = 40
	EquipPetItemPouch3   = 41
	EquipPetMesoMagnet3  = 42
	EquipPetWingBoots3   = 43
	EquipPetBinoculars3  = 44
	EquipPetMagicScales3 = 45
	EquipPetItemIgnore1  = 46
	EquipPetItemIgnore2  = 47
	EquipPetItemIgnore3  = 48
	EquipMedal           = 49
	EquipBelt            = 50
)

// Sex Id's
const (
	SexMale   = 0
	SexFemale = 1
)

// BeginnerFaces contains all allowed beginner faces for each sex
// it is mapped by id so you can just check if the map contains the id
var BeginnerFaces = [2]map[int32]bool{
	{20000: true, 20001: true, 20002: true},
	{21000: true, 21001: true, 21002: true},
}

// BeginnerHairstyles contains all allowed beginner hairstyles for each sex
// it is mapped by id so you can just check if the map contains the id
var BeginnerHairstyles = [2]map[int32]bool{
	{30000: true, 30020: true, 30030: true},
	{31000: true, 31040: true, 31050: true},
}

// BeginnerTops contains all allowed beginner tops for each sex
// it is mapped by id so you can just check if the map contains the id
var BeginnerTops = [2]map[int32]bool{
	{1040002: true, 1040006: true, 1040010: true},
	{1041002: true, 1041006: true, 1041010: true},
}

// BeginnerBottoms contains all allowed beginner bottoms for each sex
// it is mapped by id so you can just check if the map contains the id
var BeginnerBottoms = [2]map[int32]bool{
	{1060006: true, 1060002: true},
	{1061002: true, 1061008: true},
}

// Beginner skin color lower and upper bounds
const (
	BeginnerMinSkinColor = 0
	BeginnerMaxSkinColor = 3
)

// BeginnerWeapons contains a list of allowed beginner weapons
// it is mapped by id so you can just check if the map contains the id
var BeginnerWeapons = map[int32]bool{1302000: true, 1322005: true, 1312004: true}

// BeginnerShoes contains a list of allowed beginner shoes
// it is mapped by id so you can just check if the map contains the id
var BeginnerShoes = map[int32]bool{1072001: true, 1072005: true, 1072037: true, 1072038: true}

// BeginnerHairColors contains a list of allowed beginner hair colors
// it is mapped by id so you can just check if the map contains the id
var BeginnerHairColors = map[int32]bool{0: true, 1: true, 2: true, 3: true, 7: true}

const BeginnersGuidebook = 4161001

// Default Worlds --------------------------------------------------------------
const WorldCount = 1 // WorldCount is the number of worlds that will connect to the loginserver

var WorldName = [WorldCount]string{"Scania"}     // WorldName contains a list of the world names
var WorldChannelCount = [WorldCount]byte{19}     // WorldChannelCount contains a list of the channel count for each world
var WorldId = [WorldCount]int8{0}                // WorldId contains a list of the world id's
var WorldRibbon = [WorldCount]byte{0}            // WorldRibbon contains a list of each world's ribbon. 0 = None, 1 = E, 2 = N, 3 = H
var WorldDefaultGMChat = [WorldCount]bool{false} // WorldDefaultGMChat contains a list of each world's GM chat enabled flag

var WorldMobExp = [WorldCount]int32{1}   // WorldMobExp is a list of each world's mob exp rate
var WorldQuestExp = [WorldCount]int32{1} // WorldQuestExp is a list of each world's quest exp rate
var WorldMeso = [WorldCount]int32{1}     // WorldMeso is a list of each world's mob meso drop rate
var WorldDrop = [WorldCount]int32{1}     // WorldDrop is a list of each world's drop rate

var WorldMaxCharSlots = [WorldCount]byte{6}        // WorldMaxCharSlots is a list of each world's maximum char slots
var WorldDefaultCharSlots = [WorldCount]byte{3}    // WorldDefaultCharSlots is a list of each world's initial char slots
var WorldDefaultStorageSlots = [WorldCount]byte{4} // WorldDefaultStorageSlots is a list of each world's max storage slots
var WorldMaxStats = [WorldCount]uint16{999}        // WorldMaxStats is a list of each world's stat limit
var WorldMaxMultiLevel = [WorldCount]byte{1}       // WorldMaxMultiLevel is a list of each world's max multiple level gain

var WorldListenPort = [WorldCount]int16{7100} // WorldListenPort is a list of each world's listen port

// WorldEventMessage is a list of each world's event message
var WorldEventMessage = [WorldCount]string{"Top fucking kek"}

// WorldScrollingHeader is a list of each world's scrolling header message
var WorldScrollingHeader = [WorldCount]string{"Totsugeki~"}

var WorldMaxPlayerLoad = [WorldCount]int32{1000} // WorldMaxPlayerLoad is a list of each world's player cap

// WorldFameDelay contains how many seconds you need to wait before you can fame someone for each world
var WorldFameDelay = [WorldCount]int64{86400}

// WorldFameDelay contains how many seconds you need to wait before you can fame the same person again for each world
var WorldFameResetTime = [WorldCount]int64{2592000}

// WorldMapUnloadTime contains the map unload time in seconds for each world
var WorldMapUnloadTime = [WorldCount]int64{3600}

// WorldPianusChannels contains a list of channels where Pianus spawns, 255 = all
var WorldPianusChannels = [WorldCount][]byte{
	[]byte{0xFF},
}

// WorldPapChannels contains a list of channels where Papulatus spawns, 255 = all
var WorldPapChannels = [WorldCount][]byte{
	[]byte{0xFF},
}

// WorldZakumChannels contains a list of channels where Zakum spawns, 255 = all
var WorldZakumChannels = [WorldCount][]byte{
	[]byte{4, 5, 6},
}

// WorldHorntailChannels contains a list of channels where Horntail spawns, 255 = all
var WorldHorntailChannels = [WorldCount][]byte{
	[]byte{8},
}

// WorldMaxPianusAttempts contains a list of the maximum Pianus attempts allowed on each world, -1 = unlimited
var WorldMaxPianusAttempts = [WorldCount]int16{-1}

// WorldMaxPapAttempts contains a list of the maximum Pianus attempts allowed on each world, -1 = unlimited
var WorldMaxPapAttempts = [WorldCount]int16{2}

// WorldMaxZakumAttempts contains a list of the maximum Pianus attempts allowed on each world, -1 = unlimited
var WorldMaxZakumAttempts = [WorldCount]int16{2}

// WorldMaxHorntailAttempts contains a list of the maximum Pianus attempts allowed on each world, -1 = unlimited
var WorldMaxHorntailAttempts = [WorldCount]int16{-1}
