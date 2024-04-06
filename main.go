package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Static port for now for indexer and supernova to find RPC of node
		nodeServicePort := pulumi.Int(26657)

		// Deploy gno.land node
		_, _, err := NewNode(ctx)
		if err != nil {
			return err
		}

		// Set the namespace and service name as static values
		// TODO dynamic update
		namespace := pulumi.String("default")
		serviceName := pulumi.String("node-service")

		// Deploys the tx-indexer
		_, _, err = NewIndexer(ctx, namespace, serviceName, nodeServicePort)
		if err != nil {
			return err
		}

		var sleepSeconds = pulumi.Int(10)
		// Runs a Supernova job
		_, err = NewSupernova(ctx, namespace, serviceName, nodeServicePort, sleepSeconds)
		if err != nil {
			return err
		}

		// Deploys gnoland-metrics
		indexerPort := pulumi.Int(8545)
		indexerServiceName := pulumi.String("tx-indexer-service")
		_, _, err = NewMetrics(ctx, namespace, indexerServiceName, indexerPort)
		if err != nil {
			return err
		}

		return nil
	})
}
