/*
Copyright AppsCode Inc.

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

package verifier

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/reference"

	"github.com/appscode/go/log"
	"github.com/pkg/errors"
	"gomodules.xyz/sets"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	clientscheme "k8s.io/client-go/kubernetes/scheme"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/dynamic"
	"kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/tools/clusterid"
)

type LicenseOptions struct {
	clusterUID  string
	productName string
	caCert      []byte
	license     []byte
	config      *rest.Config
	k8sClient   kubernetes.Interface
	licenseFile string
}

func (opt *LicenseOptions) validateLicense() error {
	block, _ := pem.Decode(opt.license)
	if block == nil {
		// This probably is a JWT token, should be check for that when ready
		return errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.Wrap(err, "failed to parse certificate")
	}

	// First, create the set of root certificates. For this example we only
	// have one. It's also possible to omit this in order to use the
	// default root set of the current operating system.
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(opt.caCert)
	if !ok {
		return errors.New("failed to parse root certificate")
	}

	opts := x509.VerifyOptions{
		DNSName: opt.clusterUID,
		Roots:   roots,
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
	}
	if _, err := cert.Verify(opts); err != nil {
		return errors.Wrap(err, "failed to verify certificate")
	}
	if !sets.NewString(cert.Subject.Organization...).Has(opt.productName) {
		return fmt.Errorf("license was not issued for %s", opt.productName)
	}
	return nil
}

// VerifyLicensePeriodically periodically verifies whether the provided license is valid for the current cluster or not.
func VerifyLicensePeriodically(config *rest.Config, licenseFile string, stopCh <-chan struct{}) error {
	opt := &LicenseOptions{
		licenseFile: licenseFile,
		config:      config,
	}
	// Create Kubernetes client
	err := opt.createClients()
	if err != nil {
		return opt.handleLicenseVerificationFailure(err)
	}
	// Read cluster UID (UID of the "kube-system" namespace)
	err = opt.readClusterUID()
	if err != nil {
		return opt.handleLicenseVerificationFailure(err)
	}

	// Periodically verify license with 1 hour interval
	return wait.PollUntil(1*time.Hour, func() (done bool, err error) {
		log.Infof("Verifying license.......")
		// Read license from file
		err = opt.readLicenseFromFile()
		if err != nil {
			return false, opt.handleLicenseVerificationFailure(err)
		}
		// Validate license
		err = opt.validateLicense()
		if err != nil {
			return false, opt.handleLicenseVerificationFailure(err)
		}
		log.Infof("Successfully verified license!")
		return true, nil
	}, stopCh)
}

// VerifyLicense verifies whether the provided license is valid for the current cluster or not.
func VerifyLicense(config *rest.Config, licenseFile string) error {
	log.Infof("Verifying license.......")
	opt := &LicenseOptions{
		licenseFile: licenseFile,
		config:      config,
	}
	// Create Kubernetes client
	err := opt.createClients()
	if err != nil {
		return opt.handleLicenseVerificationFailure(err)
	}
	// Read cluster UID (UID of the "kube-system" namespace)
	err = opt.readClusterUID()
	if err != nil {
		return opt.handleLicenseVerificationFailure(err)
	}
	// Read license from file
	err = opt.readLicenseFromFile()
	if err != nil {
		return opt.handleLicenseVerificationFailure(err)
	}
	// Validate license
	err = opt.validateLicense()
	if err != nil {
		return opt.handleLicenseVerificationFailure(err)
	}
	log.Infof("Successfully verified license!")
	return nil
}

func (opt *LicenseOptions) createClients() (err error) {
	opt.k8sClient, err = kubernetes.NewForConfig(opt.config)
	return err
}

func (opt *LicenseOptions) readLicenseFromFile() (err error) {
	opt.license, err = ioutil.ReadFile(opt.licenseFile)
	return err
}

func (opt *LicenseOptions) readClusterUID() (err error) {
	opt.clusterUID, err = clusterid.ClusterUID(opt.k8sClient.CoreV1().Namespaces())
	return err
}

func (opt *LicenseOptions) handleLicenseVerificationFailure(licenseErr error) error {
	defer func() {
		// Send interrupt so that all go-routines shut-down gracefully
		//nolint:errcheck
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		// Exit the process so that the pod crash.
		os.Exit(1)
	}()

	// Log licenseInfo verification failure
	log.Errorln("Failed to verify license. Reason: ", licenseErr.Error())

	// Don't write event if not running inside a cluster
	if !meta.PossiblyInCluster() {
		return nil
	}

	// Read current pod name
	podName, err := os.Hostname()
	if err != nil {
		return err
	}
	// Read the namespace of current pod
	namespace := meta.Namespace()

	// Find the root parent of this pod
	parent, _, err := dynamic.DetectWorkload(
		context.TODO(),
		opt.config,
		core.SchemeGroupVersion.WithResource(core.ResourcePods.String()),
		namespace,
		podName,
	)
	if err != nil {
		return err
	}
	ref, err := reference.GetReference(clientscheme.Scheme, parent)
	if err != nil {
		return err
	}
	eventMeta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(parent.GetName(), "license"),
		Namespace: namespace,
	}
	// Create an event against the root parent specifying that the license verification failed
	_, _, err = core_util.CreateOrPatchEvent(context.TODO(), opt.k8sClient, eventMeta, func(in *core.Event) *core.Event {
		in.InvolvedObject = *ref
		in.Type = core.EventTypeWarning
		in.Source = core.EventSource{Component: EventSourceLicenseVerifier}
		in.Reason = EventReasonLicenseVerificationFailed
		in.Message = fmt.Sprintf("Failed to verify license. Reason: %s", err.Error())

		if in.FirstTimestamp.IsZero() {
			in.FirstTimestamp = metav1.Now()
		}
		in.LastTimestamp = metav1.Now()
		in.Count = in.Count + 1

		return in
	}, metav1.PatchOptions{})
	return err
}

const (
	EventSourceLicenseVerifier           = "License Verifier"
	EventReasonLicenseVerificationFailed = "License Verification Failed"
)
