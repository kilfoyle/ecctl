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

package depresource

import (
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_configuration_templates"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

const (
	// DefaultAppSearchRefID is used when the RefID is not specified.
	DefaultAppSearchRefID = "main-appsearch"
)

// NewAppSearch creates a *models.AppSearchPayload from the parameters.
// It relies on a simplified single dimension memory size and zone count to
// construct the AppSearch's topology.
func NewAppSearch(params NewStateless) (*models.AppSearchPayload, error) {
	params.fillDefaults(DefaultAppSearchRefID)
	if err := params.Validate(); err != nil {
		return nil, err
	}

	// When either not set, we obtain the values from the running deployment.
	// Overriding either or both is done at the end of the if.
	if err := getTemplateAndRefID(&params); err != nil {
		return nil, err
	}

	// Obtain the deployment template so we can create the appsearch topology from
	// the specified sizes. The sizing overrides are done in newAppSearchPayload.
	res, err := params.V1API.PlatformConfigurationTemplates.GetDeploymentTemplate(
		platform_configuration_templates.NewGetDeploymentTemplateParams().
			WithTemplateID(params.TemplateID).
			WithShowInstanceConfigurations(ec.Bool(true)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	if res.Payload.ClusterTemplate.Appsearch == nil {
		return nil, fmt.Errorf("deployment: the %s template is not configured for AppSearch. Please use another template if you wish to start AppSearch instances",
			params.TemplateID)
	}

	var clusterTopology = res.Payload.ClusterTemplate.Appsearch.Plan.ClusterTopology
	var topology = models.AppSearchTopologyElement{Size: new(models.TopologySize)}
	if len(clusterTopology) > 0 {
		topology = *clusterTopology[0]
	}
	var payload = newAppSearchPayload(params, topology)

	return &payload, nil
}

func newAppSearchPayload(params NewStateless, topology models.AppSearchTopologyElement) models.AppSearchPayload {
	if params.Size > 0 {
		topology.Size.Value = ec.Int32(params.Size)
	}
	if params.ZoneCount > 0 {
		topology.ZoneCount = params.ZoneCount
	}

	return models.AppSearchPayload{
		ElasticsearchClusterRefID: ec.String(params.ElasticsearchRefID),
		DisplayName:               params.Name,
		Region:                    ec.String(params.Region),
		RefID:                     ec.String(params.RefID),
		Plan: &models.AppSearchPlan{
			Appsearch:       &models.AppSearchConfiguration{Version: params.Version},
			ClusterTopology: []*models.AppSearchTopologyElement{&topology},
		},
	}
}
