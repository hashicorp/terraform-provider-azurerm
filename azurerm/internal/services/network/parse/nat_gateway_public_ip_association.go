package parse

import (
	"fmt"
	"strings"
)

type NatGatewayPublicIPAddressAssociationId struct {
	NatGateway        NatGatewayId
	PublicIPAddressID string
}

func NatGatewayPublicIPAddressAssociationID(input string) (*NatGatewayPublicIPAddressAssociationId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected an ID in the format `{natGatewayID}|{publicIPAddressID} but got %q", input)
	}

	natGatewayId, err := NatGatewayID(segments[0])
	if err != nil {
		return nil, fmt.Errorf("parsing NAT Gateway ID %q: %+v", segments[0], err)
	}

	// whilst we need the Resource ID, we may as well validate it
	publicIPAddress := segments[1]
	if _, err := PublicIpAddressID(publicIPAddress); err != nil {
		return nil, fmt.Errorf("parsing Public IP Address ID %q: %+v", publicIPAddress, err)
	}

	return &NatGatewayPublicIPAddressAssociationId{
		NatGateway:        *natGatewayId,
		PublicIPAddressID: publicIPAddress,
	}, nil
}
