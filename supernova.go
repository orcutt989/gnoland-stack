package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewSupernova(ctx *pulumi.Context, namespace pulumi.StringInput, nodeServiceName pulumi.StringInput, nodeServicePort pulumi.IntInput, sleepSeconds pulumi.IntInput) (*appsv1.Deployment, error) {
	// Define labels for the app
	appLabels := pulumi.StringMap{
		"app": pulumi.String("supernova"),
	}

	// Create a Kubernetes Deployment for the pod
	deployment, err := appsv1.NewDeployment(ctx, "supernova-deployment", &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:   pulumi.String("supernova-deployment"),
			Labels: appLabels,
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: appLabels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							ImagePullPolicy: pulumi.String("Always"),
							Name:            pulumi.String("supernova"),
							Image:           pulumi.String("ghcr.io/orcutt989/supernova-container:main"),
							Args: pulumi.StringArray{
								pulumi.String("sh"),
								pulumi.String("-c"),
								pulumi.Sprintf("while true; do ./supernova -sub-accounts 5 -transactions 100 -url http://%s.%s.svc.cluster.local:%d -mnemonic \"source bonus chronic canvas draft south burst lottery vacant surface solve popular case indicate oppose farm nothing bullet exhibit title speed wink action roast\" -output result.json; sleep %d; done", nodeServiceName, namespace, nodeServicePort, sleepSeconds),
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

	return deployment, nil
}
