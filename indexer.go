package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewIndexer(ctx *pulumi.Context, nodeServiceName pulumi.StringInput, nodeServicePort pulumi.IntInput) (*appsv1.StatefulSet, error) {
	statefulSet, err := appsv1.NewStatefulSet(ctx, "tx-indexer-statefulset", &appsv1.StatefulSetArgs{
		Spec: appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String("tx-indexer-statefulset"),
				},
			},
			ServiceName: nodeServiceName,
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String("tx-indexer-service"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  pulumi.String("tx-indexer"),
							Image: pulumi.String("ghcr.io/orcutt989/tx-indexer:main"),
							Command: pulumi.StringArray{
								pulumi.String("indexer"),
								pulumi.String("start"),
								pulumi.Sprintf("--remote http://%s:%d --db-path indexer-db", nodeServiceName, nodeServicePort),
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return statefulSet, nil
}
