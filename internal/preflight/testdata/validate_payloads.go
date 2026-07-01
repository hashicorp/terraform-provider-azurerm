package testdata

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func WebAppExample(subscriptionId string) map[string]any {
	exampleJSON := `{"parameters": {
    "subscriptionId": "34adfa4f-cedf-4dc0-ba29-b6d1a69ab345",
    "resourceGroupName": "testrg123",
    "name": "sitef6141",
    "api-version": "2024-11-01",
    "siteEnvelope": {
      "kind": "app",
      "location": "westeurope",
      "properties": {
        "serverFarmId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testrg123/providers/Microsoft.Web/serverfarms/DefaultAsp"
      }
    }
  }
}
`
	return normaliseExample(exampleJSON, subscriptionId)
}

func VirtualNetworkExampleMissingRequiredProperty(subscriptionId string) map[string]any {
	exampleJSON := `{"parameters": {
    "api-version": "2024-01-01",
    "subscriptionId": "00000000-0000-0000-0000-000000000000",
    "resourceGroupName": "rg1",
    "virtualNetworkName": "test-vnet",
    "parameters": {
      "properties": {
        "addressSpace": {},
        "flowTimeoutInMinutes": 10
      },
      "location": "westeurope"
    }
  }
}
`
	return normaliseExample(exampleJSON, subscriptionId)
}

