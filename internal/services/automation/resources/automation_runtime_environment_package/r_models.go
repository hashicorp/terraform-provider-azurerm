// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment_package

type AutomationRuntimeEnvironmentPackageModel struct {
	Name                           string `tfschema:"name"`
	AutomationRuntimeEnvironmentId string `tfschema:"automation_runtime_environment_id"`
	ContentUri                     string `tfschema:"content_uri"`
	ContentVersion                 string `tfschema:"content_version"`
	HashAlgorithm                  string `tfschema:"hash_algorithm"`
	HashValue                      string `tfschema:"hash_value"`
	SizeInBytes                    int64  `tfschema:"size_in_bytes"`
	Version                        string `tfschema:"version"`
	Default                        bool   `tfschema:"default"`
}

func (r Resource) ModelObject() any {
	return &AutomationRuntimeEnvironmentPackageModel{}
}
