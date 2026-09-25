package gomonkey

// writeIsolationSize must match the NOP padding around START in
// write_darwin_amd64.s. Darwin AMD64 currently uses 4 KiB pages; larger pages
// are rejected before write changes code-page permissions.
const writeIsolationSize = 4 * 1024
