package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	telcolabs1alpha1 "tlhelper_operator/api/v1alpha1"
)

func labels(v *telcolabs1alpha1.Tlhelper, tier string) map[string]string {
	// Fetches and sets labels

	return map[string]string{
		"app":  "tlhelper",
		"tier": tier,
	}
}

// ensureDeployment ensures Deployment resource presence in given namespace.
func (r *TlhelperReconciler) ensureDeployment(request reconcile.Request,
	instance *telcolabs1alpha1.Tlhelper,
	dep *appsv1.Deployment,
	ctx context.Context,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Create(ctx, dep)

		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}

	return nil, nil
}

// backendDeployment is a code for Creating Deployment
func (r *TlhelperReconciler) tlhelperDeployment(instance *telcolabs1alpha1.Tlhelper) *appsv1.Deployment {

	labels := labels(instance, "tlhelper")
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           generateImage(instance.Spec.Image.Repository, instance.Spec.Image.Tag),
						ImagePullPolicy: corev1.PullPolicy(instance.Spec.Image.PullPolicy),
						Name:            instance.Name,
						Args:            generatePodArgs(instance),
						Env: []corev1.EnvVar{
							{
								Name: "POSTGRES_HOST",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: instance.Name,
										},
										Key: "postgresHost",
									},
								},
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
							Protocol:      "TCP",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, dep, r.Scheme)
	return dep
}
