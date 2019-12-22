// Copyright 2019 IBM Corp
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package securecontainerimage

import (
	"bytes"
	"context"
	"text/template"

	securecontainersv1alpha1 "github.com/ibm/raksh/pkg/apis/securecontainers/v1alpha1"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_securecontainerimage")

const script = `#!/bin/sh
set -euo pipefail

if [ ! -d "/var/lib/securecontainer/" ]
then
    echo "ERROR: Directory /var/lib/securecontainer/ does not exists in the job, \
please verify the field imageDir in SecureContainerImageConfig. It's necessary to mount this directory across all worker nodes."
    exit 1
fi

if [ -d "/var/lib/securecontainer/{{.PodName}}" ]
then
    echo "INFO: Directory /var/lib/securecontainer/{{.PodName}} is getting overwritten."
else
    mkdir -p /var/lib/securecontainer/{{.PodName}}
fi

if [ ! -d "/securecontainer/" ]
then
    echo "ERROR: Directory /securecontainer/ doesn't exit, please check the {{.SecureContainerVMImage}}"
    exit 1
fi

if [ ! -f "/securecontainer/initrd.img" ]
then
    echo "ERROR: Unable to find initrd.img in the folder /securecontainer/, please check the {{.SecureContainerVMImage}}"
    exit 1
fi

if [ ! -f "/securecontainer/vmlinux" ]
then
    echo "ERROR: Unable to find vmlinux in the folder /securecontainer/, please check the {{.SecureContainerVMImage}}"
    exit 1
fi

echo "INFO: Copying the content from /securecontainer/ to /var/lib/securecontainer/{{.PodName}}"

ls -l /securecontainer/*
cp /securecontainer/initrd.img /var/lib/securecontainer/{{.PodName}}/initrd.img
cp /securecontainer/vmlinux /var/lib/securecontainer/{{.PodName}}/vmlinux`

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SecureContainerImage Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSecureContainerImage{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("securecontainerimage-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SecureContainerImage
	err = c.Watch(&source.Kind{Type: &securecontainersv1alpha1.SecureContainerImage{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Jobs and requeue the owner SecureContainerImage
	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &securecontainersv1alpha1.SecureContainerImage{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSecureContainerImage{}

// ReconcileSecureContainerImage reconciles a SecureContainerImage object
type ReconcileSecureContainerImage struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SecureContainerImage object and makes changes based on the state read
// and what is in the SecureContainerImage.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Job as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSecureContainerImage) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SecureContainerImage")

	// Fetch the SecureContainerImage instance
	instance := &securecontainersv1alpha1.SecureContainerImage{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	secureContainerImageConfigRef := instance.Spec.SecureContainerImageConfigRef
	secureContainerImageConfig := &securecontainersv1alpha1.SecureContainerImageConfig{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: secureContainerImageConfigRef.Name, Namespace: request.Namespace}, secureContainerImageConfig)
	if err == nil {
		instance.Spec.SecureContainerImageConfigSpec = secureContainerImageConfig.Spec
	}

	err = r.client.Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Define a new Job object
	job, err := newJobForCR(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Set SecureContainerImage instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, job, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Job already exists
	found := &batchv1.Job{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		err = r.client.Create(context.TODO(), job)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Job created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Job already exists - don't requeue
	reqLogger.Info("Skip reconcile: Job already exists", "Job.Namespace", found.Namespace, "Job.Name", found.Name)
	return reconcile.Result{}, nil
}

func newJobForCR(cr *securecontainersv1alpha1.SecureContainerImage) (*batchv1.Job, error) {
	labels := map[string]string{
		"app": cr.Name,
	}
	hostPathDirectory := corev1.HostPathDirectory
	buf := new(bytes.Buffer)
	pod := struct {
		PodName                string
		SecureContainerVMImage string
	}{
		cr.Name,
		cr.Spec.VMImage,
	}
	t := template.Must(template.New("podscript").Parse(script))
	err := t.Execute(buf, pod)
	if err != nil {
		return nil, err
	}
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-job",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "securecontainer-image",
							Image: cr.Spec.VMImage,
							Command: []string{
								"/bin/sh",
								"-c",
								buf.String(),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "securecontainer-vol",
									MountPath: "/var/lib/securecontainer",
								},
							},
						},
					},
					ImagePullSecrets: cr.Spec.ImagePullSecrets,
					RestartPolicy:    corev1.RestartPolicyOnFailure,
					Volumes: []corev1.Volume{
						{
							Name: "securecontainer-vol",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: cr.Spec.SecureContainerImageConfigSpec.ImageDir,
									Type: &hostPathDirectory,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
