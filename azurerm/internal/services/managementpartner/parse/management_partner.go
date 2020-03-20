package parse

import (
	"fmt"
	"strings"
)

type ManagementPartnerId struct {
	PartnerId string
}

func ManagementPartnerID(input string) (*ManagementPartnerId, error) {
	// /providers/Microsoft.ManagementPartner/partners/5127255
	segments := strings.Split(input, "/providers/Microsoft.ManagementPartner/partners/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Failure determining target resource ID, resource ID in unexpected format: %q", input)
	}

	partnerId := segments[1]
	managementPartnerID := ManagementPartnerId{
		PartnerId: partnerId,
	}
	return &managementPartnerID, nil
}
