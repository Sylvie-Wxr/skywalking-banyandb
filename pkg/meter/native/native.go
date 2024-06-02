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
package native

import (
	"context"

	commonv1 "github.com/apache/skywalking-banyandb/api/proto/banyandb/common/v1"
	databasev1 "github.com/apache/skywalking-banyandb/api/proto/banyandb/database/v1"
	"github.com/apache/skywalking-banyandb/banyand/metadata"
	"github.com/apache/skywalking-banyandb/pkg/logger"
	"github.com/apache/skywalking-banyandb/pkg/meter"
)

const (
	NativeObservabilityGroupName = "_monitoring"
	defaultTagFamily = "default"
)

var log = logger.GetLogger("observability", "metrics", "system")

type noopInstrument struct{}

func (noopInstrument) Inc(_ float64, _ ...string)     {}
func (noopInstrument) Set(_ float64, _ ...string)     {}
func (noopInstrument) Add(_ float64, _ ...string)     {}
func (noopInstrument) Observe(_ float64, _ ...string) {}
func (noopInstrument) Delete(_ ...string) bool        { return false }

type provider struct {
	metadata   metadata.Repo
}

func NewProvider(metadata metadata.Repo) meter.Provider {
	return &provider{
		metadata: metadata,
	}
}

func(p *provider) createMeasure(metric string, labels ...string) error {
	var tags []*databasev1.TagSpec
	for _, label := range labels {
		tag := &databasev1.TagSpec{
			Name: label,
			Type: databasev1.TagType_TAG_TYPE_STRING,
		}
		tags = append(tags, tag)
	}
	_, err := p.metadata.MeasureRegistry().CreateMeasure(context.Background(), &databasev1.Measure{
		Metadata: &commonv1.Metadata{
			Name:  metric,
			Group: NativeObservabilityGroupName,
		},
		Entity: &databasev1.Entity{
			TagNames: []string{},
		},
		TagFamilies: []*databasev1.TagFamilySpec{
			{
				Name: defaultTagFamily,
				Tags: tags,
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
func (p *provider) RegisterCounter(name string, labelNames ...string) {
	err := p.createMeasure(name, labelNames...)
	if err != nil {
		log.Error().Err(err).Msgf("Failure to createMeasure for RegisterCounter %s", name)
	}
}

// Gauge returns a native implementation of the Gauge interface.
func (p *provider) RegisterGauge(_ string, _ ...string) {}

// Histogram returns a native implementation of the Histogram interface.
func (p *provider) RegisterHistogram(_ string, _ meter.Buckets, _ ...string) {}

// Counter returns a native implementation of the Counter interface.
func (p *provider) GetCounter(name string, labelNames ...string) meter.Counter {
	return noopInstrument{}
}

// Gauge returns a native implementation of the Gauge interface.
func (p *provider) GetGauge(_ string, _ ...string) meter.Gauge {
	return noopInstrument{}
}

// Histogram returns a native implementation of the Histogram interface.
func (p *provider) GetHistogram(_ string, _ meter.Buckets, _ ...string) meter.Histogram {
	return noopInstrument{}
}