package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenMCPMigrationSpec defines the desired state of OpenMCPMigration
type OpenMCPMigrationSpec struct {
	// container service list spec
	MigrationServiceSources []MigrationServiceSource `json:"MigrationServiceSource"`
}
type MigrationServiceSource struct {
	// container service spec
	MigrationSources []MigrationSource `json:"MigrationSource"`
	VolumePath       string            `json:"VolumePath"`
	LinkShareStatus  bool              `json:"LinkShareStatus"`
}
type MigrationSource struct {
	// Migration source
	TargetCluster string `json:"TargetCluster"`
	SourceCluster string `json:"SourceCluster"`
	NameSpace     string `json:"NameSpace"`
	ResourceType  string `json:"ResourceType"`
	ResourceName  string `json:"ResourceName"`
}

// OpenMCPMigrationStatus defines the observed state of OpenMCPMigration
type OpenMCPMigrationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	MigrationStatus bool   `json:"MigrationSpec"`
	Reason          string `json:"Reason"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPMigration is the Schema for the openmcpmigrations API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=openmcpmigrations,scope=Namespaced
type OpenMCPMigration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenMCPMigrationSpec   `json:"spec,omitempty"`
	Status OpenMCPMigrationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OpenMCPMigrationList contains a list of OpenMCPMigration
type OpenMCPMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenMCPMigration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenMCPMigration{}, &OpenMCPMigrationList{})
}
