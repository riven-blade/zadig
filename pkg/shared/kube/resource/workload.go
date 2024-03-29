/*
Copyright 2021 The KodeRover Authors.

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

package resource

type Workload struct {
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Images   []ContainerImage `json:"images"`
	Pods     []*Pod           `json:"pods"`
	Replicas int32            `json:"replicas"`
	// ZadigXReleaseType represent the release type of workload created by zadigx when it is not empty
	// frontend should limit or allow some operations on these workloads
	ZadigXReleaseType string `json:"zadigx_release_type"`
	ZadigXReleaseTag  string `json:"zadigx_release_tag"`
}

type ContainerImage struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	ImageName string `json:"image_name"`
}
