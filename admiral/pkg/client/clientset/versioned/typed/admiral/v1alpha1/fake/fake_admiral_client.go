/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/istio-ecosystem/admiral/admiral/pkg/client/clientset/versioned/typed/admiral/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAdmiralV1alpha1 struct {
	*testing.Fake
}

func (c *FakeAdmiralV1alpha1) ClientConnectionConfigs(namespace string) v1alpha1.ClientConnectionConfigInterface {
	return &FakeClientConnectionConfigs{c, namespace}
}

func (c *FakeAdmiralV1alpha1) Dependencies(namespace string) v1alpha1.DependencyInterface {
	return &FakeDependencies{c, namespace}
}

func (c *FakeAdmiralV1alpha1) GlobalTrafficPolicies(namespace string) v1alpha1.GlobalTrafficPolicyInterface {
	return &FakeGlobalTrafficPolicies{c, namespace}
}

func (c *FakeAdmiralV1alpha1) OutlierDetections(namespace string) v1alpha1.OutlierDetectionInterface {
	return &FakeOutlierDetections{c, namespace}
}

func (c *FakeAdmiralV1alpha1) RoutingPolicies(namespace string) v1alpha1.RoutingPolicyInterface {
	return &FakeRoutingPolicies{c, namespace}
}

func (c *FakeAdmiralV1alpha1) TrafficConfigs(namespace string) v1alpha1.TrafficConfigInterface {
	return &FakeTrafficConfigs{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAdmiralV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
