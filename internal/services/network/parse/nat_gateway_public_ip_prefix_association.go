package parse

import (
	"fmt"
	"strings"
)

type NatGatewayPublicIPPrefixAssociationId struct {
	NatGateway       NatGatewayId
	PublicIPPrefixID string
}

func NatGatewayPublicIPPrefixAssociationID(input string) (*NatGatewayPublicIPPrefixAssociationId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected an ID in the format `{natGatewayID}|{publicIPPrefixID} but got %q", input)
	}

	natGatewayId, err := NatGatewayID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing NAT Gateway ID %q: %+v", segments[0], err)
	}

	// whilst we need the Resource ID, we may as well validate it
	publicIPPrefix := segments[1]
	if _, err := PublicIpPrefixID(publicIPPrefix); err != nil {
		return nil, fmt.Errorf("parsing Public IP Address ID %q: %+v", publicIPPrefix, err)
	}

	return &NatGatewayPublicIPPrefixAssociationId{
		NatGateway:       *natGatewayId,
		PublicIPPrefixID: publicIPPrefix,
	}, nil
}
