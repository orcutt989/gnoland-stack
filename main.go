package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Define labels for the app
		appLabels := pulumi.StringMap{
			"app": pulumi.String("gnoland"),
		}

		// Deploy the node StatefulSet and Service
		nodeServicePort := pulumi.Int(80)

		// Deploy nodeStatefulSet and nodeService
		statefulSet, nodeService, err := NewNode(ctx, appLabels, nodeServicePort)
		if err != nil {
			return err
		}

		// Export outputs if needed
		ctx.Export("nodeStatefulSetName", statefulSet.Metadata.Name())
		ctx.Export("nodeServiceName", nodeService.Metadata.Name())

		return nil
	})
}
