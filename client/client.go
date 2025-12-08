/*
Copyright AppsCode Inc. and Contributors

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

package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"

	"go.bytebuilders.dev/license-verifier/apis/licenses"
	"go.bytebuilders.dev/license-verifier/apis/licenses/v1alpha1"
	"go.bytebuilders.dev/license-verifier/info"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"moul.io/http2curl/v2"
)

type Client struct {
	url        string
	token      string
	clusterUID string
	caCert     []byte
	client     *http.Client
	userAgent  string
}

func NewClient(baseURL, token, clusterUID string, caCert []byte, insecureSkipVerifyTLS bool, userAgent string) (*Client, error) {
	u, err := info.LicenseIssuerAPIEndpoint(baseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		url:        u,
		token:      token,
		clusterUID: clusterUID,
		client:     http.DefaultClient,
		userAgent:  userAgent,
	}
	if len(caCert) > 0 || insecureSkipVerifyTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: insecureSkipVerifyTLS,
		}
		if len(c.caCert) > 0 {
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		}
		c.client = &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	}
	return c, nil
}

func (c *Client) AcquireLicense(features []string) ([]byte, *v1alpha1.Contract, error) {
	opts := struct {
		Cluster  string   `json:"cluster"`
		Features []string `json:"features"`
	}{
		Cluster:  c.clusterUID,
		Features: features,
	}
	data, err := json.Marshal(opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	// add authorization header to the req
	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}
	if klog.V(8).Enabled() {
		command, _ := http2curl.GetCurlCommand(req)
		klog.V(8).Infoln(command.String())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		var ce *tls.CertificateVerificationError
		if errors.As(err, &ce) {
			klog.ErrorS(err, "UnverifiedCertificates")
			for _, cert := range ce.UnverifiedCertificates {
				klog.Errorln(string(encodeCertPEM(cert)))
			}
		}
		return nil, nil, err
	}
	defer resp.Body.Close() // nolint:errcheck

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, apierrors.NewGenericServerResponse(
			resp.StatusCode,
			http.MethodPost,
			schema.GroupResource{Group: licenses.GroupName, Resource: "License"},
			"",
			string(body),
			0,
			false,
		)
	}

	lc := struct {
		Contract *v1alpha1.Contract `json:"contract,omitempty"`
		License  []byte             `json:"license"`
	}{}
	err = json.Unmarshal(body, &lc)
	if err != nil {
		return nil, nil, err
	}
	return lc.License, lc.Contract, nil
}

func encodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}
