package gomonkey

// writeIsolationSize must match the NOP padding around START in
// write_darwin_arm64.s. Current Apple Silicon Macs use 16 KiB pages. The
// runtime check in validateWritePageSize fails safely if a future Darwin ARM64
// system reports a larger page instead of allowing an in-page SIGBUS.
const writeIsolationSize = 16 * 1024
