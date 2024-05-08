package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlinkJobProfile struct {
	Args            *string     `json:"args,omitempty"`
	EntryClass      *string     `json:"entryClass,omitempty"`
	JarName         string      `json:"jarName"`
	JobJarDirectory string      `json:"jobJarDirectory"`
	SavePointName   *string     `json:"savePointName,omitempty"`
	UpgradeMode     UpgradeMode `json:"upgradeMode"`
}
