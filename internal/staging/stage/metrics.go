// Copyright 2023 The Cockroach Authors
//
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
//
// SPDX-License-Identifier: Apache-2.0

package stage

import (
	"github.com/cockroachdb/cdc-sink/internal/util/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	stageRetireDurations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "stage_retire_duration_seconds",
		Help:    "the length of time it took to successfully retire applied mutations",
		Buckets: metrics.LatencyBuckets,
	}, metrics.TableLabels)
	stageRetireErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "stage_retire_errors_total",
		Help: "the number of times an error was encountered while retiring mutations",
	}, metrics.TableLabels)

	stageSelectCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "stage_select_mutations_total",
		Help: "the number of mutations read for this table",
	}, metrics.TableLabels)
	stageSelectDurations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "stage_select_duration_seconds",
		Help:    "the length of time it took to successfully select mutations",
		Buckets: metrics.LatencyBuckets,
	}, metrics.TableLabels)
	stageSelectErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "stage_select_errors_total",
		Help: "the number of times an error was encountered while selecting mutations",
	}, metrics.TableLabels)

	stageStoreCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "stage_store_mutations_total",
		Help: "the number of mutations stored for this table",
	}, metrics.TableLabels)
	stageStoreDurations = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "stage_store_duration_seconds",
		Help:    "the length of time it took to successfully store mutations",
		Buckets: metrics.LatencyBuckets,
	}, metrics.TableLabels)
	stageStoreErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "stage_store_errors_total",
		Help: "the number of times an error was encountered while storing mutations",
	}, metrics.TableLabels)
)
