package controller

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/objectiser/scribble-operator/pkg/apis/io/v1alpha1"
)

// AddToManager creates a new Scribble Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func AddToManager(m manager.Manager) error {
	return add(m, newReconciler(m))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileScribble{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("scribble-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Scribble
	err = c.Watch(&source.Kind{Type: &v1alpha1.Scribble{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileScribble{}

// ReconcileScribble reconciles a Scribble object
type ReconcileScribble struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Scribble object and makes changes based on the state read
// and what is in the Scribble.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileScribble) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.WithFields(log.Fields{
		"namespace": request.Namespace,
		"name":      request.Name,
	}).Print("Reconciling Scribble")

	// Fetch the Scribble instance
	instance := &v1alpha1.Scribble{}
	err := r.client.Get(context.Background(), request.NamespacedName, instance)
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

	// workaround for https://github.com/kubernetes-sigs/controller-runtime/issues/202
	// see also: https://github.com/kubernetes-sigs/controller-runtime/pull/212
	// once there's a version incorporating the PR above, the manual setting of the GKV can be removed
	instance.APIVersion = fmt.Sprintf("%s/%s", v1alpha1.SchemeGroupVersion.Group, v1alpha1.SchemeGroupVersion.Version)
	instance.Kind = "Scribble"

	// wait for all the dependencies to succeed
	if err := r.handleDependencies(); err != nil {
		return reconcile.Result{}, err
	}

	created, err := r.handleCreate()
	if err != nil {
		log.WithField("instance", instance).WithError(err).Error("failed to create")
		return reconcile.Result{}, err
	}

	if created {
		log.WithField("name", instance.Name).Info("Configured Scribble instance")
	}

	if err := r.handleUpdate(); err != nil {
		return reconcile.Result{}, err
	}

	// we store back the changed CR, so that what is stored reflects what is being used
	if err := r.client.Update(context.Background(), instance); err != nil {
		log.WithError(err).Error("failed to update")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileScribble) handleCreate() (bool, error) {
	os := []runtime.Object{}

	// TODO: Need to create deployment for each monitor - but it is dependent
	// upon the role to service bindings - which could be statically defined
	// in the CR - but then also needs to be verified that the roles all exist
	// and that atleast one service binding exists per role in the protocol.

	created := false
	for _, obj := range os {
		err := r.client.Create(context.Background(), obj)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			log.WithError(err).Error("failed to create")
			return false, err
		}

		if err == nil {
			created = true
		}
	}

	return created, nil
}

func (r *ReconcileScribble) handleUpdate() error {
	// Not currently handled
	return nil
}

func (r *ReconcileScribble) handleDependencies() error {
	/* TODO: Need a dependency that will validate the protocol description and possibly
	   return the list of roles, which can be verified against those defined in the CR.

		for _, dep := range str.Dependencies() {
			err := r.client.Create(context.Background(), &dep)
			if err != nil && !apierrors.IsAlreadyExists(err) {
				log.WithError(err).Error("failed to create")
				return err
			}

			// we probably want to add a couple of seconds to this deadline, but for now, this should be sufficient
			deadline := time.Duration(*dep.Spec.ActiveDeadlineSeconds)
			return wait.Poll(time.Second, deadline*time.Second, func() (done bool, err error) {
				batch := &batchv1.Job{}
				err = r.client.Get(context.Background(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, batch)
				if err != nil {
					log.WithField("dependency", dep.Name).WithError(err).Error("failed to get the status of the dependency")
					return false, err
				}

				// for now, we just assume each batch job has one pod
				if batch.Status.Succeeded != 1 {
					log.WithField("dependency", dep.Name).Info("Waiting for dependency to complete")
					return false, nil
				}

				return true, nil
			})
		}
	*/
	return nil
}
