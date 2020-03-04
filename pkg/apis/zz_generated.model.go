// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by model-api-gen. DO NOT EDIT.

package apis

import (
	time "time"
)

// SAdminSharableVirtualResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SAdminSharableVirtualResourceBase.
type SAdminSharableVirtualResourceBase struct {
	SSharableVirtualResourceBase
	Records string `json:"records"`
}

// SDomainLevelResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SDomainLevelResourceBase.
type SDomainLevelResourceBase struct {
	SStandaloneResourceBase
	SDomainizedResourceBase
	DomainSrc string `json:"domain_src"`
}

// SDomainizedResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SDomainizedResourceBase.
type SDomainizedResourceBase struct {
	DomainId string `json:"domain_id"`
}

// SEnabledResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SEnabledResourceBase.
type SEnabledResourceBase struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// SEnabledStatusDomainLevelResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SEnabledStatusDomainLevelResourceBase.
type SEnabledStatusDomainLevelResourceBase struct {
	SStatusDomainLevelResourceBase
	SEnabledResourceBase
}

// SEnabledStatusStandaloneResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SEnabledStatusStandaloneResourceBase.
type SEnabledStatusStandaloneResourceBase struct {
	SStatusStandaloneResourceBase
	SEnabledResourceBase
}

// SExternalizedResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SExternalizedResourceBase.
type SExternalizedResourceBase struct {
	ExternalId string `json:"external_id"`
}

// SJointResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SJointResourceBase.
type SJointResourceBase struct {
	SResourceBase
	RowId int64 `json:"row_id"`
}

// SKeystoneCacheObject is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SKeystoneCacheObject.
type SKeystoneCacheObject struct {
	SStandaloneResourceBase
	DomainId  string    `json:"domain_id"`
	Domain    string    `json:"domain"`
	LastCheck time.Time `json:"last_check"`
}

// SProjectizedResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SProjectizedResourceBase.
type SProjectizedResourceBase struct {
	SDomainizedResourceBase
	ProjectId string `json:"tenant_id"`
}

// SResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SResourceBase.
type SResourceBase struct {
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdateVersion int       `json:"update_version"`
	DeletedAt     time.Time `json:"deleted_at"`
	Deleted       bool      `json:"deleted"`
}

// SScopedResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SScopedResourceBase.
type SScopedResourceBase struct {
	DomainId  string `json:"domain_id"`
	ProjectId string `json:"tenant_id"`
}

// SSharableBaseResource is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SSharableBaseResource.
type SSharableBaseResource struct {
	IsPublic bool `json:"is_public"`
}

// SSharableVirtualResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SSharableVirtualResourceBase.
type SSharableVirtualResourceBase struct {
	SVirtualResourceBase
	IsPublic    bool   `json:"is_public"`
	PublicScope string `json:"public_scope"`
}

// SSharedResource is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SSharedResource.
type SSharedResource struct {
	SResourceBase
	Id              int64  `json:"id"`
	ResourceType    string `json:"resource_type"`
	ResourceId      string `json:"resource_id"`
	OwnerProjectId  string `json:"owner_project_id"`
	TargetProjectId string `json:"target_project_id"`
}

// SStandaloneResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStandaloneResourceBase.
type SStandaloneResourceBase struct {
	SResourceBase
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsEmulated  bool   `json:"is_emulated"`
}

// SStatusDomainLevelResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStatusDomainLevelResourceBase.
type SStatusDomainLevelResourceBase struct {
	SDomainLevelResourceBase
	SStatusResourceBase
}

// SStatusResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStatusResourceBase.
type SStatusResourceBase struct {
	Status string `json:"status"`
}

// SStatusStandaloneResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SStatusStandaloneResourceBase.
type SStatusStandaloneResourceBase struct {
	SStandaloneResourceBase
	SStatusResourceBase
}

// SVirtualJointResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SVirtualJointResourceBase.
type SVirtualJointResourceBase struct {
	SJointResourceBase
}

// SVirtualResourceBase is an autogenerated struct via yunion.io/x/onecloud/pkg/cloudcommon/db.SVirtualResourceBase.
type SVirtualResourceBase struct {
	SStatusStandaloneResourceBase
	SProjectizedResourceBase
	ProjectSrc       string    `json:"project_src"`
	IsSystem         bool      `json:"is_system"`
	PendingDeletedAt time.Time `json:"pending_deleted_at"`
	PendingDeleted   bool      `json:"pending_deleted"`
}
