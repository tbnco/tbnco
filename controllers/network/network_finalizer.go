package network

import (
	"context"

	networkv1alpha1 "github.com/tbnco/tbnco/apis/network/v1alpha1"
)

// NetworkFinalizer defines the finalizer string for the CR.
// Follows <kind>.<group>/finalizer naming scheme.
const NetworkFinalizer = "network.network.tbnco.github.io/finalizer"

// finalize runs the finalizer logic for Network.
func (r *NetworkReconciler) finalize(_ context.Context, cr *networkv1alpha1.Network) error {
	// TODO: Add finalizer logic.

	return nil
}
