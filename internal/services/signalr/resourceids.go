package signalr

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebPubsub -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/webPubSub/Webpubsub1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebPubsubHub -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/webPubSub/Webpubsub1/hubs/Webpubsubhub1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=WebPubsubSharedPrivateLinkResource -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.SignalRService/webPubSub/Webpubsub1/sharedPrivateLinkResources/resource1
