// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package tfjson

const (
	InitOutput LogMessageType = "init_output"

	MessageInitializingStateStoreStart LogMessageType = "initializing_state_store_start"
)

// InitOutputMessage represents messages of type "init_output"
// Note that:
// - This is a subset of the full set of messages emitted by terraform init. Other messages are logged with type "log".
// - The message_code field is a string that identifies the specific type of "init_output" message it is
//   - In an ideal world the "type" field would be set to this value, as that's how it is intended to be used.
type InitOutputMessage struct {
	baseLogMessage
	MessageCode string `json:"message_code"`
}

type InitializingStateStoreStartMessage struct {
	baseLogMessage
}
