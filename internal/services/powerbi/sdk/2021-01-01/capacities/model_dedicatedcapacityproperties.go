package capacities

type DedicatedCapacityProperties struct {
	Administration    *DedicatedCapacityAdministrators `json:"administration,omitempty"`
	Mode              *Mode                            `json:"mode,omitempty"`
	ProvisioningState *CapacityProvisioningState       `json:"provisioningState,omitempty"`
	State             *State                           `json:"state,omitempty"`
}
