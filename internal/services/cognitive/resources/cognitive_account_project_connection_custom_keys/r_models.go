// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_custom_keys

type CognitiveAccountProjectConnectionCustomKeysModel struct {
	Name                      string            `tfschema:"name"`
	CognitiveAccountProjectId string            `tfschema:"cognitive_account_project_id"`
	AuthType                  string            `tfschema:"authentication_type"`
	Category                  string            `tfschema:"category"`
	CustomKeys                map[string]string `tfschema:"custom_keys"`
	Target                    string            `tfschema:"target"`
	Metadata                  map[string]string `tfschema:"metadata"`
}
