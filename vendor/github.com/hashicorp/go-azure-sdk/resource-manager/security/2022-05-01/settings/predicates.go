package settings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SettingOperationPredicate struct {
}

func (p SettingOperationPredicate) Matches(input Setting) bool {

	return true
}
