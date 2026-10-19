// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_python3_package

type Python3PackageModel struct {
	Name                  string            `tfschema:"name"`
	ResourceGroupName     string            `tfschema:"resource_group_name"`
	AutomationAccountName string            `tfschema:"automation_account_name"`
	ContentUri            string            `tfschema:"content_uri"`
	ContentVersion        string            `tfschema:"content_version"`
	HashAlgorithm         string            `tfschema:"hash_algorithm"`
	HashValue             string            `tfschema:"hash_value"`
	Tags                  map[string]string `tfschema:"tags"`
}

func (m Resource) ModelObject() any {
	return &Python3PackageModel{}
}
