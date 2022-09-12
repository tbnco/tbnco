package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkConditionType defines the underlying type for conditions.
type NetworkConditionType string

// Definition of condition types.
// Follows https://github.com/openshift/custom-resource-status/tree/master/conditions
// recommendations at the moment.
const (
	// NetworkAvailable means the network and associated
	// infrastructure is configured and usable.
	NetworkAvailable NetworkConditionType = "Available"
	// NetworkProgressing means the network is progressing.
	// A network progresses, when its associated configuration
	// is not fully achieved. The controller actively
	// works on reconciling the required state.
	NetworkProgressing NetworkConditionType = "Progressing"
	// NetworkDegraded means the network is degraded.
	// Does not imply 'Available=false', e.g. when the network
	// in general is usable, but with limited performance.
	NetworkDegraded NetworkConditionType = "Degraded"
)

// StatusReason defines the underlying type for condition reasons.
type StatusReason string

const (
	ReasonInitializing     StatusReason = "Initializing"
	ReasonNetworkAvailable StatusReason = "NetworkAvailable"
	ReasonFinalizing       StatusReason = "Finalizing"
)

// SetConditionDefaults initializes all conditions.
// Usually called at the beginning of the reconcile loop to
// have all conditions in a defined state.
func (s *NetworkStatus) SetConditionDefaults() {
	const msg string = "Controller is initializing CR"

	if cond := meta.FindStatusCondition(s.Conditions, string(NetworkAvailable)); cond == nil {
		meta.SetStatusCondition(&s.Conditions, metav1.Condition{
			Type:    string(NetworkAvailable),
			Status:  metav1.ConditionFalse,
			Reason:  string(ReasonInitializing),
			Message: msg,
		})
	}

	if cond := meta.FindStatusCondition(s.Conditions, string(NetworkProgressing)); cond == nil {
		meta.SetStatusCondition(&s.Conditions, metav1.Condition{
			Type:    string(NetworkProgressing),
			Status:  metav1.ConditionTrue,
			Reason:  string(ReasonInitializing),
			Message: msg,
		})
	}

	if cond := meta.FindStatusCondition(s.Conditions, string(NetworkDegraded)); cond == nil {
		meta.SetStatusCondition(&s.Conditions, metav1.Condition{
			Type:    string(NetworkDegraded),
			Status:  metav1.ConditionFalse,
			Reason:  string(ReasonInitializing),
			Message: msg,
		})
	}
}

// SetConditionAvailable sets the available condition to true
// and progressing or degraded conditions to false.
// Meant to be used when network is fully usable.
func (s *NetworkStatus) SetConditionAvailable() {
	meta.SetStatusCondition(&s.Conditions, metav1.Condition{
		Type:    string(NetworkAvailable),
		Status:  metav1.ConditionTrue,
		Reason:  string(ReasonNetworkAvailable),
		Message: "Network configured and ready to use",
	})

	meta.SetStatusCondition(&s.Conditions, metav1.Condition{
		Type:    string(NetworkProgressing),
		Status:  metav1.ConditionFalse,
		Reason:  string(ReasonNetworkAvailable),
		Message: "",
	})

	meta.SetStatusCondition(&s.Conditions, metav1.Condition{
		Type:    string(NetworkDegraded),
		Status:  metav1.ConditionFalse,
		Reason:  string(ReasonNetworkAvailable),
		Message: "",
	})
}

// SetConditionProgressing sets the progressing condition.
// Does not modify any other condition.
func (s *NetworkStatus) SetConditionProgressing(reason StatusReason, message string) {
	meta.SetStatusCondition(&s.Conditions, metav1.Condition{
		Type:    string(NetworkProgressing),
		Status:  metav1.ConditionTrue,
		Reason:  string(reason),
		Message: message,
	})
}

// SetConditionDegraded sets the degraded condition.
// Does not modify any other condition.
func (s *NetworkStatus) SetConditionDegraded(reason StatusReason, message string) {
	meta.SetStatusCondition(&s.Conditions, metav1.Condition{
		Type:    string(NetworkDegraded),
		Status:  metav1.ConditionTrue,
		Reason:  string(reason),
		Message: message,
	})
}
