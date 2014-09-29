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

package packets

import "github.com/Francesco149/maplelib"

// ***********************************************************************
// Login Server

// NewEncryptedPacket creates a new packet and appends a placeholder for
// the encrypted header plus the given header to it
func NewEncryptedPacket(header uint16) (p maplelib.Packet) {
	p = maplelib.NewPacket()
	p.Encode4(0x00000000) // placeholder for the encrypted header
	p.Encode2(header)
	return
}

// Ping returns a ping packet
func Ping() (p maplelib.Packet) {
	p = NewEncryptedPacket(OPing)
	return
}

// Pong returns a pong packet
func Pong() (p maplelib.Packet) {
	p = NewEncryptedPacket(IPong)
	return
}

// AuthSuccessRequestPin returns a login success packet that requests pin from the client
func AuthSuccessRequestPin(username string) (p maplelib.Packet) {
	tacos := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // dunno
		0xFF, 0x6A, 0x01, 0x00, // possibly account id but it doesn't seem to matter in v62
		0x00, // player status (set gender, set pin) but I don't give a shit for now
		0x00, // isAdmin: enables client-side gm commands and disables trading
		0x4E} // some kind of gm-related flag

	pizza := []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xDC, 0x3D, 0x0B,
		0x28, 0x64, 0xC5, 0x01, 0x08, 0x00, 0x00, 0x00}

	p = NewEncryptedPacket(OLoginStatus)
	p.Append(tacos)
	p.EncodeString(username)
	p.Append(pizza)
	return
}

// Login failed reasons for LoginFailed()
const (
	LoginIDDeleted          = 3  // ID deleted or blocked
	LoginIncorrectPassword  = 4  // Incorrect password
	LoginNotRegistered      = 5  // Not a registered id
	LoginSystemError        = 6  // System error
	LoginAlreadyLoggedIn    = 7  // Already logged in
	LoginSystemError2       = 8  // System error
	LoginSystemError3       = 9  // System error
	LoginTooManyConnection  = 10 // Cannot process so many connections
	LoginMustBeOver20       = 11 // Only users older than 20 can use this channel
	LoginCannotLogAsMaster  = 13 // Unable to log on as master at this ip
	LoginWrongGateway       = 14 // Wrong gateway or personal info and weird korean button
	LoginTooManyConnection2 = 15 // Processing request with that korean button!
	LoginMustVerifyEmail    = 16 // Please verify your account through email...
	LoginWrongGateway2      = 17 // Wrong gateway or personal info
	LoginMustVerifyEmail2   = 21 // Please verify your account through email...
	LoginShowLicense        = 23 // License agreement
	LoginMapleEuropeNotice  = 25 // Maple Europe notice =[
	LoginTrialVersionNotice = 27 // Some weird full client notice, probably for trial versions
)

/*
   LoginFailed returns a login failed packet
   reason:
   LoginIDDeleted          = 3  // ID deleted or blocked
   LoginIncorrectPassword  = 4  // Incorrect password
   LoginNotRegistered      = 5  // Not a registered id
   LoginSystemError        = 6  // System error
   LoginAlreadyLoggedIn    = 7  // Already logged in
   LoginSystemError2       = 8  // System error
   LoginSystemError3       = 9  // System error
   LoginTooManyConnection  = 10 // Cannot process so many connections
   LoginMustBeOver20       = 11 // Only users older than 20 can use this channel
   LoginCannotLogAsMaster  = 13 // Unable to log on as master at this ip
   LoginWrongGateway       = 14 // Wrong gateway or personal info and weird korean button
   LoginTooManyConnection2 = 15 // Processing request with that korean button!
   LoginMustVerifyEmail    = 16 // Please verify your account through email...
   LoginWrongGateway2      = 17 // Wrong gateway or personal info
   LoginMustVerifyEmail2   = 21 // Please verify your account through email...
   LoginShowLicense        = 23 // License agreement
   LoginMapleEuropeNotice  = 25 // Maple Europe notice =[
   LoginTrialVersionNotice = 27 // Some weird full client notice, probably for trial versions
*/
func LoginFailed(reason int32) (p maplelib.Packet) {
	p = NewEncryptedPacket(OLoginStatus)
	p.Encode4(uint32(reason))
	p.Encode2(0x0000)
	return
}

