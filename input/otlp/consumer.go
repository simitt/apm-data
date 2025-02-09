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

package otlp

import (
	"sync/atomic"

	"github.com/elastic/apm-data/model"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

// ConsumerConfig holds configuration for Consumer.
type ConsumerConfig struct {
	// Processor holds the model.BatchProcessor which will be invoked
	// with event batches when consuming OTLP payloads.
	Processor model.BatchProcessor

	// Logger holds a logger for the consumer. If this is nil, then
	// no logging will be performed.
	Logger *zap.Logger
}

// Consumer transforms OpenTelemetry data to the Elastic APM data model,
// sending each payload as a batch to the configured BatchProcessor.
type Consumer struct {
	stats consumerStats

	config ConsumerConfig
}

// NewConsumer returns a new Consumer with the given configuration.
func NewConsumer(config ConsumerConfig) *Consumer {
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	} else {
		config.Logger = config.Logger.Named("otel")
	}
	return &Consumer{config: config}
}

// ConsumerStats holds a snapshot of statistics about data consumption.
type ConsumerStats struct {
	// UnsupportedMetricsDropped records the number of unsupported metrics
	// that have been dropped by the consumer.
	UnsupportedMetricsDropped int64
}

// consumerStats holds the current statistics, which must be accessed and
// modified using atomic operations.
type consumerStats struct {
	unsupportedMetricsDropped int64
}

// Stats returns a snapshot of the current statistics about data consumption.
func (c *Consumer) Stats() ConsumerStats {
	return ConsumerStats{
		UnsupportedMetricsDropped: atomic.LoadInt64(&c.stats.unsupportedMetricsDropped),
	}
}

// Capabilities is part of the consumer interfaces.
func (c *Consumer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: false,
	}
}
