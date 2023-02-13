/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	telcolabs1alpha1 "tlhelper_operator/api/v1alpha1"
)

// TlhelperReconciler reconciles a Tlhelper object
type TlhelperReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=telco.labs,resources=tlhelpers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=telco.labs,resources=tlhelpers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=telco.labs,resources=tlhelpers/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps,resources=deployments;replicasets;,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods;services;services/finalizers;endpoints;events;configmaps,verbs=create;update;get;list;watch;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Tlhelper object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TlhelperReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("Tlhelper", req.NamespacedName)

	// Fetch the Thelper instance
	instance := &telcolabs1alpha1.Tlhelper{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Tlhelper resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if this Deployment already exists
	var result *reconcile.Result
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		result, err = r.ensureDeployment(req, instance, r.tlhelperDeployment(instance), ctx)
		if result != nil {
			log.Error(err, "Deployment Not ready")
			return *result, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// This point, we have the deployment object created
	// Ensure the deployment size is same as the spec
	replicas := instance.Spec.Replicas
	if *found.Spec.Replicas != *replicas {
		found.Spec.Replicas = replicas
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated return and requeue
		// Requeue for any reason other than an error
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the app status with pod names
	// List the pods for this app's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(instance.GetLabels()),
	}

	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Falied to list pods", "Tlhelper.Namespace", instance.Namespace, "Tlhelper.Name", instance.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Pods if needed
	if !reflect.DeepEqual(podNames, instance.Status.Pods) {
		instance.Status.Pods = podNames
		err := r.Update(ctx, instance)
		if err != nil {
			log.Error(err, "Failed to update InhouseApp status")
			return ctrl.Result{}, err
		}
	}

	// Check if this Service already exists
	result, err = r.ensureService(req, instance, r.tlhelperService(instance), ctx)
	if result != nil {
		log.Error(err, "Service Not ready")
		return *result, err
	}

	// Create ConfigMap first
	cm, err := configMapForTlhelperController(instance, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, cm, configMapMutateFn(cm, instance.Spec.ConfigMapData))
	if err != nil {
		log.Error(err, "ConfigMap Not ready")
		return ctrl.Result{}, err
	}

	// Deployment and Service already exists - don't requeue
	log.Info("Skip reconcile: Deployment and service already exists",
		"Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)

	return ctrl.Result{}, nil
}

// Utility function to iterate over pods and return the names slice
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *TlhelperReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&telcolabs1alpha1.Tlhelper{}).
		Complete(r)
}
