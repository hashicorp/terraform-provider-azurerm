// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_entra_id

func (r Resource) ModelObject() any {
	return &CognitiveAccountConnectionEntraIDModel{}
}

type CognitiveAccountConnectionEntraIDModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	Category           string            `tfschema:"category"`
	Target             string            `tfschema:"target"`
	Metadata           map[string]string `tfschema:"metadata"`
}