// ContainerAppExample - Example in spec is broken, leaving this case here to fix later
func ContainerAppExample(subscriptionId string) map[string]any {
	exampleJson := fmt.Sprintf(`
{
  "parameters": {
    "subscriptionId": "%[1]s",
    "resourceGroupName": "rg",
    "containerAppName": "testcontainerapp0",
    "api-version": "2025-07-01",
    "containerAppEnvelope": {
      "location": "East US",
      "identity": {
        "type": "SystemAssigned,UserAssigned",
        "userAssignedIdentities": {
          "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity": {}
        }
      },
      "properties": {
        "environmentId": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube",
        "workloadProfileName": "My-GP-01",
        "configuration": {
          "ingress": {
            "external": true,
            "targetPort": 3000,
            "customDomains": [
              {
                "name": "www.my-name.com",
                "bindingType": "SniEnabled",
                "certificateId": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube/certificates/my-certificate-for-my-name-dot-com"
              },
              {
                "name": "www.my-other-name.com",
                "bindingType": "SniEnabled",
                "certificateId": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube/certificates/my-certificate-for-my-other-name-dot-com"
              }
            ],
            "traffic": [
              {
                "weight": 100,
                "revisionName": "testcontainerapp0-ab1234",
                "label": "production"
              }
            ],
            "ipSecurityRestrictions": [
              {
                "name": "Allow work IP A subnet",
                "description": "Allowing all IP's within the subnet below to access containerapp",
                "ipAddressRange": "192.168.1.1/32",
                "action": "Allow"
              },
              {
                "name": "Allow work IP B subnet",
                "description": "Allowing all IP's within the subnet below to access containerapp",
                "ipAddressRange": "192.168.1.1/8",
                "action": "Allow"
              }
            ],
            "stickySessions": {
              "affinity": "sticky"
            },
            "clientCertificateMode": "accept",
            "corsPolicy": {
              "allowedOrigins": [
                "https://a.test.com",
                "https://b.test.com"
              ],
              "allowedMethods": [
                "GET",
                "POST"
              ],
              "allowedHeaders": [
                "HEADER1",
                "HEADER2"
              ],
              "exposeHeaders": [
                "HEADER3",
                "HEADER4"
              ],
              "maxAge": 1234,
              "allowCredentials": true
            },
            "additionalPortMappings": [
              {
                "external": true,
                "targetPort": 1234
              },
              {
                "external": false,
                "targetPort": 2345,
                "exposedPort": 3456
              }
            ]
          },
          "dapr": {
            "enabled": true,
            "appPort": 3000,
            "appProtocol": "http",
            "httpReadBufferSize": 30,
            "httpMaxRequestSize": 10,
            "logLevel": "debug",
            "enableApiLogging": true,
            "appHealth": {
              "enabled": true,
              "path": "/health",
              "probeIntervalSeconds": 3,
              "probeTimeoutMilliseconds": 1000,
              "threshold": 3
            },
            "maxConcurrency": 10
          },
          "maxInactiveRevisions": 10,
          "service": {
            "type": "redis"
          },
          "runtime": {
            "java": {
              "enableMetrics": true
            }
          },
          "identitySettings": [
            {
              "identity": "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity",
              "lifecycle": "All"
            },
            {
              "identity": "system",
              "lifecycle": "Init"
            }
          ]
        },
        "template": {
          "containers": [
            {
              "image": "repo/testcontainerapp0:v1",
              "name": "testcontainerapp0",
              "probes": [
                {
                  "type": "Liveness",
                  "httpGet": {
                    "path": "/health",
                    "port": 8080,
                    "httpHeaders": [
                      {
                        "name": "Custom-Header",
                        "value": "Awesome"
                      }
                    ]
                  },
                  "initialDelaySeconds": 3,
                  "periodSeconds": 3
                }
              ],
              "volumeMounts": [
                {
                  "volumeName": "azurefile",
                  "mountPath": "/mnt/path1",
                  "subPath": "subPath1"
                },
                {
                  "volumeName": "nfsazurefile",
                  "mountPath": "/mnt/path2",
                  "subPath": "subPath2"
                }
              ]
            }
          ],
          "initContainers": [
            {
              "image": "repo/testcontainerapp0:v4",
              "name": "testinitcontainerApp0",
              "resources": {
                "cpu": 0.5,
                "memory": "1Gi"
              },
              "command": [
                "/bin/sh"
              ],
              "args": [
                "-c",
                "while true; do echo hello; sleep 10;done"
              ]
            }
          ],
          "scale": {
            "minReplicas": 1,
            "maxReplicas": 5,
            "cooldownPeriod": 350,
            "pollingInterval": 35,
            "rules": [
              {
                "name": "httpscalingrule",
                "custom": {
                  "type": "http",
                  "metadata": {
                    "concurrentRequests": "50"
                  }
                }
              },
              {
                "name": "servicebus",
                "custom": {
                  "type": "azure-servicebus",
                  "metadata": {
                    "queueName": "myqueue",
                    "namespace": "mynamespace",
                    "messageCount": "5"
                  },
                  "identity": "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity"
                }
              },
              {
                "name": "azure-queue",
                "azureQueue": {
                  "accountName": "account1",
                  "queueName": "queue1",
                  "queueLength": 1,
                  "identity": "system"
                }
              }
            ]
          },
          "serviceBinds": [
            {
              "serviceId": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/containerApps/redisService",
              "name": "redisService"
            }
          ],
          "volumes": [
            {
              "name": "azurefile",
              "storageType": "AzureFile",
              "storageName": "storage"
            },
            {
              "name": "nfsazurefile",
              "storageType": "NfsAzureFile",
              "storageName": "nfsStorage"
            }
          ]
        }
      }
    }
  },
  "responses": {
    "200": {
      "headers": {},
      "body": {
        "id": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/containerApps/testcontainerapp0",
        "name": "testcontainerapp0",
        "type": "Microsoft.App/containerApps",
        "location": "East US",
        "identity": {
          "type": "SystemAssigned,UserAssigned",
          "principalId": "24adfa4f-dedf-8dc0-ca29-b6d1a69ab319",
          "tenantId": "23adfa4f-eedf-1dc0-ba29-a6d1a69ab3d0",
          "userAssignedIdentities": {
            "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity": {
              "principalId": "74adfa4f-dedf-8dc0-ca29-b6d1a69ab312",
              "clientId": "14adfa4f-eedf-1dc0-ba29-a6d1a69ab3df"
            }
          }
        },
        "properties": {
          "provisioningState": "Succeeded",
          "runningStatus": "Running",
          "managedEnvironmentId": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube",
          "workloadProfileName": "My-GP-01",
          "latestRevisionFqdn": "testcontainerapp0-pjxhsye.demokube-t24clv0g.eastus.containerApps.k4apps.io",
          "latestReadyRevisionName": "testcontainerapp0-pjxhsye",
          "configuration": {
            "ingress": {
              "fqdn": "testcontainerapp0.demokube-t24clv0g.eastus.containerApps.k4apps.io",
              "external": true,
              "targetPort": 3000,
              "transport": "auto",
              "customDomains": [
                {
                  "name": "www.my-name.com",
                  "bindingType": "SniEnabled",
                  "certificateId": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube/certificates/my-certificate-for-my-name-dot-com"
                },
                {
                  "name": "www.my-other-name.com",
                  "bindingType": "SniEnabled",
                  "certificateId": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube/certificates/my-certificate-for-my-other-name-dot-com"
                }
              ],
              "traffic": [
                {
                  "weight": 80,
                  "revisionName": "testcontainerapp0-ab1234"
                },
                {
                  "weight": 20,
                  "revisionName": "testcontainerapp0-ab4321",
                  "label": "staging"
                }
              ],
              "ipSecurityRestrictions": [
                {
                  "name": "Allow work IP A subnet",
                  "description": "Allowing all IP's within the subnet below to access containerapp",
                  "ipAddressRange": "192.168.1.1/32",
                  "action": "Allow"
                },
                {
                  "name": "Allow work IP B subnet",
                  "description": "Allowing all IP's within the subnet below to access containerapp",
                  "ipAddressRange": "192.168.1.1/8",
                  "action": "Allow"
                }
              ],
              "stickySessions": {
                "affinity": "sticky"
              }
            },
            "dapr": {
              "enabled": true,
              "appPort": 3000,
              "appProtocol": "http",
              "httpReadBufferSize": 30,
              "appHealth": {
                "enabled": true,
                "path": "/health",
                "probeIntervalSeconds": 3,
                "probeTimeoutMilliseconds": 1000,
                "threshold": 3
              },
              "maxConcurrency": 10
            },
            "runtime": {
              "java": {
                "enableMetrics": true
              }
            },
            "identitySettings": [
              {
                "identity": "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity",
                "lifecycle": "All"
              },
              {
                "identity": "system",
                "lifecycle": "Init"
              }
            ]
          },
          "template": {
            "containers": [
              {
                "image": "repo/testcontainerapp0:v4",
                "name": "testcontainerapp0",
                "resources": {
                  "cpu": 0.5,
                  "memory": "1Gi"
                },
                "probes": [
                  {
                    "type": "Liveness",
                    "httpGet": {
                      "path": "/health",
                      "port": 8080,
                      "httpHeaders": [
                        {
                          "name": "Custom-Header",
                          "value": "Awesome"
                        }
                      ]
                    },
                    "initialDelaySeconds": 3,
                    "periodSeconds": 3
                  }
                ],
                "volumeMounts": [
                  {
                    "volumeName": "azurefile",
                    "mountPath": "/mnt/path1",
                    "subPath": "subPath1"
                  },
                  {
                    "volumeName": "nfsazurefile",
                    "mountPath": "/mnt/path2",
                    "subPath": "subPath2"
                  }
                ]
              }
            ],
            "initContainers": [
              {
                "image": "repo/testcontainerapp0:v4",
                "name": "testinitcontainerApp0",
                "resources": {
                  "cpu": 0.5,
                  "memory": "1Gi"
                },
                "command": [
                  "/bin/sh"
                ],
                "args": [
                  "-c",
                  "while true; do echo hello; sleep 10;done"
                ]
              }
            ],
            "scale": {
              "minReplicas": 1,
              "maxReplicas": 5,
              "cooldownPeriod": 350,
              "pollingInterval": 35,
              "rules": [
                {
                  "name": "httpscalingrule",
                  "http": {
                    "metadata": {
                      "concurrentRequests": "50"
                    }
                  }
                },
                {
                  "name": "servicebus",
                  "custom": {
                    "type": "azure-servicebus",
                    "metadata": {
                      "queueName": "myqueue",
                      "namespace": "mynamespace",
                      "messageCount": "5"
                    },
                    "identity": "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity"
                  }
                },
                {
                  "name": "azure-queue",
                  "azureQueue": {
                    "accountName": "account1",
                    "queueName": "queue1",
                    "queueLength": 1,
                    "identity": "system"
                  }
                }
              ]
            },
            "volumes": [
              {
                "name": "azurefile",
                "storageType": "AzureFile",
                "storageName": "storage"
              },
              {
                "name": "nfsazurefile",
                "storageType": "NfsAzureFile",
                "storageName": "nfsStorage"
              }
            ]
          },
          "eventStreamEndpoint": "testEndpoint"
        }
      }
    },
    "201": {
      "headers": {},
      "body": {
        "id": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/containerApps/testcontainerapp0",
        "name": "testcontainerapp0",
        "type": "Microsoft.App/containerApps",
        "location": "East US",
        "identity": {
          "type": "SystemAssigned,UserAssigned",
          "principalId": "24adfa4f-dedf-8dc0-ca29-b6d1a69ab319",
          "tenantId": "23adfa4f-eedf-1dc0-ba29-a6d1a69ab3d0",
          "userAssignedIdentities": {
            "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity": {
              "principalId": "74adfa4f-dedf-8dc0-ca29-b6d1a69ab312",
              "clientId": "14adfa4f-eedf-1dc0-ba29-a6d1a69ab3df"
            }
          }
        },
        "properties": {
          "provisioningState": "InProgress",
          "runningStatus": "Running",
          "managedEnvironmentId": "/subscriptions/%[1]s/resourceGroups/rg/providers/Microsoft.App/managedEnvironments/demokube",
          "latestRevisionFqdn": "testcontainerapp0-pjxhsye.demokube-t24clv0g.eastus.containerApps.k4apps.io",
          "configuration": {
            "ingress": {
              "fqdn": "testcontainerapp0.demokube-t24clv0g.eastus.containerApps.k4apps.io",
              "external": true,
              "targetPort": 3000,
              "transport": "auto",
              "traffic": [
                {
                  "weight": 80,
                  "revisionName": "testcontainerapp0-ab1234"
                },
                {
                  "weight": 20,
                  "revisionName": "testcontainerapp0-ab4321",
                  "label": "staging"
                }
              ]
            },
            "dapr": {
              "enabled": true,
              "appPort": 3000,
              "appProtocol": "http",
              "httpReadBufferSize": 30,
              "appHealth": {
                "enabled": true,
                "path": "/health",
                "probeIntervalSeconds": 3,
                "probeTimeoutMilliseconds": 1000,
                "threshold": 3
              },
              "maxConcurrency": 10
            },
            "runtime": {
              "java": {
                "enableMetrics": true
              }
            },
            "identitySettings": [
              {
                "identity": "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity",
                "lifecycle": "All"
              },
              {
                "identity": "system",
                "lifecycle": "Init"
              }
            ]
          },
          "template": {
            "containers": [
              {
                "image": "repo/testcontainerapp0:v4",
                "name": "testcontainerapp0",
                "resources": {
                  "cpu": 0.5,
                  "memory": "1Gi"
                },
                "probes": [
                  {
                    "type": "Liveness",
                    "httpGet": {
                      "path": "/health",
                      "port": 8080,
                      "httpHeaders": [
                        {
                          "name": "Custom-Header",
                          "value": "Awesome"
                        }
                      ]
                    },
                    "initialDelaySeconds": 3,
                    "periodSeconds": 3
                  }
                ],
                "volumeMounts": [
                  {
                    "volumeName": "azurefile",
                    "mountPath": "/mnt/path1",
                    "subPath": "subPath1"
                  },
                  {
                    "volumeName": "nfsazurefile",
                    "mountPath": "/mnt/path2",
                    "subPath": "subPath2"
                  }
                ]
              }
            ],
            "initContainers": [
              {
                "image": "repo/testcontainerapp0:v4",
                "name": "testinitcontainerApp0",
                "resources": {
                  "cpu": 0.5,
                  "memory": "1Gi"
                },
                "command": [
                  "/bin/sh"
                ],
                "args": [
                  "-c",
                  "while true; do echo hello; sleep 10;done"
                ]
              }
            ],
            "scale": {
              "minReplicas": 1,
              "maxReplicas": 5,
              "cooldownPeriod": 350,
              "pollingInterval": 35,
              "rules": [
                {
                  "name": "httpscalingrule",
                  "http": {
                    "metadata": {
                      "concurrentRequests": "50"
                    }
                  }
                },
                {
                  "name": "servicebus",
                  "custom": {
                    "type": "azure-servicebus",
                    "metadata": {
                      "queueName": "myqueue",
                      "namespace": "mynamespace",
                      "messageCount": "5"
                    },
                    "identity": "/subscriptions/%[1]s/resourcegroups/rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/myidentity"
                  }
                },
                {
                  "name": "azure-queue",
                  "azureQueue": {
                    "accountName": "account1",
                    "queueName": "queue1",
                    "queueLength": 1,
                    "identity": "system"
                  }
                }
              ]
            },
            "volumes": [
              {
                "name": "azurefile",
                "storageType": "AzureFile",
                "storageName": "storage"
              },
              {
                "name": "nfsazurefile",
                "storageType": "NfsAzureFile",
                "storageName": "nfsStorage"
              }
            ]
          },
          "eventStreamEndpoint": "testEndpoint"
        }
      }
    }
  }
}`, subscriptionId)

	return normaliseExample(exampleJson, subscriptionId)
}

