// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package features

import (
	"os"
	"strings"
)

// SixPointOh returns whether this provider is running in 6.0 mode
// that is to say - the final 6.0 release
//
// This exists to allow breaking changes to be piped through the provider
// during the development of 5.x until 6.0 is ready.
// The environment variable `ARM_SIXPOINTZERO_BETA` has been added
// to facilitate testing. But it should be noted that
// `ARM_SIXPOINTZERO_BETA` is ** NOT READY FOR PUBLIC USE ** and
// ** SHOULD NOT BE SET IN PRODUCTION ENVIRONMENTS **
// Setting `ARM_SIXPOINTZERO_BETA` will cause irreversible changes
// to your state.
func SixPointOh() bool {
	return strings.EqualFold(os.Getenv("ARM_SIXPOINTZERO_BETA"), "true")
}
