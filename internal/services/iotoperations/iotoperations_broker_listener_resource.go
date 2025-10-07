package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerlistener"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BrokerListenerResource struct{}

var _ sdk.ResourceWithUpdate = BrokerListenerResource{}

type BrokerListenerModel struct {
	Name              string                      `tfschema:"name"`
	ResourceGroupName string                      `tfschema:"resource_group_name"`
	InstanceName      string                      `tfschema:"instance_name"`
	BrokerName        string                      `tfschema:"broker_name"`
	ServiceName       *string                     `tfschema:"service_name"`
	ServiceType       *string                     `tfschema:"service_type"`
	Ports             []BrokerListenerPortModel   `tfschema:"ports"`
	ExtendedLocation  *ExtendedLocationModel      `tfschema:"extended_location"`
	Tags              map[string]string           `tfschema:"tags"`
	ProvisioningState *string                     `tfschema:"provisioning_state"`
}

type BrokerListenerPortModel struct {
	Port              int                           `tfschema:"port"`
	NodePort          *int                          `tfschema:"node_port"`
	Protocol          *string                       `tfschema:"protocol"`
	AuthenticationRef *string                       `tfschema:"authentication_ref"`
	AuthorizationRef  *string                       `tfschema:"authorization_ref"`
	Tls               *BrokerListenerTlsModel       `tfschema:"tls"`
}

type BrokerListenerTlsModel struct {
	Mode                        string                                         `tfschema:"mode"`
	CertManagerCertificateSpec  *BrokerListenerCertManagerCertificateSpecModel `tfschema:"cert_manager_certificate_spec"`
	Manual                      *BrokerListenerManualModel                     `tfschema:"manual"`
}

type BrokerListenerCertManagerCertificateSpecModel struct {
	Duration    *string                                    `tfschema:"duration"`
	SecretName  *string                                    `tfschema:"secret_name"`
	RenewBefore *string                                    `tfschema:"renew_before"`
	IssuerRef   *BrokerListenerIssuerRefModel              `tfschema:"issuer_ref"`
	PrivateKey  *BrokerListenerPrivateKeyModel             `tfschema:"private_key"`
	San         *BrokerListenerSanModel                    `tfschema:"san"`
}

type BrokerListenerIssuerRefModel struct {
	Group *string `tfschema:"group"`
	Kind  *string `tfschema:"kind"`
	Name  *string `tfschema:"name"`
}

type BrokerListenerPrivateKeyModel struct {
	Algorithm      *string `tfschema:"algorithm"`
	RotationPolicy *string `tfschema:"rotation_policy"`
}

type BrokerListenerSanModel struct {
	Dns *[]string `tfschema:"dns"`
	Ip  *[]string `tfschema:"ip"`
}

type BrokerListenerManualModel struct {
	SecretRef *string `tfschema:"secret_ref"`
}

func (r BrokerListenerResource) ModelObject() interface{} {
	return &BrokerListenerModel{}
}

func (r BrokerListenerResource) ResourceType() string {
	return "azurerm_iotoperations_broker_listener"
}

func (r BrokerListenerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return brokerlistener.ValidateListenerID
}

func (r BrokerListenerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 90),
		},
		"instance_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"broker_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 63),
				validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
			),
		},
		"service_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 63),
		},
		"service_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"LoadBalancer",
				"NodePort",
				"ClusterIP",
			}, false),
		},
		"ports": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"port": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 65535),
					},
					"node_port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(30000, 32767),
					},
					"protocol": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"MQTT",
							"WebSockets",
						}, false),
					},
					"authentication_ref": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"authorization_ref": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 253),
					},
					"tls": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"mode": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Disabled",
										"Manual",
										"Automatic",
									}, false),
								},
								"cert_manager_certificate_spec": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"duration": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 50),
											},
											"secret_name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 253),
											},
											"renew_before": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 50),
											},
											"issuer_ref": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"group": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringLenBetween(1, 253),
														},
														"kind": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringLenBetween(1, 63),
														},
														"name": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ValidateFunc: validation.StringLenBetween(1, 253),
														},
													},
												},
											},
											"private_key": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"algorithm": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															ValidateFunc: validation.StringInSlice([]string{
																"RSA",
																"ECDSA",
																"Ed25519",
															}, false),
														},
														"rotation_policy": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															ValidateFunc: validation.StringInSlice([]string{
																"Always",
																"Never",
															}, false),
														},
													},
												},
											},
											"san": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"dns": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															Elem: &pluginsdk.Schema{
																Type:         pluginsdk.TypeString,
																ValidateFunc: validation.StringLenBetween(1, 253),
															},
														},
														"ip": {
															Type:     pluginsdk.TypeList,
															Optional: true,
															Elem: &pluginsdk.Schema{
																Type:         pluginsdk.TypeString,
																ValidateFunc: validation.IsIPAddress,
															},
														},
													},
												},
											},
										},
									},
								},
								"manual": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"secret_ref": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringLenBetween(1, 253),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"extended_location": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"CustomLocation",
						}, false),
					},
				},
			},
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r BrokerListenerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"provisioning_state": {
			Type:     pluginsdk.TypeString,
			// NOTE: O+C Azure automatically assigns provisioning state during resource lifecycle
			Computed: true,
		},
	}
}

