/*


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

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	customimagedeployv1 "example.com/custom-image-deploy/api/v1"
)

// CustomImageDeployReconciler reconciles a CustomImageDeploy object
type CustomImageDeployReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=customimagedeploy.example.com,resources=customimagedeploys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=customimagedeploy.example.com,resources=customimagedeploys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=services;configmaps;pods;ingresses;deployments;serviceaccounts,verbs=get;list;watch;create;update;patch;delete

func (r *CustomImageDeployReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	log := r.Log.WithValues("customimagedeploy", req.NamespacedName)

	//
	cid := &customimagedeployv1.CustomImageDeploy{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, cid)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile req.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("CustomImageDeploy resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the req.
		log.Error(err, "Failed to get Memcached.")
		return ctrl.Result{}, err
	}

	// check if Deployment already exists, if not crate a new one
	deployment := &appsv1.Deployment{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: cid.Name, Namespace: cid.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deploymentForCustomImageDeploy(cid)
		log.Info("Creating a new deployment.", "Namespace: ", dep.Namespace, "Name: ", dep.Name)
		err = r.Client.Create(context.TODO(), dep)
		if err != nil {
			log.Error(err, "Failed to create a new deployment", "Namespace: ", dep.Namespace, "Name: ", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		// NOTE: that the requeue is made with the purpose to provide the deployment object for the next step
		// to ensure the deployment size is the same as the spec.
		// Also, you could GET the deployment object again instead of requeue if you wish.
		// See more over it here: https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler
		return reconcile.Result{}, nil
	}

	// ensure the size
	size := cid.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		err = r.Client.Update(context.TODO(), deployment)
		if err != nil {
			log.Error(err, "Failed to udpate deployment", "Namespace: ", deployment.Namespace, "Name: ", deployment.Name)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// deploymentForMemcached returns a memcached Deployment object
func (r *CustomImageDeployReconciler) deploymentForCustomImageDeploy(c *customimagedeployv1.CustomImageDeploy) *appsv1.Deployment {
	replicas := c.Spec.Size
	image := c.Spec.Image
	name := c.Name
	port := c.Spec.Port

	ls := labelsForCustomImageDeploy(name)

	dep := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      c.Name,
			Namespace: c.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &v1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  name,
						// Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []corev1.ContainerPort{{
							ContainerPort: port,
							// Name:          name, // Name is optinal, no more than 15 characters
						}},
					}},
				},
			},
		},
	}

	// Set Memcached instance as the owner of the Deployment.
	// ctrl.SetControllerReference(m, dep, r.Scheme) //todo check how to get the schema

	return dep
}

// labelsForCustomImageDeploy returns the labels for selecting the resources
// belonging to the given custom-image-deploy CR name.
func labelsForCustomImageDeploy(name string) map[string]string {
	return map[string]string{"app": name, "managed_by": "custom-image-deploy"}
}

func (r *CustomImageDeployReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&customimagedeployv1.CustomImageDeploy{}).
		Complete(r)
}
