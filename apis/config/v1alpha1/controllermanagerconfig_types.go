package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cfg "sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
)

// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.

// ControllerManagerConfigSpec defines the desired state of ControllerManagerConfig.
type ControllerManagerConfigSpec struct {
	// ControllerManagerConfigurationSpec returns the contfigurations for controllers.
	cfg.ControllerManagerConfigurationSpec `json:",inline"`

	// NamespaceSelector supports label selecting to include / exclude namespaces
	// from the controller. Use case is to exclude specific namespaces e.g. for
	// development purposes.
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
}

// ControllerManagerConfigStatus defines the observed state of ControllerManagerConfig.
type ControllerManagerConfigStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ControllerManagerConfig is the Schema for the controllermanagerconfigs API.
//
// Currently all fields have to be defined although not really necessary for the
// use case of this object. See https://github.com/operator-framework/operator-sdk/issues/5584.
type ControllerManagerConfig struct {
	metav1.TypeMeta `json:",inline"`

	Spec   ControllerManagerConfigSpec   `json:"spec,omitempty"`
	Status ControllerManagerConfigStatus `json:"status,omitempty"`
}

func init() {
	SchemeBuilder.Register(&ControllerManagerConfig{})
}

// Complete returns the configuration for controller-runtime.
func (c *ControllerManagerConfig) Complete() (cfg.ControllerManagerConfigurationSpec, error) {
	return c.Spec.ControllerManagerConfigurationSpec, nil
}
