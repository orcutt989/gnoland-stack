package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewNode(ctx *pulumi.Context) (*appsv1.StatefulSet, *corev1.Service, error) {
	// Deploy the StatefulSet for pods
	// Define labels for the app
	appLabels := pulumi.StringMap{
		"app": pulumi.String("gnoland"),
	}
	var port = pulumi.Int(26657)
	var portName = pulumi.String("gnoland-rpc")
	statefulSet, err := appsv1.NewStatefulSet(ctx, "node", &appsv1.StatefulSetArgs{
		Spec: appsv1.StatefulSetSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: appLabels,
					Name:   pulumi.String("node"),
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  pulumi.String("gnoland"),
							Image: pulumi.String("ghcr.io/gnolang/gno:latest"),
							Args: pulumi.StringArray{
								pulumi.String("sh"),
								pulumi.String("-c"),
								pulumi.String("gnoland start --skip-start=true && sed -i \"s#laddr = \\\".*:26657\\\"#laddr = \\\"tcp://0.0.0.0:26657\\\"#\" ./testdir/config/config.toml && gnoland start"),
							},
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									ContainerPort: port,
									Name:          portName,
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
	nodeService, err := corev1.NewService(ctx, "node-service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels: appLabels,
			Name:   pulumi.String("node-service"),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: appLabels,
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Name:       portName,
					Port:       pulumi.Int(26657),
					TargetPort: portName,
				},
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return statefulSet, nodeService, nil
}
