// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_account_key

func (r Resource) ModelObject() any {
	return &CognitiveAccountConnectionAccountKeyModel{}
}

type CognitiveAccountConnectionAccountKeyModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	AccountKey         string            `tfschema:"account_key"`
	Category           string            `tfschema:"category"`
	Metadata           map[string]string `tfschema:"metadata"`
	Target             string            `tfschema:"target"`
}
