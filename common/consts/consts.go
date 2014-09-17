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

const MySQLUser = "kagami"         // MySQLUser is the MySQL username
const MySQLPassword = "testing"    // MySQLPassword is the MySQL password
const MySQLHost = "127.0.0.1:3306" // MySQLHost contains the ip:port of the MySQL database
const MySQLDB = "my_kagami"        // MySQLDB contains the name of the used MySQL database

const Loginport = 8484 // Loginport is the port the Login Server will listen on

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
