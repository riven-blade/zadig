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

package aslan

import (
	"net/http"
	"net/url"

	"github.com/koderover/zadig/v2/pkg/tool/httpclient"
)

func (c *Client) Overview(header http.Header, qs url.Values) ([]byte, error) {
	url := "/stat/dashboard/overview"
	res, err := c.Get(url, httpclient.SetHeadersFromHTTPHeader(header), httpclient.SetQueryParamsFromValues(qs))
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}
func (c *Client) Test(header http.Header, qs url.Values) ([]byte, error) {
	url := "/stat/dashboard/test"
	res, err := c.Get(url, httpclient.SetHeadersFromHTTPHeader(header), httpclient.SetQueryParamsFromValues(qs))
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}
func (c *Client) Deploy(header http.Header, qs url.Values) ([]byte, error) {
	url := "/stat/dashboard/deploy"
	res, err := c.Get(url, httpclient.SetHeadersFromHTTPHeader(header), httpclient.SetQueryParamsFromValues(qs))
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}
func (c *Client) Build(header http.Header, qs url.Values) ([]byte, error) {
	url := "/stat/dashboard/build"
	res, err := c.Get(url, httpclient.SetHeadersFromHTTPHeader(header), httpclient.SetQueryParamsFromValues(qs))
	if err != nil {
		return nil, err
	}

	return res.Body(), nil
}
