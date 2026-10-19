// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_custom_keys

func (r Resource) ModelObject() any {
	return &CognitiveAccountConnectionCustomKeysModel{}
}

type CognitiveAccountConnectionCustomKeysModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	Category           string            `tfschema:"category"`
	CustomKeys         map[string]string `tfschema:"custom_keys"`
	Target             string            `tfschema:"target"`
	Metadata           map[string]string `tfschema:"metadata"`
}
