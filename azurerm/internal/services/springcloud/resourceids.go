package springcloud

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SpringCloudApp -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/Spring/spring1/apps/app1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SpringCloudCertificate -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/Spring/spring1/certificates/cert1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=SpringCloudService -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/Spring/spring1
