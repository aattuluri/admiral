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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/istio-ecosystem/admiral/admiral/pkg/apis/admiral/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RoutingPolicyLister helps list RoutingPolicies.
// All objects returned here must be treated as read-only.
type RoutingPolicyLister interface {
	// List lists all RoutingPolicies in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.RoutingPolicy, err error)
	// RoutingPolicies returns an object that can list and get RoutingPolicies.
	RoutingPolicies(namespace string) RoutingPolicyNamespaceLister
	RoutingPolicyListerExpansion
}

// routingPolicyLister implements the RoutingPolicyLister interface.
type routingPolicyLister struct {
	indexer cache.Indexer
}

// NewRoutingPolicyLister returns a new RoutingPolicyLister.
func NewRoutingPolicyLister(indexer cache.Indexer) RoutingPolicyLister {
	return &routingPolicyLister{indexer: indexer}
}

// List lists all RoutingPolicies in the indexer.
func (s *routingPolicyLister) List(selector labels.Selector) (ret []*v1alpha1.RoutingPolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.RoutingPolicy))
	})
	return ret, err
}

// RoutingPolicies returns an object that can list and get RoutingPolicies.
func (s *routingPolicyLister) RoutingPolicies(namespace string) RoutingPolicyNamespaceLister {
	return routingPolicyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RoutingPolicyNamespaceLister helps list and get RoutingPolicies.
// All objects returned here must be treated as read-only.
type RoutingPolicyNamespaceLister interface {
	// List lists all RoutingPolicies in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.RoutingPolicy, err error)
	// Get retrieves the RoutingPolicy from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.RoutingPolicy, error)
	RoutingPolicyNamespaceListerExpansion
}

// routingPolicyNamespaceLister implements the RoutingPolicyNamespaceLister
// interface.
type routingPolicyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all RoutingPolicies in the indexer for a given namespace.
func (s routingPolicyNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.RoutingPolicy, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.RoutingPolicy))
	})
	return ret, err
}

// Get retrieves the RoutingPolicy from the indexer for a given namespace and name.
func (s routingPolicyNamespaceLister) Get(name string) (*v1alpha1.RoutingPolicy, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("routingpolicy"), name)
	}
	return obj.(*v1alpha1.RoutingPolicy), nil
}
