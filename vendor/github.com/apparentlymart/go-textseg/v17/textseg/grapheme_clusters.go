package textseg

import (
	"github.com/apparentlymart/go-textseg/v17/textseg/internal/charprops"
	"github.com/apparentlymart/go-textseg/v17/textseg/internal/machine"
)

// ScanGraphemeClusters is a split function for [bufio.Scanner] that splits
// on grapheme cluster boundaries.
//
// Note that while this function does some minimal UTF-8 validation as part
// of its work, it never actually returns an error and instead just reports
// invalid bytes as individual grapheme clusters. Also in particular it does
// not currently check for overlong encodings or out-of-range Unicode scalar
// values, and so callers should perform their own validation if needed. The
// handling of overlong and out-of-range input may change in future releases.
func ScanGraphemeClusters(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	properties, count := charprops.LookupFirstChar(data)
	if properties == charprops.Error {
		return 1, data[:1], nil
	}
	if count == 0 {
		if len(data) > count && atEOF {
			// If we're already at EOF then this is invalid input, which
			// we treat as each invalid byte being a separate grapheme cluster.
			return 1, data[:1], nil
		}
		return 0, nil, nil
	}
	remain := data[count:]

	state := machine.Begin(properties)
	prev := properties
	for {
		if len(remain) == 0 {
			if atEOF {
				// If we're at the end of the file then whatever we've
				// accumulated so far is a grapheme cluster.
				return count, data[:count], nil
			}
			// If we're not at the end of the file then we'll need more
			// bytes before we can decide if we've reached the end of a
			// grapheme cluster.
			return 0, nil, nil
		}

		next, moreCount := charprops.LookupFirstChar(remain)
		if next == charprops.Error {
			// If the next sequence is invalid then we'll just return here
			// and let the next call deal with that.
			return count, data[:count], nil
		}
		if moreCount == 0 {
			// More bytes required to complete the next UTF-8 sequence.
			if len(remain) != 0 && atEOF {
				// If we're already at EOF then the next UTF-8 sequence is
				// invalid, so we'll report what we already accumulated and
				// then let a subsequent call deal with the invalid byte.
				return count, data[:count], nil
			}
			return 0, nil, nil
		}
		remain = remain[moreCount:]

		split, nextState := state.Transition(prev, next)
		if split {
			// We've found the next split point, so we'll just return what
			// we have and then let the next call pick up from here.
			return count, data[:count], nil
		}
		count += moreCount
		state = nextState
		prev = next
	}
}