// RedisEnterpriseExample - Example in spec is broken, leaving this case here to fix later
func RedisEnterpriseExample(subscriptionId string) map[string]any {
	exampleJSON := `{
"parameters": {
    "clusterName": "testRedisEnterpriseValid",
    "resourceGroupName": "testResourceGroup",
    "api-version": "2025-07-01",
    "subscriptionId": "00000000-0000-0000-0000-000000000000",
    "parameters": {
      "location": "uksouth",
      "sku": {
        "name": "Balanced_B1"
      },
      "zones": [
        "1",
        "2",
        "3"
      ],
      "identity": {
        "type": "UserAssigned",
        "userAssignedIdentities": {
          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/your-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/your-identity": {}
        }
      },
      "properties": {
        "minimumTlsVersion": "1.2",
        "encryption": {
          "customerManagedKeyEncryption": {
            "keyEncryptionKeyIdentity": {
              "identityType": "userAssignedIdentity",
              "userAssignedIdentityResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/your-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/your-identity"
            },
            "keyEncryptionKeyUrl": "https://your-kv.vault.azure.net/keys/your-key/your-key-version"
          }
        },
        "publicNetworkAccess": "Disabled"
      },
      "tags": {
        "tag1": "value1"
      }
    }
  }
}
`
	return normaliseExample(exampleJSON, subscriptionId)
}

func normaliseExample(input string, subscriptionId string) map[string]any {
	result := make(map[string]any)

	input = strings.ReplaceAll(input, "\"subscriptionId\": \"00000000-0000-0000-0000-000000000000\"", "\"subscriptionId\": \"%s\"")

	re := regexp.MustCompile(`(?i)/subscriptions/[a-f0-9-]{36}`)
	input = re.ReplaceAllString(input, "/subscriptions/%[1]s")

	input = fmt.Sprintf(input, subscriptionId)

	trimmed := strings.ReplaceAll(input, "\n", "")
	trimmed = strings.ReplaceAll(trimmed, "\t", "")

	err := json.Unmarshal([]byte(trimmed), &result)
	if err != nil {
		panic(err)
	}

	return result
}
