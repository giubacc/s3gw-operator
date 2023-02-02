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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	s3v1 "github.com/giubacc/s3gw-operator/api/v1"
	"github.com/sirupsen/logrus"
)

var DebugLogger *logrus.Logger

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
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *BucketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var bucket s3v1.Bucket
	if err := r.Get(ctx, req.NamespacedName, &bucket); err != nil {
		if err != nil {
			if errors.IsNotFound(err) {
				//This is a cancellation.
				//Cancellations should be handled with a finalizer like in
				//Kubewarden: https://github.com/kubewarden/kubewarden-controller

				//*********************************
				//This approach is only for demoing
				//*********************************

				DebugLogger.Tracef("deleting bucket:%s", req.Name)
				bucket.Spec.Name = req.Name
				if err := S3Manager.EnsureBucketDeleted(ctx, &bucket); err != nil {
					log.Error(err, "unable to delete bucket")
					return ctrl.Result{}, err
				}
			} else {
				// Error reading the object - requeue the request.
				return reconcile.Result{}, err
			}
		}
	} else {
		DebugLogger.Tracef("ensuring bucket:%s", req.Name)
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
