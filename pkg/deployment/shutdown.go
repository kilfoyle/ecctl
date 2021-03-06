// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package deployment

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// ShutdownParams is consumed by Shutdown.
type ShutdownParams struct {
	*api.API
	DeploymentID string

	SkipSnapshot bool
	Hide         bool
}

// Validate ensures the parameters are usable by Shutdown.
func (params ShutdownParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.DeploymentID) != 32 {
		merr = multierror.Append(merr, util.ErrDeploymentID)
	}

	return merr.ErrorOrNil()
}

// Shutdown shuts down a deployment and all of its associated resources. To
// shutdown individual deployment resources use the kind specific APIs.
func Shutdown(params ShutdownParams) (*models.DeploymentShutdownResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.ShutdownDeployment(
		deployments.NewShutdownDeploymentParams().
			WithDeploymentID(params.DeploymentID).
			WithSkipSnapshot(ec.Bool(params.SkipSnapshot)).
			WithHide(ec.Bool(params.Hide)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
