package charprops

//go:generate go run ./generate.go

var zeroProperties CharProperties

// LookupFirstChar identifies the grapheme cluster segmentation properties and
// the byte length of the first UTF-8 sequence in the given buffer.
//
// If there are not enough bytes in the buffer to recognize a complete UTF-8
// sequence then this returns a length of zero and a meaningless
// [CharProperties] value.
//
// If the buffer begins with something that cannot possibly be extended to a
// valid UTF-8 sequence then the returned properties are [Error] and length
// as 1 to represent treating just the next byte as an error.
func LookupFirstChar(p []byte) (props CharProperties, length int) {
	if len(p) == 0 {
		return zeroProperties, 0
	}

	first := p[0]
	if first < 128 {
		length = 1
	} else if (first & 0b11100000) == 0b11000000 {
		length = 2
	} else if (first & 0b11110000) == 0b11100000 {
		length = 3
	} else if (first & 0b11111000) == 0b11110000 {
		length = 4
	} else {
		// Invalid initial byte.
		return Error, 1
	}

	if len(p) < length {
		// Buffer begins with incomplete UTF-8 sequence.
		return zeroProperties, 0
	}
	// All of the subsequent bytes of the sequence we're decoding (if any)
	// must be UTF-8 continuation bytes.
	for i := 1; i < length; i++ {
		if (p[i] & 0b11000000) != 0b10000000 {
			return Error, 1
		}
	}

	switch length {
	case 1:
		return lookupProps[first], length
	case 2:
		blockIdx := int(lookupIndices[first&0b111111])
		return lookupProps[(blockIdx<<6)+int(p[1]&0b111111)], length
	case 3:
		blockIdx := int(lookupIndices[first&0b111111])
		blockIdx = int(lookupIndices[(blockIdx<<6)+int(p[1]&0b111111)])
		return lookupProps[(blockIdx<<6)+int(p[2]&0b111111)], length
	case 4:
		blockIdx := int(lookupIndices[first&0b111111])
		blockIdx = int(lookupIndices[(blockIdx<<6)+int(p[1]&0b111111)])
		blockIdx = int(lookupIndices[(blockIdx<<6)+int(p[2]&0b111111)])
		return lookupProps[(blockIdx<<6)+int(p[3]&0b111111)], length
	default:
		panic("unreachable") // (because we should've caught this case above)
	}
}