func (r BrokerListenerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerListenerClient

			var model BrokerListenerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := brokerlistener.NewListenerID(subscriptionId, model.ResourceGroupName, model.InstanceName, model.BrokerName, model.Name)

			// Build payload
			payload := brokerlistener.BrokerListenerResource{
				Properties: expandBrokerListenerProperties(model),
			}

			if model.ExtendedLocation != nil {
				payload.ExtendedLocation = expandExtendedLocation(model.ExtendedLocation)
			}

			if len(model.Tags) > 0 {
				payload.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BrokerListenerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerListenerClient

			id, err := brokerlistener.ParseListenerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := BrokerListenerModel{
				Name:              id.ListenerName,
				ResourceGroupName: id.ResourceGroupName,
				InstanceName:      id.InstanceName,
				BrokerName:        id.BrokerName,
			}

			if respModel := resp.Model; respModel != nil {
				if respModel.ExtendedLocation != nil {
					model.ExtendedLocation = flattenExtendedLocation(respModel.ExtendedLocation)
				}

				if respModel.Tags != nil {
					model.Tags = *respModel.Tags
				}

				if respModel.Properties != nil {
					flattenBrokerListenerProperties(respModel.Properties, &model)
					
					if respModel.Properties.ProvisioningState != nil {
						provisioningState := string(*respModel.Properties.ProvisioningState)
						model.ProvisioningState = &provisioningState
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r BrokerListenerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerListenerClient

			id, err := brokerlistener.ParseListenerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BrokerListenerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Check if anything actually changed before making API call
			if !metadata.ResourceData.HasChange("tags") && 
			   !metadata.ResourceData.HasChange("service_name") &&
			   !metadata.ResourceData.HasChange("service_type") &&
			   !metadata.ResourceData.HasChange("ports") {
				return nil
			}

			payload := brokerlistener.BrokerListenerPatchModel{}
			hasChanges := false

			// Only include tags if they changed
			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = &model.Tags
				hasChanges = true
			}

			// Only include properties if they changed
			if metadata.ResourceData.HasChange("service_name") ||
			   metadata.ResourceData.HasChange("service_type") ||
			   metadata.ResourceData.HasChange("ports") {
				payload.Properties = &brokerlistener.BrokerListenerPropertiesPatch{
					ServiceName: model.ServiceName,
					ServiceType: model.ServiceType,
					Ports:       expandBrokerListenerPorts(model.Ports),
				}
				hasChanges = true
			}

			// Only make API call if something actually changed
			if !hasChanges {
				return nil
			}

			if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r BrokerListenerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTOperations.BrokerListenerClient

			id, err := brokerlistener.ParseListenerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

// Helper functions for expand/flatten operations
func expandBrokerListenerProperties(model BrokerListenerModel) *brokerlistener.BrokerListenerProperties {
	props := &brokerlistener.BrokerListenerProperties{
		Ports: expandBrokerListenerPorts(model.Ports),
	}

	if model.ServiceName != nil {
		props.ServiceName = model.ServiceName
	}

	if model.ServiceType != nil {
		serviceType := brokerlistener.ServiceType(*model.ServiceType)
		props.ServiceType = &serviceType
	}

	return props
}

func expandBrokerListenerPorts(ports []BrokerListenerPortModel) *[]brokerlistener.ListenerPort {
	if len(ports) == 0 {
		return nil
	}

	result := make([]brokerlistener.ListenerPort, 0, len(ports))

	for _, port := range ports {
		listenerPort := brokerlistener.ListenerPort{
			Port: int64(port.Port),
		}

		if port.NodePort != nil {
			listenerPort.NodePort = func(i int) *int64 { v := int64(i); return &v }(*port.NodePort)
		}

		if port.Protocol != nil {
			protocol := brokerlistener.BrokerProtocolType(*port.Protocol)
			listenerPort.Protocol = &protocol
		}

		if port.AuthenticationRef != nil {
			listenerPort.AuthenticationRef = port.AuthenticationRef
		}

		if port.AuthorizationRef != nil {
			listenerPort.AuthorizationRef = port.AuthorizationRef
		}

		if port.Tls != nil {
			listenerPort.Tls = expandBrokerListenerTls(*port.Tls)
		}

		result = append(result, listenerPort)
	}

	return &result
}

func expandBrokerListenerTls(tls BrokerListenerTlsModel) *brokerlistener.TlsCertMethod {
	tlsMode := brokerlistener.TlsCertMethodMode(tls.Mode)
	result := &brokerlistener.TlsCertMethod{
		Mode: tlsMode,
	}

	if tls.CertManagerCertificateSpec != nil {
		result.CertManagerCertificateSpec = expandBrokerListenerCertManagerSpec(*tls.CertManagerCertificateSpec)
	}

	if tls.Manual != nil {
		result.Manual = expandBrokerListenerManual(*tls.Manual)
	}

	return result
}

func expandBrokerListenerCertManagerSpec(spec BrokerListenerCertManagerCertificateSpecModel) *brokerlistener.CertManagerCertificateSpec {
	result := &brokerlistener.CertManagerCertificateSpec{}

	if spec.Duration != nil {
		result.Duration = spec.Duration
	}

	if spec.SecretName != nil {
		result.SecretName = spec.SecretName
	}

	if spec.RenewBefore != nil {
		result.RenewBefore = spec.RenewBefore
	}

	if spec.IssuerRef != nil {
		result.IssuerRef = expandBrokerListenerIssuerRef(*spec.IssuerRef)
	}

	if spec.PrivateKey != nil {
		result.PrivateKey = expandBrokerListenerPrivateKey(*spec.PrivateKey)
	}

	if spec.San != nil {
		result.San = expandBrokerListenerSan(*spec.San)
	}

	return result
}

func expandBrokerListenerIssuerRef(issuerRef BrokerListenerIssuerRefModel) *brokerlistener.CertManagerIssuerRef {
	result := &brokerlistener.CertManagerIssuerRef{}

	if issuerRef.Group != nil {
		result.Group = issuerRef.Group
	}

	if issuerRef.Kind != nil {
		result.Kind = issuerRef.Kind
	}

	if issuerRef.Name != nil {
		result.Name = issuerRef.Name
	}

	return result
}

func expandBrokerListenerPrivateKey(privateKey BrokerListenerPrivateKeyModel) *brokerlistener.CertManagerPrivateKey {
	result := &brokerlistener.CertManagerPrivateKey{}

	if privateKey.Algorithm != nil {
		algorithm := brokerlistener.PrivateKeyAlgorithm(*privateKey.Algorithm)
		result.Algorithm = &algorithm
	}

	if privateKey.RotationPolicy != nil {
		rotationPolicy := brokerlistener.PrivateKeyRotationPolicy(*privateKey.RotationPolicy)
		result.RotationPolicy = &rotationPolicy
	}

	return result
}

func expandBrokerListenerSan(san BrokerListenerSanModel) *brokerlistener.SanForCert {
	result := &brokerlistener.SanForCert{}

	if san.Dns != nil {
		result.Dns = san.Dns
	}

	if san.Ip != nil {
		result.Ip = san.Ip
	}

	return result
}

func expandBrokerListenerManual(manual BrokerListenerManualModel) *brokerlistener.ManualCertMethod {
	result := &brokerlistener.ManualCertMethod{}

	if manual.SecretRef != nil {
		result.SecretRef = manual.SecretRef
	}

	return result
}

func flattenBrokerListenerProperties(props *brokerlistener.BrokerListenerProperties, model *BrokerListenerModel) {
	if props == nil {
		return
	}

	if props.ServiceName != nil {
		model.ServiceName = props.ServiceName
	}

	if props.ServiceType != nil {
		serviceType := string(*props.ServiceType)
		model.ServiceType = &serviceType
	}

	if props.Ports != nil {
		model.Ports = flattenBrokerListenerPorts(*props.Ports)
	}
}

func flattenBrokerListenerPorts(ports []brokerlistener.ListenerPort) []BrokerListenerPortModel {
	result := make([]BrokerListenerPortModel, 0, len(ports))

	for _, port := range ports {
		portModel := BrokerListenerPortModel{
			Port: int(port.Port),
		}

		if port.NodePort != nil {
			nodePort := int(*port.NodePort)
			portModel.NodePort = &nodePort
		}

		if port.Protocol != nil {
			protocol := string(*port.Protocol)
			portModel.Protocol = &protocol
		}

		if port.AuthenticationRef != nil {
			portModel.AuthenticationRef = port.AuthenticationRef
		}

		if port.AuthorizationRef != nil {
			portModel.AuthorizationRef = port.AuthorizationRef
		}

		if port.Tls != nil {
			portModel.Tls = flattenBrokerListenerTls(*port.Tls)
		}

		result = append(result, portModel)
	}

	return result
}

func flattenBrokerListenerTls(tls brokerlistener.TlsCertMethod) *BrokerListenerTlsModel {
	result := &BrokerListenerTlsModel{
		Mode: string(tls.Mode),
	}

	if tls.CertManagerCertificateSpec != nil {
		result.CertManagerCertificateSpec = flattenBrokerListenerCertManagerSpec(*tls.CertManagerCertificateSpec)
	}

	if tls.Manual != nil {
		result.Manual = flattenBrokerListenerManual(*tls.Manual)
	}

	return result
}

func flattenBrokerListenerCertManagerSpec(spec brokerlistener.CertManagerCertificateSpec) *BrokerListenerCertManagerCertificateSpecModel {
	result := &BrokerListenerCertManagerCertificateSpecModel{}

	if spec.Duration != nil {
		result.Duration = spec.Duration
	}

	if spec.SecretName != nil {
		result.SecretName = spec.SecretName
	}

	if spec.RenewBefore != nil {
		result.RenewBefore = spec.RenewBefore
	}

	if spec.IssuerRef != nil {
		result.IssuerRef = flattenBrokerListenerIssuerRef(*spec.IssuerRef)
	}

	if spec.PrivateKey != nil {
		result.PrivateKey = flattenBrokerListenerPrivateKey(*spec.PrivateKey)
	}

	if spec.San != nil {
		result.San = flattenBrokerListenerSan(*spec.San)
	}

	return result
}

func flattenBrokerListenerIssuerRef(issuerRef brokerlistener.CertManagerIssuerRef) *BrokerListenerIssuerRefModel {
	result := &BrokerListenerIssuerRefModel{}

	if issuerRef.Group != nil {
		result.Group = issuerRef.Group
	}

	if issuerRef.Kind != nil {
		result.Kind = issuerRef.Kind
	}

	if issuerRef.Name != nil {
		result.Name = issuerRef.Name
	}

	return result
}

func flattenBrokerListenerPrivateKey(privateKey brokerlistener.CertManagerPrivateKey) *BrokerListenerPrivateKeyModel {
	result := &BrokerListenerPrivateKeyModel{}

	if privateKey.Algorithm != nil {
		algorithm := string(*privateKey.Algorithm)
		result.Algorithm = &algorithm
	}

	if privateKey.RotationPolicy != nil {
		rotationPolicy := string(*privateKey.RotationPolicy)
		result.RotationPolicy = &rotationPolicy
	}

	return result
}

func flattenBrokerListenerSan(san brokerlistener.SanForCert) *BrokerListenerSanModel {
	result := &BrokerListenerSanModel{}

	if san.Dns != nil {
		result.Dns = san.Dns
	}

	if san.Ip != nil {
		result.Ip = san.Ip
	}

	return result
}

func flattenBrokerListenerManual(manual brokerlistener.ManualCertMethod) *BrokerListenerManualModel {
	result := &BrokerListenerManualModel{}

	if manual.SecretRef != nil {
		result.SecretRef = manual.SecretRef
	}

	return result
}