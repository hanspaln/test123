package crdconversion

import (
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// serveIngressBackendsPolicyConversion servers endpoint for the converter defined as convertIngressBackendsPolicy function.
func serveIngressBackendsPolicyConversion(w http.ResponseWriter, r *http.Request) {
	serve(w, r, convertIngressBackendsPolicy)
}

// convertIngressBackendsPolicy contains the business logic to convert ingressbackends.policy.openservicemesh.io CRD
// Example implementation reference : https://github.com/kubernetes/kubernetes/blob/release-1.22/test/images/agnhost/crd-conversion-webhook/converter/example_converter.go
func convertIngressBackendsPolicy(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("IngressBackendsPolicy: conversion from a version to itself should not call the webhook: %s", toVersion)
	}

	log.Debug().Msg("IngressBackendsPolicy: successfully converted object")
	return convertedObject, statusSucceed()
}
