/*
Copyright 2023 The KodeRover Authors.

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

package models

type UserGroup struct {
	Model
	GroupID     string `gorm:"column:group_id"    json:"group_id"`
	GroupName   string `gorm:"column:group_name"  json:"group_name"`
	Description string `gorm:"column:description" json:"description"`
	Type        int64  `gorm:"column:type"        json:"type"`
	// used to mention the foreign key relationship between userGroup and groupBinding
	// and specify the onDelete action.
	GroupBindings     []GroupBinding     `gorm:"foreignKey:GroupID;references:GroupID;constraint:OnDelete:CASCADE;" json:"-"`
	GroupRoleBindings []GroupRoleBinding `gorm:"foreignKey:GroupID;references:GroupID;constraint:OnDelete:CASCADE;" json:"-"`
}

// TableName sets the insert table name for this struct type
func (UserGroup) TableName() string {
	return "user_group"
}
