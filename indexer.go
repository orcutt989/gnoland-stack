package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewIndexer(ctx *pulumi.Context, namespace pulumi.StringInput, nodeServiceName pulumi.StringInput, nodeServicePort pulumi.IntInput) (*appsv1.StatefulSet, *corev1.Service, error) {
	appLabels := pulumi.StringMap{
		"app": pulumi.String("tx-indexer"),
	}
	var indexerListenPort = pulumi.Int(8545)
	var indexerPortName = pulumi.String("indexer-rpc")
	statefulSet, err := appsv1.NewStatefulSet(ctx, "tx-indexer", &appsv1.StatefulSetArgs{
		Spec: appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			ServiceName: pulumi.String("tx-indexer-service"),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: appLabels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  pulumi.String("tx-indexer"),
							Image: pulumi.String("ghcr.io/orcutt989/tx-indexer:main"),
							Command: pulumi.StringArray{
								pulumi.String("sh"),
								pulumi.String("-c"),
								pulumi.Sprintf("indexer start --remote http://%s.%s.svc.cluster.local:%d --db-path indexer-db --listen-address 0.0.0.0:%d", nodeServiceName, namespace, nodeServicePort, indexerListenPort),
							},
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									ContainerPort: indexerListenPort,
									Name:          indexerPortName,
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	// Create a Kubernetes Service for the pods
	indexerService, err := corev1.NewService(ctx, "tx-indexer-service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels: appLabels,
			Name:   pulumi.String("tx-indexer-service"),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: appLabels,
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       indexerPortName,
					Port:       indexerListenPort,
					TargetPort: indexerPortName,
				},
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return statefulSet, indexerService, nil
}
