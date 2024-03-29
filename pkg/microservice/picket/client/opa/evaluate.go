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

package opa

import (
	"fmt"
	"strings"

	"github.com/koderover/zadig/v2/pkg/tool/httpclient"
)

type InputGenerator func() (*Input, error)

// Evaluate evaluates the query with the given input and return a json response which has a field called "result"
func (c *Client) Evaluate(query string, result interface{}, ig InputGenerator) error {
	input, err := ig()
	if err != nil {
		return err
	}

	url := "v1/data"
	req := struct {
		Input interface{} `json:"input"`
	}{
		Input: input,
	}

	queryURL := fmt.Sprintf("%s/%s", url, strings.ReplaceAll(query, ".", "/"))
	_, err = c.Post(queryURL, httpclient.SetBody(req), httpclient.SetResult(result))
	if err != nil {
		return err
	}

	return nil
}
