// Copyright Splunk Inc.
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

package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultRealm = "us0"
const defaultIngestURL = ""
const defaultToken = ""
const defaultFastIngest = true
const defaultReportingDuration = time.Duration(15) * time.Second
const defaultReportingTimeout = time.Duration(5) * time.Second
const defaultVerbose = false
const defaultHttpTracing = false

const ingestUrlFormat = "https://ingest.%s.signalfx.com"

const minTokenLength = 10 // SFx Access Tokens are 22 chars long in 2019 but accept 10 or more chars just in case

const realmEnv = "SPLUNK_REALM"
const ingestURLEnv = "SPLUNK_INGEST_URL"
const tokenEnv = "SPLUNK_ACCESS_TOKEN"
const fastIngestEnv = "FAST_INGEST"
const reportingDelayEnv = "REPORTING_RATE"
const reportingTimeoutEnv = "REPORTING_TIMEOUT"
const verboseEnv = "VERBOSE"
const httpTracingEnv = "HTTP_TRACING"

type Configuration struct {
	SplunkRealm      string
	SplunkMetricsUrl string
	SplunkToken      string
	FastIngest       bool
	ReportingDelay   time.Duration
	ReportingTimeout time.Duration
	Verbose          bool
	HttpTracing      bool
}

func New() Configuration {
	configuration := Configuration{
		SplunkRealm:      strOrDefault(realmEnv, defaultRealm),
		SplunkMetricsUrl: strOrDefault(ingestURLEnv, defaultIngestURL),
		SplunkToken:      strOrDefault(tokenEnv, defaultToken),
		FastIngest:       boolOrDefault(fastIngestEnv, defaultFastIngest),
		ReportingDelay:   durationOrDefault(reportingDelayEnv, defaultReportingDuration),
		ReportingTimeout: durationOrDefault(reportingTimeoutEnv, defaultReportingTimeout),
		Verbose:          boolOrDefault(verboseEnv, defaultVerbose),
		HttpTracing:      boolOrDefault(httpTracingEnv, defaultHttpTracing),
	}

	if configuration.SplunkMetricsUrl == "" {
		configuration.SplunkMetricsUrl = fmt.Sprintf(ingestUrlFormat, configuration.SplunkRealm)
	}

	configuration.SplunkMetricsUrl += "/v2/datapoint"

	return configuration
}

func (c Configuration) String() string {
	builder := strings.Builder{}
	addLine := func(format string, arg interface{}) { builder.WriteString(fmt.Sprintf(format+"\n", arg)) }

	addLine("Splunk Realm       = %v", c.SplunkRealm)
	addLine("Splunk Metrics URL = %v", c.SplunkMetricsUrl)
	addLine("Splunk Token       = %v", obfuscatedToken(c.SplunkToken))
	addLine("Fast Ingest        = %v", c.FastIngest)
	addLine("Reporting Delay    = %v", c.ReportingDelay.Seconds())
	addLine("Reporting Timeout  = %v", c.ReportingTimeout.Seconds())
	addLine("Verbose            = %v", c.Verbose)
	addLine("HTTP Tracing       = %v", c.HttpTracing)

	return builder.String()
}

func obfuscatedToken(token string) string {
	if len(token) < minTokenLength {
		return fmt.Sprintf("<invalid token> minimum %v chars required", minTokenLength)
	}
	return fmt.Sprintf("%s...%s", token[0:2], token[len(token)-2:])
}

func strOrDefault(key, d string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return d
}

func durationOrDefault(key string, d time.Duration) time.Duration {
	str := strOrDefault(key, "nan")

	if seconds, err := strconv.Atoi(str); err == nil {
		return time.Second * time.Duration(seconds)
	}

	log.Printf("can't parse number of seconds: %s\n", str)
	return d
}

func boolOrDefault(key string, d bool) bool {
	str := strOrDefault(key, "")

	if trueOrFalse, err := strconv.ParseBool(str); err == nil {
		return trueOrFalse
	}

	log.Printf("can't parse bool: %s\n", str)
	return d
}
