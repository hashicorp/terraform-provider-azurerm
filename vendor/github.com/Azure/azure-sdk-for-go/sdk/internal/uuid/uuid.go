// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package uuid

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// The UUID reserved variants.
const (
	reservedNCS       byte = 0x80
	reservedRFC4122   byte = 0x40
	reservedMicrosoft byte = 0x20
	reservedFuture    byte = 0x00
)

func init() {
	rand.Seed(time.Now().Unix())
}

// A UUID representation compliant with specification in RFC 4122 document.
type UUID [16]byte

// New returns a new uuid using RFC 4122 algorithm.
func New() UUID {
	u := UUID{}
	// Set all bits to randomly (or pseudo-randomly) chosen values.
	// math/rand.Read() is no-fail so we omit any error checking.
	// NOTE: this takes a process-wide lock
	rand.Read(u[:])
	u[8] = (u[8] | reservedRFC4122) & 0x7F // u.setVariant(ReservedRFC4122)

	var version byte = 4
	u[6] = (u[6] & 0xF) | (version << 4) // u.setVersion(4)
	return u
}

// String returns an unparsed version of the generated UUID sequence.
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// Parse parses a string formatted as "003020100-0504-0706-0809-0a0b0c0d0e0f"
// or "{03020100-0504-0706-0809-0a0b0c0d0e0f}" into a UUID.
func Parse(uuidStr string) UUID {
	char := func(hexString string) byte {
		i, _ := strconv.ParseUint(hexString, 16, 8)
		return byte(i)
	}
	if uuidStr[0] == '{' {
		uuidStr = uuidStr[1:] // Skip over the '{'
	}
	// 03020100 - 05 04 - 07 06 - 08 09 - 0a 0b 0c 0d 0e 0f
	//             1 11 1 11 11 1 12 22 2 22 22 22 33 33 33
	// 01234567 8 90 12 3 45 67 8 90 12 3 45 67 89 01 23 45
	uuidVal := UUID{
		char(uuidStr[0:2]),
		char(uuidStr[2:4]),
		char(uuidStr[4:6]),
		char(uuidStr[6:8]),

		char(uuidStr[9:11]),
		char(uuidStr[11:13]),

		char(uuidStr[14:16]),
		char(uuidStr[16:18]),

		char(uuidStr[19:21]),
		char(uuidStr[21:23]),

		char(uuidStr[24:26]),
		char(uuidStr[26:28]),
		char(uuidStr[28:30]),
		char(uuidStr[30:32]),
		char(uuidStr[32:34]),
		char(uuidStr[34:36]),
	}
	return uuidVal
}

func (u UUID) bytes() []byte {
	return u[:]
}
