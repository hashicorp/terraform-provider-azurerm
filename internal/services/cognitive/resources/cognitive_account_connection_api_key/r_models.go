// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_api_key

func (r Resource) ModelObject() any {
	return &CognitiveAccountConnectionApiKeyModel{}
}

type CognitiveAccountConnectionApiKeyModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	ApiKey             string            `tfschema:"api_key"`
	Category           string            `tfschema:"category"`
	Metadata           map[string]string `tfschema:"metadata"`
	Target             string            `tfschema:"target"`
}
