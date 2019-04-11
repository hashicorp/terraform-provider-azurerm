package azure

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest"
)

func CosmosGetIDFromResponse(resp autorest.Response) (string, error) {
	if resp.Response == nil {
		return "", fmt.Errorf("Error: Unable to get Cosmos ID from Response: http response is nil")
	}

	if resp.Response.Request == nil {
		return "", fmt.Errorf("Error: Unable to get Cosmos ID from Response: Request is nil")
	}

	if resp.Response.Request.URL == nil {
		return "", fmt.Errorf("Error: Unable to get Cosmos ID from Response: URL is nil")
	}

	return resp.Response.Request.URL.Path, nil
}

type CosmosAccountResourceID struct {
	ResourceID
	Account string
}

func ParseCosmosAccountResourceID(id string) (*CosmosAccountResourceID, error) {
	subid, err := ParseAzureResourceID(id)
	if err != nil {
		return nil, err
	}

	account, ok := subid.Path["databaseAccounts"]
	if !ok {
		return nil, fmt.Errorf("Error: Unable to parse Cosmos Database Resource ID: databaseAccounts is missing from: %s", id)
	}

	return &CosmosAccountResourceID{
		ResourceID: *subid,
		Account:    account,
	}, nil
}

type CosmosDatabaseResourceID struct {
	CosmosAccountResourceID
	Database string
}

func ParseCosmosDatabaseResourceID(id string) (*CosmosDatabaseResourceID, error) {
	subid, err := ParseCosmosAccountResourceID(id)
	if err != nil {
		return nil, err
	}

	db, ok := subid.Path["databases"]
	if !ok {
		return nil, fmt.Errorf("Error: Unable to parse Cosmos Database Resource ID: databases is missing from: %s", id)
	}

	return &CosmosDatabaseResourceID{
		CosmosAccountResourceID: *subid,
		Database:                db,
	}, nil
}

type CosmosKeyspaceResourceID struct {
	CosmosAccountResourceID
	Keyspace string
}

func ParseCosmosKeyspaceResourceID(id string) (*CosmosKeyspaceResourceID, error) {
	subid, err := ParseCosmosAccountResourceID(id)
	if err != nil {
		return nil, err
	}

	ks, ok := subid.Path["keyspaces"]
	if !ok {
		return nil, fmt.Errorf("Error: Unable to parse Cosmos Keyspace Resource ID: keyspaces is missing from: %s", id)
	}

	return &CosmosKeyspaceResourceID{
		CosmosAccountResourceID: *subid,
		Keyspace:                ks,
	}, nil
}

type CosmosTableResourceID struct {
	CosmosAccountResourceID
	Table string
}

func ParseCosmosTableResourceID(id string) (*CosmosTableResourceID, error) {
	subid, err := ParseCosmosAccountResourceID(id)
	if err != nil {
		return nil, err
	}

	table, ok := subid.Path["keyspaces"]
	if !ok {
		return nil, fmt.Errorf("Error: Unable to parse Cosmos Table Resource ID: keyspaces is missing from: %s", id)
	}

	return &CosmosTableResourceID{
		CosmosAccountResourceID: *subid,
		Table:                   table,
	}, nil
}
