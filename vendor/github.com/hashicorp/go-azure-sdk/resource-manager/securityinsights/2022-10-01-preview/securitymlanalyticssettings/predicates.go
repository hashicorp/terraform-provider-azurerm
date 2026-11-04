package securitymlanalyticssettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityMLAnalyticsSettingOperationPredicate struct {
}

func (p SecurityMLAnalyticsSettingOperationPredicate) Matches(input SecurityMLAnalyticsSetting) bool {

	return true
}
