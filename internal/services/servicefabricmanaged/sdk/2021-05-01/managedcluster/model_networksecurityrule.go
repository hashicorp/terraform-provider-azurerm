package managedcluster

type NetworkSecurityRule struct {
	Access                     Access      `json:"access"`
	Description                *string     `json:"description,omitempty"`
	DestinationAddressPrefixes *[]string   `json:"destinationAddressPrefixes,omitempty"`
	DestinationPortRanges      *[]string   `json:"destinationPortRanges,omitempty"`
	Direction                  Direction   `json:"direction"`
	Name                       string      `json:"name"`
	Priority                   int64       `json:"priority"`
	Protocol                   NsgProtocol `json:"protocol"`
	SourceAddressPrefixes      *[]string   `json:"sourceAddressPrefixes,omitempty"`
	SourcePortRanges           *[]string   `json:"sourcePortRanges,omitempty"`
}
