// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_account_key

func (r Resource) ModelObject() any {
	return &CognitiveAccountProjectConnectionAccountKeyModel{}
}

type CognitiveAccountProjectConnectionAccountKeyModel struct {
	Name                      string            `tfschema:"name"`
	CognitiveAccountProjectId string            `tfschema:"cognitive_account_project_id"`
	AccountKey                string            `tfschema:"account_key"`
	Category                  string            `tfschema:"category"`
	Metadata                  map[string]string `tfschema:"metadata"`
	Target                    string            `tfschema:"target"`
}