// Ban reasons for LoginBanned()
const (
	BanDeleted            = 0  // id has been deleted or blocked (used for ip bans, perma bans, chainbans...)
	BanHacking            = 1  // hacking or illegal use of third party programs
	BanMacro              = 2  // using macro/auto keyboard
	BanAd                 = 3  // illicit promotion and advertising
	BanHarassment         = 4  // harassment
	BanProfane            = 5  // using profane language
	BanScam               = 6  // scamming
	BanMisconduct         = 7  // misconduct
	BanIllegalTransaction = 8  // illegal cash transaction
	BanIllegalCharging    = 9  // illegal charging/funding
	BanTemporary          = 10 // temporary request
	BanImpersonatingGM    = 11 // impersonating GM
	BanIllegalPrograms    = 12 // using illegal programs or violating the game policy
	BanMegaphone          = 13 // cursing, scamming or illegal trading via megaphones
	BanNull               = 14 // empty message
)

/*
   LoginBanned returns a login failed packet that tells the user he's temporarily banned.
   If the timestamp is large enough, it will show as a perma ban.
   reason:
   BanDeleted            = 0 // id has been deleted or blocked (used for ip bans, perma bans, chainbans...)
   BanHacking            = 1 // hacking or illegal use of third party programs
   BanMacro              = 2 // using macro/auto keyboard
   BanAd                 = 3 // illicit promotion and advertising
   BanHarassment         = 4 // harassment
   BanProfane            = 5 // using profane language
   BanScam               = 6 // scamming
   BanMisconduct         = 7 // misconduct
   BanIllegalTransaction = 8 // illegal cash transaction
   BanIllegalCharging    = 9 // illegal charging/funding
   BanTemporary          = 10 // temporary request
   BanImpersonatingGM    = 11 // impersonating GM
   BanIllegalPrograms    = 12 // using illegal programs or violating the game policy
   BanMegaphone          = 13 // cursing, scamming or illegal trading via megaphones
   BanNull               = 14 // empty message
*/
func LoginBanned(koreanTimeExpire uint64, reason byte) (p maplelib.Packet) {
	huahuehua := [5]byte{0x00, 0x00, 0x00, 0x00, 0x00}
	p = NewEncryptedPacket(OLoginStatus)
	p.Encode1(0x02)
	p.Append(huahuehua[:])
	p.Encode1(reason)
	p.Encode8(koreanTimeExpire)
	return
}

// Pin operation ids for PinOperation()
const (
	PinOpAccepted    = 0 // PIN was accepted
	PinOpNew         = 1 // Register a new PIN
	PinOpInvalid     = 2 // Invalid pin / Reenter
	PinOpSystemError = 3 // Connection failed due to system error
	PinOpEnter       = 4 // Enter the pin
)

// PinOperation returns a packet that updates the pin operation status of the client
// mode:
// PinOpAccepted    = 0 // PIN was accepted
// PinOpNew         = 1 // Register a new PIN
// PinOpInvalid     = 2 // Invalid pin / Reenter
// PinOpSystemError = 3 // Connection failed due to system error
// PinOpEnter       = 4 // Enter the pin
func PinOperation(mode byte) (p maplelib.Packet) {
	p = NewEncryptedPacket(OPinOperation)
	p.Encode1(mode)
	return
}

// PinAccepted returns a packet that notifies the client that the pin has been accepted
func PinAccepted() maplelib.Packet {
	return PinOperation(PinOpAccepted)
}

// RequestPinAfterFailure returns a packet that notifies the client that the pin
// is wrong and must be re-entered
func RequestPinAfterFailure() maplelib.Packet {
	return PinOperation(PinOpInvalid)
}

// RequestPin returns a packet that tells the client to request a pin from the user
func RequestPin() maplelib.Packet {
	return PinOperation(PinOpEnter)
}

// WorldListEnd returns a packet that indicates the end of a world list
func WorldListEnd() (p maplelib.Packet) {
	p = NewEncryptedPacket(OServerList)
	p.Encode1(0xFF)
	return
}

// Possible values for ServerStatus()
const (
	ServerNormal = 0 // Normal load
	ServerHigh   = 1 // Highly populated
	ServerFull   = 2 // Full
)

// ServerStatus returns a packet that tells the client how full the world is
// possible values for status:
// ServerNormal = 0 // Normal load
// ServerHigh = 1 // Highly populated
// ServerFull = 2 // Full
func ServerStatus(status uint16) (p maplelib.Packet) {
	p = NewEncryptedPacket(OServerStatus)
	p.Encode2(status)
	return
}

