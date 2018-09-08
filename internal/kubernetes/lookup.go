package kubernetes

import "strings"

type GroupFinder struct {
	namespace string
}

func NewKubernetesGroupFinder(ns string) *GroupFinder {
	return &GroupFinder{
		namespace: ns,
	}
}

func (g *GroupFinder) GroupFind(apiVersion, kind string) string {
	if apiVersion == "v1" {
		apiVersion = "core/v1"
	}

	apiVersion = strings.Replace(apiVersion, "/", ".", -1)

	return strings.Join([]string{g.namespace, apiVersion, kind}, ".")
}
