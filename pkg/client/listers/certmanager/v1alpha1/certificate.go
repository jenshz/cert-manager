/*
Copyright 2019 The Jetstack cert-manager contributors.

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

package v1alpha1

import (
	v1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CertificateLister helps list Certificates.
type CertificateLister interface {
	// List lists all Certificates in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Certificate, err error)
	// Certificates returns an object that can list and get Certificates.
	Certificates(namespace string) CertificateNamespaceLister
	CertificateListerExpansion
}

// certificateLister implements the CertificateLister interface.
type certificateLister struct {
	indexer cache.Indexer
}

// NewCertificateLister returns a new CertificateLister.
func NewCertificateLister(indexer cache.Indexer) CertificateLister {
	return &certificateLister{indexer: indexer}
}

// List lists all Certificates in the indexer.
func (s *certificateLister) List(selector labels.Selector) (ret []*v1alpha1.Certificate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Certificate))
	})
	return ret, err
}

// Certificates returns an object that can list and get Certificates.
func (s *certificateLister) Certificates(namespace string) CertificateNamespaceLister {
	return certificateNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CertificateNamespaceLister helps list and get Certificates.
type CertificateNamespaceLister interface {
	// List lists all Certificates in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Certificate, err error)
	// Get retrieves the Certificate from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Certificate, error)
	CertificateNamespaceListerExpansion
}

// certificateNamespaceLister implements the CertificateNamespaceLister
// interface.
type certificateNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Certificates in the indexer for a given namespace.
func (s certificateNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Certificate, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Certificate))
	})
	return ret, err
}

// Get retrieves the Certificate from the indexer for a given namespace and name.
func (s certificateNamespaceLister) Get(name string) (*v1alpha1.Certificate, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("certificate"), name)
	}
	return obj.(*v1alpha1.Certificate), nil
}
