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

package info

import (
	"net/url"
	"path"
	"strconv"
	"strings"
	"unicode"

	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	EnforceLicense string
	LicenseCA      string

	ProductOwnerName string
	ProductOwnerUID  string

	ProductName string // This has been renamed to Features
	ProductUID  string

	prodAddress          = "https://byte.builders"
	qaAddress            = "https://appscode.ninja"
	registrationAPIPath  = "api/v1/register"
	licenseIssuerAPIPath = "api/v1/license/issue"
)

func Features() []string {
	return ParseFeatures(ProductName)
}

func ParseFeatures(features string) []string {
	out := strings.FieldsFunc(features, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == ';'
	})
	return sets.NewString(out...).List()
}

func SkipLicenseVerification() bool {
	v, _ := strconv.ParseBool(EnforceLicense)
	return !v
}

func RegistrationAPIEndpoint() string {
	u := APIServerAddress()
	u.Path = path.Join(u.Path, registrationAPIPath)
	return u.String()
}

func LicenseIssuerAPIEndpoint() string {
	u := APIServerAddress()
	u.Path = path.Join(u.Path, licenseIssuerAPIPath)
	return u.String()
}

func APIServerAddress() *url.URL {
	if SkipLicenseVerification() {
		u, _ := url.Parse(qaAddress)
		return u
	}
	u, _ := url.Parse(prodAddress)
	return u
}
