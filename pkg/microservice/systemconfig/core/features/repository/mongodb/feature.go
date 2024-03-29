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

package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/koderover/zadig/v2/pkg/microservice/systemconfig/config"
	"github.com/koderover/zadig/v2/pkg/microservice/systemconfig/core/features/repository/models"
	mongotool "github.com/koderover/zadig/v2/pkg/tool/mongo"
)

type FeatureColl struct {
	*mongo.Collection

	coll string
}

func NewFeatureColl() *FeatureColl {
	name := models.Feature{}.TableName()
	coll := &FeatureColl{Collection: mongotool.Database(config.MongoDatabase()).Collection(name), coll: name}

	return coll
}

func (c *FeatureColl) GetCollectionName() string {
	return c.coll
}

func (c *FeatureColl) UpdateOrCreateFeature(ctx context.Context, obj *models.Feature) error {
	q := bson.M{"name": obj.Name}
	opts := options.Replace().SetUpsert(true)
	_, err := c.ReplaceOne(ctx, q, obj, opts)
	return err
}

func (c *FeatureColl) ListFeatures() ([]*models.Feature, error) {
	var fs []*models.Feature
	cursor, err := c.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &fs); err != nil {
		return nil, err
	}
	return fs, nil
}
