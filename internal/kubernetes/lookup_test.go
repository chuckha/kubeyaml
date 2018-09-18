package kubernetes_test

import (
	"testing"

	"github.com/chuckha/kube-validate/internal/kubernetes"
)

func TestAPIKey(t *testing.T) {
	testcases := []struct {
		name      string
		namespace string
		suffix    string
		version   string
		kind      string
		expected  string
	}{
		{
			name:      "simple test",
			namespace: "simple.test",
			suffix:    "hello",
			version:   "v1alpha1",
			kind:      "Deployment",
			expected:  "simple.test.v1alpha1.Deployment",
		},
		{
			name:      "v1 test",
			namespace: "k8s.io",
			suffix:    "hello",
			version:   "v1",
			kind:      "Pod",
			expected:  "k8s.io.core.v1.Pod",
		},
		{
			name:      "suffix test",
			namespace: "io.k8s.api",
			suffix:    ".k8s.io",
			version:   "certificates.k8s.io/v1beta1",
			kind:      "CertificateSigningRequest",
			expected:  "io.k8s.api.certificates.v1beta1.CertificateSigningRequest",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			apikeyer := kubernetes.NewAPIKeyer(tc.namespace, tc.suffix)
			key := apikeyer.APIKey(tc.version, tc.kind)
			if key != tc.expected {
				t.Fatalf("expected %v found %v", tc.expected, key)
			}
		})
	}
}
