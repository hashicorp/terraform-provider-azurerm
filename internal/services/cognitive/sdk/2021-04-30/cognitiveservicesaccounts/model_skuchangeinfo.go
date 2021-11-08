package cognitiveservicesaccounts

type SkuChangeInfo struct {
	CountOfDowngrades              *float64 `json:"countOfDowngrades,omitempty"`
	CountOfUpgradesAfterDowngrades *float64 `json:"countOfUpgradesAfterDowngrades,omitempty"`
	LastChangeDate                 *string  `json:"lastChangeDate,omitempty"`
}
