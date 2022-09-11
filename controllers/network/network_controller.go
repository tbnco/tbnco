package network

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	networkv1alpha1 "github.com/tbnco/tbnco/apis/network/v1alpha1"
)

// NetworkReconciler reconciles a Network object.
type NetworkReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=network.tbnco.github.io,resources=networks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=network.tbnco.github.io,resources=networks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=network.tbnco.github.io,resources=networks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *NetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch Network instance.
	network := &networkv1alpha1.Network{}
	if err := r.Get(ctx, req.NamespacedName, network); err != nil {
		// Object not found error occurs if the object is deleted.
		// No need to reconcile here again.
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		log.Error(err, "Failed to get Network instance")
		return ctrl.Result{}, err
	}

	// Initialize conditions if necessary.
	network.Status.SetConditionDefaults()
	if err := r.Status().Update(ctx, network); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	// Set finalizer if not already set.
	if !controllerutil.ContainsFinalizer(network, NetworkFinalizer) {
		log.Info("Add finalizer")

		if ok := controllerutil.AddFinalizer(network, NetworkFinalizer); !ok {
			log.Error(nil, "Failed to add finalizer")
			return ctrl.Result{Requeue: true}, nil
		}

		if err := r.Update(ctx, network); err != nil {
			log.Error(err, "Failed to update CR to add finalizer")
			return ctrl.Result{}, err
		}
	}

	// Check for CR deletion and run finalizer logic.
	isMarkedToBeDeleted := network.GetDeletionTimestamp() != nil
	if isMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(network, NetworkFinalizer) {
			log.Info("Run finalizer")

			// Set CR to degraded, as deletion is in progress and finalizer logic runs.
			network.Status.SetConditionDegraded(networkv1alpha1.ReasonFinalizing, "Finalizer running")
			if err := r.Status().Update(ctx, network); err != nil {
				log.Error(err, "Failed to update status")
				return ctrl.Result{}, err
			}

			// Run logic for finalizer. In case of errors
			// don't remove the finalizer so that it can
			// be retried at the next reconciliation loop.
			if err := r.finalize(ctx, network); err != nil {
				log.Error(err, "Failed to run finalizer")
				return ctrl.Result{}, err
			}

			// Set CR to degraded, as finalizer logic executed successfully.
			network.Status.SetConditionDegraded(networkv1alpha1.ReasonFinalizing, "Finalizer successfully executed")
			if err := r.Status().Update(ctx, network); err != nil {
				log.Error(err, "Failed to update status")
				return ctrl.Result{}, err
			}

			// Remove finalizer. Once all finalizers have been
			// removed, the object will be deleted.
			if ok := controllerutil.RemoveFinalizer(network, NetworkFinalizer); !ok {
				log.Error(nil, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, nil
			}

			// Update CR to apply finalizer deletion.
			if err := r.Update(ctx, network); err != nil {
				log.Error(err, "Failed to update resource to remove finalizer")
				return ctrl.Result{}, err
			}
		}

		// Nothing to do, end reconcile loop.
		return ctrl.Result{}, nil
	}

	// TODO: Run reconciler logic.

	// Reconciler loop successful. For now it is assumed, that
	// the Network is fully working.
	network.Status.SetConditionAvailable()
	if err := r.Status().Update(ctx, network); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetworkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1alpha1.Network{}).
		Complete(r)
}
