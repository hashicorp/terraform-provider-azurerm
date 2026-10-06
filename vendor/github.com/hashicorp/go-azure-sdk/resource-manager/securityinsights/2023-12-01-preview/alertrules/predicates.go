package alertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleOperationPredicate struct {
}

func (p AlertRuleOperationPredicate) Matches(input AlertRule) bool {

	return true
}
