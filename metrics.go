package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewMetrics(ctx *pulumi.Context, namespace pulumi.StringInput, indexerServiceName pulumi.StringInput, indexerPort pulumi.IntInput) (*appsv1.StatefulSet, *corev1.Service, error) {
	appLabels := pulumi.StringMap{
		"app": pulumi.String("gnoland-metrics"),
	}
	var metricsListenPort = pulumi.Int(8080)
	var metricsPortName = pulumi.String("metrics")
	statefulSet, err := appsv1.NewStatefulSet(ctx, "gnoland-metrics", &appsv1.StatefulSetArgs{
		Spec: appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			ServiceName: pulumi.String("gnoland-metrics-service"),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: appLabels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  pulumi.String("gnoland-metrics"),
							Image: pulumi.String("ghcr.io/orcutt989/gnoland-metrics:latest"),
							Args: pulumi.StringArray{
								pulumi.Sprintf("-jsonrpc-url=http://%s.%s.svc.cluster.local:%d/query", indexerServiceName, namespace, indexerPort),
							},
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									ContainerPort: metricsListenPort,
									Name:          metricsPortName,
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
	metricsService, err := corev1.NewService(ctx, "gnoland-metrics-service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels: appLabels,
			Name:   pulumi.String("gnoland-metrics-service"),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: appLabels,
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       metricsPortName,
					Port:       metricsListenPort,
					TargetPort: metricsPortName,
				},
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return statefulSet, metricsService, nil
}
