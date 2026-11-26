// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_entra_id

func (r Resource) ModelObject() any {
	return &CognitiveAccountProjectConnectionEntraIDModel{}
}

type CognitiveAccountProjectConnectionEntraIDModel struct {
	Name                      string            `tfschema:"name"`
	CognitiveAccountProjectId string            `tfschema:"cognitive_account_project_id"`
	AuthType                  string            `tfschema:"authentication_type"`
	Category                  string            `tfschema:"category"`
	Target                    string            `tfschema:"target"`
	Metadata                  map[string]string `tfschema:"metadata"`
}
