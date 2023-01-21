package logic

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IntegrationAccountSession -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/integrationAccount1/sessions/session1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=IntegrationServiceEnvironment -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Logic/integrationServiceEnvironments/ise1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=LogicAppStandard -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Workflow -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Logic/workflows/workflow1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Trigger -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Logic/workflows/workflow1/triggers/trigger1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Action -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Logic/workflows/workflow1/actions/action1
