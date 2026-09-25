//go:build darwin && arm64
// +build darwin,arm64

#include "textflag.h"

#define NOP8 WORD $0x1f2003d5; WORD $0x1f2003d5;
#define NOP24 NOP8; NOP8; NOP8;
#define NOP192 NOP24; NOP24; NOP24; NOP24; NOP24; NOP24; NOP24; NOP24;
#define NOP1536 NOP192; NOP192; NOP192; NOP192; NOP192; NOP192; NOP192; NOP192;

// privateMethodTrampolineBase exposes the exact start of the executable slot
// area. Returning the address from assembly avoids an ABI wrapper changing the
// address observed by Go.
TEXT ·privateMethodTrampolineBase(SB),NOSPLIT,$0-8
	MOVD $privateMethodTrampolineSpace<>(SB), R0
	MOVD R0, ret+0(FP)
	RET

// 128 fixed-size slots. Each slot holds one 24-byte buildJmpDirective result.
TEXT privateMethodTrampolineSpace<>(SB),NOSPLIT|NOFRAME,$0-0
	NOP1536
	NOP1536
	RET
