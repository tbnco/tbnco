package predicate

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// GetNamespaceFunc defines a getter function to get a
// namespace object.
type GetNamespaceFunc func(namespace string) (*corev1.Namespace, error)

// FilterNamespace sets up an event filter to select
// namespaces processed by the controller.
func FilterNamespace(selector labels.Selector, f GetNamespaceFunc) predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			ns, err := f(e.Object.GetNamespace())
			if err != nil {
				return false
			}

			labelsAslabels := labels.Set(ns.GetLabels())
			return selector.Matches(labelsAslabels)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			ns, err := f(e.ObjectNew.GetNamespace())
			if err != nil {
				return false
			}

			labelsAslabels := labels.Set(ns.GetLabels())
			return selector.Matches(labelsAslabels)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			ns, err := f(e.Object.GetNamespace())
			if err != nil {
				return false
			}

			labelsAslabels := labels.Set(ns.GetLabels())
			return selector.Matches(labelsAslabels)
		},
		GenericFunc: func(e event.GenericEvent) bool {
			ns, err := f(e.Object.GetNamespace())
			if err != nil {
				return false
			}

			labelsAslabels := labels.Set(ns.GetLabels())
			return selector.Matches(labelsAslabels)
		},
	}
}
