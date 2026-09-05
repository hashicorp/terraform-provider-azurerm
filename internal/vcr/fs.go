// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vcr

// GzipFS is a cassette.FS implementation that reads and writes gzip-compressed files.
// Note: Not currently active as we successfully removed the "noise" of the RP status caching which was causing 3.4MiB
// bloat in the recorded cassettes. When in use the cassettes are not human-readable, so troubleshooting / checking /
// reviewing is inconvenient.
type GzipFS struct{}
