package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthentication"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerlistener"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflow"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowendpoint"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowprofile"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BrokerAuthenticationClient *brokerauthentication.BrokerAuthenticationClient
	BrokerAuthorizationClient  *brokerauthorization.BrokerAuthorizationClient
	BrokerClient               *broker.BrokerClient
	BrokerListenerClient       *brokerlistener.BrokerListenerClient
	DataflowClient             *dataflow.DataflowClient
	DataflowEndpointClient     *dataflowendpoint.DataflowEndpointClient
	DataflowProfileClient      *dataflowprofile.DataflowProfileClient
	InstanceClient             *instance.InstanceClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	brokerAuthenticationClient := brokerauthentication.NewBrokerAuthenticationClientWithBaseURI(o.ResourceManagerEndpoint)
	brokerAuthenticationClient.Client.Authorizer = o.ResourceManagerAuthorizer

	brokerAuthorizationClient := brokerauthorization.NewBrokerAuthorizationClientWithBaseURI(o.ResourceManagerEndpoint)
	brokerAuthorizationClient.Client.Authorizer = o.ResourceManagerAuthorizer

	brokerClient := broker.NewBrokerClientWithBaseURI(o.ResourceManagerEndpoint)
	brokerClient.Client.Authorizer = o.ResourceManagerAuthorizer

	brokerListenerClient := brokerlistener.NewBrokerListenerClientWithBaseURI(o.ResourceManagerEndpoint)
	brokerListenerClient.Client.Authorizer = o.ResourceManagerAuthorizer

	dataflowClient := dataflow.NewDataflowClientWithBaseURI(o.ResourceManagerEndpoint)
	dataflowClient.Client.Authorizer = o.ResourceManagerAuthorizer

	dataflowEndpointClient := dataflowendpoint.NewDataflowEndpointClientWithBaseURI(o.ResourceManagerEndpoint)
	dataflowEndpointClient.Client.Authorizer = o.ResourceManagerAuthorizer

	dataflowProfileClient := dataflowprofile.NewDataflowProfileClientWithBaseURI(o.ResourceManagerEndpoint)
	dataflowProfileClient.Client.Authorizer = o.ResourceManagerAuthorizer

	instanceClient := instance.NewInstanceClientWithBaseURI(o.ResourceManagerEndpoint)
	instanceClient.Client.Authorizer = o.ResourceManagerAuthorizer

	return &Client{
		BrokerAuthenticationClient: brokerAuthenticationClient,
		BrokerAuthorizationClient:  brokerAuthorizationClient,
		BrokerClient:               brokerClient,
		BrokerListenerClient:       brokerListenerClient,
		DataflowClient:             dataflowClient,
		DataflowEndpointClient:     dataflowEndpointClient,
		DataflowProfileClient:      dataflowProfileClient,
		InstanceClient:             instanceClient,
	}, nil
}
