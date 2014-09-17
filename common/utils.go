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
	"crypto/rand"
	"crypto/sha512"
	"fmt"
)

import "github.com/Francesco149/kagami/common/consts"

// HashPassword returns a salted sha-512 hash of the given password
func HashPassword(password, salt string) string {
	hasher := sha512.New()
	saltedpassword := fmt.Sprintf("%sIREALLYLIKELOLIS%s", password, salt)
	hasher.Write([]byte(saltedpassword))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// MakeSalt generates a random string of fixed length that will be used as a password salt
func MakeSalt() string {
	var salt [consts.SaltLength]byte
	rand.Read(salt[:])

	// make it a valid string
	for i := 0; i < consts.SaltLength; i++ {
		salt[i] %= 93 // characters will be between ascii 33 and ascii 126
		salt[i] += 33
	}

	return string(salt[:])
}

// UnixToTempBanTimestamp converts a unix timestamp (in seconds) to a temp ban timestamp
// (number of 100-ns intervals since 1/1/1601)
func UnixToTempBanTimestamp(unixSeconds int64) uint64 {
	// this should be the offset between the unix timestamp and this weird korean timestamp
	const offset = 116444736000000000
	millisecs := uint64(unixSeconds * 1000)
	nano100 := millisecs * 10000 // number of 100-ns intervals
	return nano100 + offset
}

// UnixToTempBanTimestamp converts a unix timestamp (in seconds) to a item timestamp
func UnixToItemTimestamp(unixSeconds int64) uint64 {
	const realYear2000 = 946681229830
	const itemYear2000 = 1085019342
	millisecs := uint64(unixSeconds * 1000)
	time := (millisecs - realYear2000) / 1000 / 60
	// what the fuck
	return uint64(float64(time)*35.762787) - itemYear2000
}

// UnixToQuestTimestamp converts a unix timestamp (in seconds) to a quest timestamp
func UnixToQuestTimestamp(unixSeconds int64) uint64 {
	const questUnixAge = 27111908
	millisecs := uint64(unixSeconds * 1000)
	time := millisecs / 1000 / 60
	// what the fuck
	return uint64(float64(time)*0.1396987) + questUnixAge
}
