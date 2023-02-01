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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	s3v1 "github.com/giubacc/s3gw-operator/api/v1"
)

// S3 manager
var S3Manager *Manager

// BucketReconciler reconciles a Bucket object
type BucketReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=s3.s3gw.io,resources=buckets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=s3.s3gw.io,resources=buckets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=s3.s3gw.io,resources=buckets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Bucket object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *BucketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var bucket s3v1.Bucket
	if err := r.Get(ctx, req.NamespacedName, &bucket); err != nil {
		log.Error(err, "unable to fetch Bucket")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		if err := S3Manager.EnsureBucketCreated(ctx, &bucket); err != nil {
			log.Error(err, "unable to ensure bucket")
			bucket.Status.Status = s3v1.Error
			if err := r.Status().Update(ctx, &bucket); err != nil {
				log.Error(err, "unable to update Bucket status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		} else {
			bucket.Status.Status = s3v1.Created
			if err := r.Status().Update(ctx, &bucket); err != nil {
				log.Error(err, "unable to update Bucket status")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BucketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&s3v1.Bucket{}).
		Complete(r)
}
