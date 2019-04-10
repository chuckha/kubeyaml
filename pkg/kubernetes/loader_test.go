package kubernetes_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/chuckha/kubeyaml/pkg/kubernetes"
)

type badReader struct{}

func (b *badReader) Read(p []byte) (int, error) {
	return 0, errors.New("failed to read")
}

func TestLoader(t *testing.T) {
	testcases := []struct {
		name  string
		input io.Reader
		check func(i *kubernetes.Input, err error, t *testing.T)
	}{
		{
			name:  "no apiVersion",
			input: strings.NewReader(`kind: Deployment`),
			check: func(i *kubernetes.Input, err error, t *testing.T) {
				if err == nil {
					t.Fatalf("expected an error but got nil")
				}
			},
		},
		{
			name:  "bad apiVersion",
			input: strings.NewReader(`apiVersion: 999`),
			check: func(i *kubernetes.Input, err error, t *testing.T) {
				if err == nil {
					t.Fatalf("expected an error but got nil")
				}
			},
		},
		{
			name:  "reader blows up",
			input: &badReader{},
			check: func(i *kubernetes.Input, err error, t *testing.T) {
				if err == nil {
					t.Fatalf("expected an error but got nil")
				}
			},
		},
		{
			name: "bad yaml should throw an error",
			input: strings.NewReader("	a tab"),
			check: func(i *kubernetes.Input, err error, t *testing.T) {
				if err == nil {
					t.Errorf("expected an error but did not get one")
				}
				if i != nil {
					t.Errorf("expected a nil Input but found something: %v", i)
				}
			},
		},
		{
			name: "easy test",
			input: strings.NewReader(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80`),
			check: func(i *kubernetes.Input, err error, t *testing.T) {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
			},
		},
	}

	l := kubernetes.NewLoader()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			i, err := l.Load(tc.input)
			tc.check(i, err, t)
		})
	}
}
