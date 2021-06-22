package namespaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
)

type EHNamespaceProperties struct {
	ClusterArmId               *string                      `json:"clusterArmId,omitempty"`
	CreatedAt                  *string                      `json:"createdAt,omitempty"`
	Encryption                 *Encryption                  `json:"encryption,omitempty"`
	IsAutoInflateEnabled       *bool                        `json:"isAutoInflateEnabled,omitempty"`
	KafkaEnabled               *bool                        `json:"kafkaEnabled,omitempty"`
	MaximumThroughputUnits     *int64                       `json:"maximumThroughputUnits,omitempty"`
	MetricId                   *string                      `json:"metricId,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *string                      `json:"provisioningState,omitempty"`
	ServiceBusEndpoint         *string                      `json:"serviceBusEndpoint,omitempty"`
	Status                     *string                      `json:"status,omitempty"`
	UpdatedAt                  *string                      `json:"updatedAt,omitempty"`
	ZoneRedundant              *bool                        `json:"zoneRedundant,omitempty"`
}

func (o EHNamespaceProperties) ListCreatedAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o EHNamespaceProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o EHNamespaceProperties) ListUpdatedAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o EHNamespaceProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
