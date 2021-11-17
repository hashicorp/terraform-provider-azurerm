package alerts

type AlertProperties struct {
	CloseTime                  *string                    `json:"closeTime,omitempty"`
	CostEntityId               *string                    `json:"costEntityId,omitempty"`
	CreationTime               *string                    `json:"creationTime,omitempty"`
	Definition                 *AlertPropertiesDefinition `json:"definition,omitempty"`
	Description                *string                    `json:"description,omitempty"`
	Details                    *AlertPropertiesDetails    `json:"details,omitempty"`
	ModificationTime           *string                    `json:"modificationTime,omitempty"`
	Source                     *AlertSource               `json:"source,omitempty"`
	Status                     *AlertStatus               `json:"status,omitempty"`
	StatusModificationTime     *string                    `json:"statusModificationTime,omitempty"`
	StatusModificationUserName *string                    `json:"statusModificationUserName,omitempty"`
}
