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

package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	modelv1 "github.com/apache/skywalking-banyandb/api/proto/banyandb/model/v1"
	"github.com/apache/skywalking-banyandb/banyand/measure/aggregate"
)

func TestMin(t *testing.T) {
	var err error

	// case1: input int64 values
	minInt64, _ := aggregate.NewFunction[int64, aggregate.Void, int64](modelv1.MeasureAggregate_MEASURE_AGGREGATE_MIN)
	err = minInt64.Combine(aggregate.NewMinArguments[int64](
		[]int64{1, 2, 3}, // mock the "minimum" column
	))
	assert.NoError(t, err)
	_, _, r1 := minInt64.Result()
	assert.Equal(t, int64(1), r1)

	// case2: input float64 values
	minFloat64, _ := aggregate.NewFunction[float64, aggregate.Void, float64](modelv1.MeasureAggregate_MEASURE_AGGREGATE_MIN)
	err = minFloat64.Combine(aggregate.NewMinArguments[float64](
		[]float64{1.0, 2.0, 3.0}, // mock the "minimum" column
	))
	assert.NoError(t, err)
	_, _, r2 := minFloat64.Result()
	assert.Equal(t, 1.0, r2)
}
