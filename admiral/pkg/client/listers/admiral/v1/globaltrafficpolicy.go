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

package v1

import (
	v1 "github.com/admiral/admiral/pkg/apis/admiral/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GlobalTrafficPolicyLister helps list GlobalTrafficPolicies.
type GlobalTrafficPolicyLister interface {
	// List lists all GlobalTrafficPolicies in the indexer.
	List(selector labels.Selector) (ret []*v1.GlobalTrafficPolicy, err error)
	// GlobalTrafficPolicies returns an object that can list and get GlobalTrafficPolicies.
	GlobalTrafficPolicies(namespace string) GlobalTrafficPolicyNamespaceLister
	GlobalTrafficPolicyListerExpansion
}

// globalTrafficPolicyLister implements the GlobalTrafficPolicyLister interface.
type globalTrafficPolicyLister struct {
	indexer cache.Indexer
}

// NewGlobalTrafficPolicyLister returns a new GlobalTrafficPolicyLister.
func NewGlobalTrafficPolicyLister(indexer cache.Indexer) GlobalTrafficPolicyLister {
	return &globalTrafficPolicyLister{indexer: indexer}
}

// List lists all GlobalTrafficPolicies in the indexer.
func (s *globalTrafficPolicyLister) List(selector labels.Selector) (ret []*v1.GlobalTrafficPolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.GlobalTrafficPolicy))
	})
	return ret, err
}

// GlobalTrafficPolicies returns an object that can list and get GlobalTrafficPolicies.
func (s *globalTrafficPolicyLister) GlobalTrafficPolicies(namespace string) GlobalTrafficPolicyNamespaceLister {
	return globalTrafficPolicyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GlobalTrafficPolicyNamespaceLister helps list and get GlobalTrafficPolicies.
type GlobalTrafficPolicyNamespaceLister interface {
	// List lists all GlobalTrafficPolicies in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.GlobalTrafficPolicy, err error)
	// Get retrieves the GlobalTrafficPolicy from the indexer for a given namespace and name.
	Get(name string) (*v1.GlobalTrafficPolicy, error)
	GlobalTrafficPolicyNamespaceListerExpansion
}

// globalTrafficPolicyNamespaceLister implements the GlobalTrafficPolicyNamespaceLister
// interface.
type globalTrafficPolicyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GlobalTrafficPolicies in the indexer for a given namespace.
func (s globalTrafficPolicyNamespaceLister) List(selector labels.Selector) (ret []*v1.GlobalTrafficPolicy, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.GlobalTrafficPolicy))
	})
	return ret, err
}

// Get retrieves the GlobalTrafficPolicy from the indexer for a given namespace and name.
func (s globalTrafficPolicyNamespaceLister) Get(name string) (*v1.GlobalTrafficPolicy, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("globaltrafficpolicy"), name)
	}
	return obj.(*v1.GlobalTrafficPolicy), nil
}
