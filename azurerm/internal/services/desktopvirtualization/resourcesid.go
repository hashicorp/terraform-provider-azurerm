package desktopvirtualization

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=ApplicationGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/applicationGroups/applicationGroup1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Application -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/applicationGroups/applicationGroup1/applications/application1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=HostPool -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/hostPools/pool1 -rewrite=true
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Workspace -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DesktopVirtualization/workspaces/workspace1
