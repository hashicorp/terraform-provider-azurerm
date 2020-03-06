package parse

import (
	"fmt"
	"strings"
)

type ManagementPartnerId struct {
	PartnerId string
}

func ParseManagementPartnerID(input string) (*ManagementPartnerId, error) {
	// /providers/Microsoft.ManagementPartner/partners/5127255
	segments := strings.Split(input, "/")
	if len(segments) != 5 {
		return nil, fmt.Errorf("Expected there to be 5 segments but got %d", len(segments))
	}

	id := ManagementPartnerId{
		PartnerId: segments[4],
	}
	return &id, nil
}
