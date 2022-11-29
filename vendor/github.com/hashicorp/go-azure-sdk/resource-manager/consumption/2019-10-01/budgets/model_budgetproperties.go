package budgets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BudgetProperties struct {
	Amount        float64                  `json:"amount"`
	Category      CategoryType             `json:"category"`
	CurrentSpend  *CurrentSpend            `json:"currentSpend"`
	Filter        *BudgetFilter            `json:"filter"`
	ForecastSpend *ForecastSpend           `json:"forecastSpend"`
	Notifications *map[string]Notification `json:"notifications,omitempty"`
	TimeGrain     TimeGrainType            `json:"timeGrain"`
	TimePeriod    BudgetTimePeriod         `json:"timePeriod"`
}
