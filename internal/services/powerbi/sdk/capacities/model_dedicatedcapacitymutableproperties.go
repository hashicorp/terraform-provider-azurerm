package capacities

type DedicatedCapacityMutableProperties struct {
	Administration *DedicatedCapacityAdministrators `json:"administration,omitempty"`
	Mode           *Mode                            `json:"mode,omitempty"`
}
