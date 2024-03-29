/*
Copyright 2022 The KodeRover Authors.

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

package types

import "fmt"

const ZadigDomain = "zadig.koderover.com"
const Zadig = "zadig"

var ZadigLabelKeyGlobalOwner = fmt.Sprintf("%s/owner", ZadigDomain)

const IstioLabelKeyInjection = "istio-injection"
const IstioLabelValueInjection = "enabled"

var OriginSpec = fmt.Sprintf("%s/origin", ZadigDomain)

const (
	ZadigReleaseVersionLabelKey     = "zadigx-release-version"
	ZadigReleaseTypeLabelKey        = "zadigx-release-type"
	ZadigReleaseServiceNameLabelKey = "zadigx-release-service-name"
	ZadigReleaseMSEGrayTagLabelKey  = "alicloud.service.tag"
)

const (
	ZadigReleaseVersionOriginal = "original"

	ZadigReleaseTypeMseGray   = "mse-gray"
	ZadigReleaseTypeBlueGreen = "blue-green"
)

var ZadigReleaseTypeList = []string{
	ZadigReleaseTypeMseGray,
	ZadigReleaseTypeBlueGreen,
}
