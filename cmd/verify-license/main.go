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

package main

import (
	"fmt"

	"go.bytebuilders.dev/license-verifier/info"
)

func main() {
	data := `-----BEGIN CERTIFICATE-----
MIIEgzCCA2ugAwIBAgIIaSg9nV8JDpcwDQYJKoZIhvcNAQELBQAwJTEWMBQGA1UE
ChMNQXBwc0NvZGUgSW5jLjELMAkGA1UEAxMCY2EwHhcNMjUxMDA5MDYxNDMxWhcN
MjUxMTA4MDYxNDMxWjCCARgxDzANBgNVBAYTBmt1YmVkYjETMBEGA1UECBMKZW50
ZXJwcmlzZTGBpDAXBgNVBAoTEGt1YmVkYi1jb21tdW5pdHkwFwYDVQQKExBrdWJl
ZGItZXh0LXN0YXNoMBgGA1UEChMRa3ViZWRiLWF1dG9zY2FsZXIwGAYDVQQKExFr
dWJlZGItZW50ZXJwcmlzZTAcBgNVBAoTFXBhbm9wdGljb24tZW50ZXJwcmlzZTAe
BgNVBAoTF2t1YmVkYi1tb25pdG9yaW5nLWFnZW50MRowGAYDVQQLExFrdWJlZGIt
ZW50ZXJwcmlzZTEtMCsGA1UEAxMkOTc4Nzk0ZTctODYxMC00OWEwLWEzNTQtYjRm
NTE1YzY0NGE4MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuzISekbt
QhW0CqqZa7n2sPlXclCxGtLTgSVm/rKOmmZpXRyAyRtFxg1im76Y34+07af5fuIr
1EtAPPqh95tF+rFFqbvJ+XtBp1EKJ7rTvFmB6pzCstghC04fbtKemJDHNBv4l7kP
TCfe93wf3BaseNiwOvPz6r4/Mijqb3AgiTR0+5tzG2f2HuuQR5Y34Xhhb93JpMvg
nWqdz4B8Gf0P96r06/ZA8siLU+C+9qLalmzNBqu+pnolnW3b9XF6VDCXzp2aXAQs
L+ooYCHBsLg1EnnwAKaxyncpl6z+RlBuNv5JkDzxln59LcF4cUNii2+Hd9P9rDxf
w0HdrQSs/gCVOQIDAQABo4HBMIG+MA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAK
BggrBgEFBQcDAjAfBgNVHSMEGDAWgBTZMREkIF69G3qzMCLSm9DzksSgeTB2BgNV
HREEbzBtgiQ5Nzg3OTRlNy04NjEwLTQ5YTAtYTM1NC1iNGY1MTVjNjQ0YTiBLkhp
cmFubW95IERhcyBDaG93ZGh1cnkgPGhpcmFubW95QGFwcHNjb2RlLmNvbT6BFWhp
cmFubW95QGFwcHNjb2RlLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEAlvZZacpStUTS
mC+Slmc7hqGu2vDf1m5CeGHP7mOk1tD6sD/SJAM0t+qu+sQ2KsykAu2JWgs9Ck1/
UiMMH0lvRVVXPgfYl9wcLYwioYXmB+X5U6T040N/bFr9UkcusWVgZbiPZZdPbOq5
yQ1he9FvHXQ4GajQF+RmzcT37AUn1SP1ICN84MMn1x0vlWuJDxwscIyERbxa1GUH
XWbnhcQLk4uYENgPZA2gWkWF5sEmUkbkqhLu1QWqHGGIcgHc/YUVsOUgsv0/cq1A
WV2juHeslnz/EbqyJD8xJLYdX+xFMULlSJu7PQRuLOSlYzPKfv/KR0k7ZvRGXU/f
Ta3dpe/anQ==
-----END CERTIFICATE-----`
	cert, err := info.ParseCertificate([]byte(data))
	if err != nil {
		panic(err)
	}
	fmt.Println(cert.EmailAddresses)
}