// SendAllCharsBegin returns a packet that sends the beginning of a character list packet
// unk must be charcount + (3 - charcount % 3)
func SendAllCharsBegin(worldcount, unk uint32) (p maplelib.Packet) {
	p = NewEncryptedPacket(OAllCharlist)
	p.Encode1(0x01)
	p.Encode4(worldcount)
	p.Encode4(unk)
	return
}

// RelogResponse rturns a packet that accepts a relog request
func RelogResponse() (p maplelib.Packet) {
	p = NewEncryptedPacket(ORelogResponse)
	p.Encode1(0x01)
	return
}

// ConnectIp returns a server transfer packet
func ConnectIp(ip []byte, port int16, charId int32) (p maplelib.Packet) {
	huehuahue := make([]byte, 5)
	p = NewEncryptedPacket(OServerIP)
	p.Encode2(0x0000)
	p.Append(ip)
	p.Encode2(uint16(port))
	p.Encode4(uint32(charId))
	p.Append(huehuahue)
	return
}

// CharNameResponse returns a char name check response packet
func CharNameResponse(charName string, used bool) (p maplelib.Packet) {
	p = NewEncryptedPacket(OCharNameResponse)
	p.EncodeString(charName)
	if used {
		p.Encode1(0x01)
	} else {
		p.Encode1(0x00)
	}
	return
}

// Possible statuses for DeleteCharResponse
const (
	DeleteOk          = 0x00 // ok
	DeleteFail        = 0x01 // failed to delete character
	DeleteInvalidCode = 0x12 // invalid birthday
)

// DeleteCharResponse returns a char delete response packet
// state:
// DeleteOk = 0 // ok
// DeleteFail = 1 // failed to delete character
// DeleteInvalidCode = 0x12 // invalid birthday
func DeleteCharResponse(id int32, state byte) (p maplelib.Packet) {
	p = NewEncryptedPacket(ODeleteCharResponse)
	p.Encode4(uint32(id))
	p.Encode1(state)
	return
}

// SetGenderDone returns a packet response for a set gender request
func SetGenderDone(gender byte) (p maplelib.Packet) {
	p = NewEncryptedPacket(OGenderDone)
	p.Encode1(gender)
	p.Encode1(0x01)
	return
}

// PinAssigned returns a packet that tells the client that the pin has successfully been assigned
func PinAssigned() (p maplelib.Packet) {
	p = NewEncryptedPacket(OPinAssigned)
	p.Encode1(0x01)
	return
}

// ***********************************************************************
// Channel Server

// Possible values for msgtype in ServerMessage()
const (
	ServerMessageNotice          = 0 // [Notice]
	ServerMessagePopup           = 1 // Popup
	ServerMessageMega            = 2 // Megaphone
	ServerMessageSmega           = 3 // Super Megaphone
	ServerMessageScrollingHeader = 4 // Scrolling header
	ServerMessagePinkText        = 5 // Pink text
	ServerMessageLightBlueText   = 6 // Light blue text
)

/*
	ServerMessage returns a server message packet

	Possible values for msgtype:
	ServerMessageNotice = 0 // [Notice]
	ServerMessagePopup = 1 // Popup
	ServerMessageMega = 2 // Megaphone
	ServerMessageSmega = 3 // Super Megaphone
	ServerMessageScrollingHeader = 4 // Scrolling header
	ServerMessagePinkText = 5 // Pink text
	ServerMessageLightBlueText = 6 // Light blue text
*/
func ServerMessage(msgtype, channel int8, message string,
	isScrollingHeader bool, megaEar bool) (p maplelib.Packet) {

	p = NewEncryptedPacket(OServerMessage)
	p.Encode1s(msgtype)

	if isScrollingHeader {
		p.Encode1(0x01)
	}

	p.EncodeString(message)

	if msgtype == ServerMessageSmega {
		p.Encode1s(channel)
		var tmp byte
		if megaEar {
			tmp = 1
		} else {
			tmp = 0
		}
		p.Encode1(tmp)
	}

	return
}

// ScrollingHeader returns a packet that
// updates the yellow scrolling header for the player
func ScrollingHeader(msg string) maplelib.Packet {
	return ServerMessage(ServerMessageScrollingHeader, 0, msg, true, false)
}

// ChangeChannel returns a packet that tells the client the ip of the new channel
func ChangeChannel(ip []byte, port int16) (p maplelib.Packet) {
	p = NewEncryptedPacket(OChangeChannel)
	p.Encode1(0x01)
	p.Append(ip)
	p.Encode2s(port)
	return
}
