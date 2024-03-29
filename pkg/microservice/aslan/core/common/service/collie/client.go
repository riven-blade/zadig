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

package collie

import (
	"github.com/koderover/zadig/v2/pkg/tool/httpclient"
)

type Client struct {
	*httpclient.Client

	host string
}

func New(host string) *Client {
	c := httpclient.New(
		httpclient.SetHostURL(host),
	)

	return &Client{
		Client: c,
		host:   host,
	}
}
