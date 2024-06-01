// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
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

// Package meter provides a simple meter system for metrics. The metrics are aggregated by the meter provider.
package meter

import (
	"context"

	commonv1 "github.com/apache/skywalking-banyandb/api/proto/banyandb/common/v1"
	databasev1 "github.com/apache/skywalking-banyandb/api/proto/banyandb/database/v1"
	"github.com/apache/skywalking-banyandb/banyand/metadata"
)

type noopInstrument struct{}

func (noopInstrument) Inc(_ float64, _ ...string)     {}
func (noopInstrument) Set(_ float64, _ ...string)     {}
func (noopInstrument) Add(_ float64, _ ...string)     {}
func (noopInstrument) Observe(_ float64, _ ...string) {}
func (noopInstrument) Delete(_ ...string) bool        { return false }

// NativeProvider is native implementation of the Provider interface.
type provider struct {
	metadata   metadata.Repo
}

func NewProvider(metadata metadata.Repo) Provider {
	return &provider{
		metadata: metadata,
	}
}

func(p *provider) createMeasure(metric string, labels ...string) error {
	_, err := p.metadata.MeasureRegistry().CreateMeasure(context.Background(), &databasev1.Measure{
		Metadata: &commonv1.Metadata{
			Name:  metric,
			Group: "_monitoring",
		},
		Entity: &databasev1.Entity{
			TagNames: []string{"testtag"},
		},
		TagFamilies: []*databasev1.TagFamilySpec{
			{
				Name: metric,
				Tags: []*databasev1.TagSpec{
					{
						Name: "testtag",
						Type: databasev1.TagType_TAG_TYPE_STRING,
					},
				},
			},
		},
		Fields: []*databasev1.FieldSpec{
			{
				Name:              "testfield",
				FieldType:         databasev1.FieldType_FIELD_TYPE_INT,
				EncodingMethod:    databasev1.EncodingMethod_ENCODING_METHOD_GORILLA,
				CompressionMethod: databasev1.CompressionMethod_COMPRESSION_METHOD_ZSTD,
			},
		},
	})
	return err
}

// Counter returns a native implementation of the Counter interface.
func (p *provider) Counter(name string, labelNames ...string) Counter {
	return noopInstrument{}
}

// Gauge returns a native implementation of the Gauge interface.
func (p *provider) Gauge(_ string, _ ...string) Gauge {
	return noopInstrument{}
}

// Histogram returns a native implementation of the Histogram interface.
func (p *provider) Histogram(_ string, _ Buckets, _ ...string) Histogram {
	return noopInstrument{}
}
