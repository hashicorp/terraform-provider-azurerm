package machine

import (
	"github.com/apparentlymart/go-textseg/v17/textseg/internal/charprops"
)

// The state machine implemented here is a Go port of a Rust implementation
// in the "grapheme_machine" package on crates.io. It's more general than
// the main textseg package currently needs but maybe a future version will
// expose more of this in the public API if a compelling use-case emerges for
// that.

// State represents different states we can transition through while detecting
// grapheme cluster boundaries.
//
// A State value essentially summarizes a set of category transitions that
// happened before the current one, so that we can detect arbitrary-long
// grapheme clusters using only finite storage.
//
// (In order to actually _use_ a detected grapheme cluster after its bounds
// have been found would require the caller to have buffered everything that
// appeared since the last boundary, but that's the caller's problem.)
type State int

const (
	// StateBase is the initial state at the beginning of the text or when the
	// following should be treated as if it were at the beginning of the text.
	StateBase State = iota

	// When the previous character was [charprops.GCBRegionalIndicator] but
	// its predecessor was not, and therefore if the next character
	// is also [charprops.GCBRegionalIndicator] the two together represent
	// an emoji flag under rules GB12 and GB13.
	StateAwaitEmojiFlag

	// When consecutive scalar values before the previous matched
	// [charprops.GCBExtendedPictographic] followed by any number of
	// [charprops.GCBExtend], and therefore if the next character is
	// [charprops.GCBZWJ] we should transition to [StateGB11AfterZWJ].
	StateGB11BeforeZWJ

	// The previous scalar value of category [charprops.GCBZWJ] arrived while
	// in [StateGB11BeforeZWJ], and therefore rule GB11 is active.
	StateGB11AfterZWJ

	// We encountered [charprops.InCBConsonant] followed by zero or more
	// [charprops.InCBExtend] with no [charprops.InCBLinker] in between.
	StateGB9cConsonant

	// We encountered [charprops.InCBConsonant] followed by at least one
	// [charprops.InCBLinker] and zero or more [charprops.InCBExtend] in any
	// order.
	StateGB9cLinker
)

// Begin initializes the state machine by providing the properties for the
// first character. Call [State.Transition] on the resulting state once
// the second character is known, passing the same properties as "first" in
// its "prev" argument.
func Begin(first charprops.CharProperties) State {
	return StateBase.nextState(first)
}

// Given the previous category and the next category, Transition returns whether
// there is a grapheme cluster boundary between two characters of those
// categories in the receiving state, and which state should be used for the
// next transition.
//
// Correct use requires that the "prev" of one call equal the "next" of the
// call that selected the recieving state. If that is not upheld then the
// results are unspecified.
func (s State) Transition(prev, next charprops.CharProperties) (bool, State) {
	nextState := s.nextState(next)

	// The following was slightly more readable with Rust pattern matching
	// syntax, but this is a direct translation of the same logic into Go if
	// statements.

	prevC := prev.GCBProperty()
	nextC := next.GCBProperty()

	// GB1 and GB2 aren't covered here because the caller is expected to
	// use this function only at the boundary between two scalar values, and
	// not at the beginning or end of the input. There is always a grapheme
	// cluster boundary at the beginning and end of the input.

	// GB3: Do not break between CR and LF...
	if prevC == charprops.GCBCR && nextC == charprops.GCBLF {
		return false, nextState
	}
	// GB4 and GB5: ...but otherwise, break before and after any control characters
	if prevC.IsAnyControl() || nextC.IsAnyControl() {
		return true, nextState
	}
	// GB6: Do not break Hangul syllable or other conjoining sequences.
	if prevC == charprops.GCBL && (nextC == charprops.GCBL || nextC == charprops.GCBV || nextC == charprops.GCBLV || nextC == charprops.GCBLVT) {
		return false, nextState
	}
	// GB7: Do not break Hangul syllable or other conjoining sequences.
	if (prevC == charprops.GCBLV || prevC == charprops.GCBV) && (nextC == charprops.GCBV || nextC == charprops.GCBT) {
		return false, nextState
	}
	// GB8: Do not break Hangul syllable or other conjoining sequences.
	if (prevC == charprops.GCBLVT || prevC == charprops.GCBT) && nextC == charprops.GCBT {
		return false, nextState
	}
	// GB9: Do not break before extending characters or ZWJ
	if nextC == charprops.GCBExtend || nextC == charprops.GCBZWJ {
		return false, nextState
	}
	// GB9a: Do not break before spacing marks...
	if nextC == charprops.GCBSpacingMark {
		return false, nextState
	}
	// GB9a: ...or after Prepend characters
	if prevC == charprops.GCBPrepend {
		return false, nextState
	}
	// GB9c: Do not break within certain combinations with Indic_Conjunct_Break (InCB)=Linker
	if s.gb9cActive() {
		prevI := prev.InCBProperty()
		nextI := next.InCBProperty()
		if (prevI == charprops.InCBLinker || prevI == charprops.InCBExtend) && nextI == charprops.InCBConsonant {
			return false, nextState
		}
	}
	// (GB10 was from an earlier version of the specification but is no longer used)
	// GB11: Do not break within emoji modifier sequences or emoji zwj sequences.
	if s.gb11Active() {
		if prevC == charprops.GCBZWJ && nextC == charprops.GCBExtendedPictographic {
			return false, nextState
		}
	}
	// GB12 and GB13: Do not break within emoji flag sequences.
	if s.gb13Active() {
		if prevC == charprops.GCBRegionalIndicator && nextC == charprops.GCBRegionalIndicator {
			return false, nextState
		}
	}
	// GB999: Break in all other locations
	return true, nextState
}

func (s State) nextState(next charprops.CharProperties) State {
	// Two of the multi-character prefixes can begin regardless of what
	// preceeds them. These don't need to be covered by the state-specific
	// branches that follow.
	if next.GCBProperty() == charprops.GCBExtendedPictographic {
		return StateGB11BeforeZWJ
	}
	if next.InCBProperty() == charprops.InCBConsonant {
		return StateGB9cConsonant
	}
	gcbProp := next.GCBProperty()
	incbProp := next.InCBProperty()
	switch s {
	case StateBase:
		switch gcbProp {
		case charprops.GCBRegionalIndicator:
			return StateAwaitEmojiFlag
		default:
			return StateBase
		}
	case StateAwaitEmojiFlag:
		return StateBase
	case StateGB11BeforeZWJ:
		switch gcbProp {
		case charprops.GCBZWJ:
			return StateGB11AfterZWJ
		case charprops.GCBExtend:
			return StateGB11BeforeZWJ
		default:
			return StateBase
		}
	case StateGB11AfterZWJ:
		return StateBase
	case StateGB9cConsonant:
		switch incbProp {
		case charprops.InCBLinker:
			return StateGB9cLinker
		case charprops.InCBExtend:
			return StateGB9cConsonant
		default:
			return StateBase
		}
	case StateGB9cLinker:
		switch incbProp {
		case charprops.InCBLinker, charprops.InCBExtend:
			return StateGB9cLinker
		default:
			return StateBase
		}
	default:
		panic("invalid state")
	}
}

func (s State) gb9cActive() bool {
	// GB9c is active only in [StateGB9cLinker].
	return s == StateGB9cLinker
}

func (s State) gb11Active() bool {
	// GB11 is active only in [StateGB11AfterZWJ].
	return s == StateGB11AfterZWJ
}

func (s State) gb13Active() bool {
	// GB12/GB13 is active only in [StateAwaitEmojiFlag].
	return s == StateAwaitEmojiFlag
}
