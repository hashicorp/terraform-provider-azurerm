// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_account_managed_identity

func (r Resource) ModelObject() any {
	return &CognitiveAccountConnectionAccountManagedIdentityModel{}
}

type CognitiveAccountConnectionAccountManagedIdentityModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	Category           string            `tfschema:"category"`
	Metadata           map[string]string `tfschema:"metadata"`
	Target             string            `tfschema:"target"`
}
