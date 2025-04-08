package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	loremipsum "gopkg.in/loremipsum.v1"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var fakeSpanDuration = time.Duration(50 * time.Millisecond)
var randSource *rand.Rand

var loremIpsumGenerator = loremipsum.New()
var logCount = 0

func init() {

	logCount, _ = strconv.Atoi(os.Getenv("LOG_COUNT"))
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	randSource = rand.New(rand.NewSource(rngSeed))

}

func GenerateData(typ string, key string) string {
	switch typ {
	case "logs":
		w1 := loremIpsumGenerator.Word()
		w2 := loremIpsumGenerator.Word()
		w3 := loremIpsumGenerator.Word()
		tim := strconv.FormatInt(time.Now().UnixNano(), 10)
		words := loremIpsumGenerator.Words(15)
		sid1 := trace.SpanID{}
		_, _ = randSource.Read(sid1[:])

		tid := trace.TraceID{}
		_, _ = randSource.Read(tid[:])

		str := strings.Builder{}
		str.Write([]byte(`{
  "resourceLogs": [
    {
      "resource": {
        "attributes": [
          {"key": "host.id", "value": { "stringValue": "` + w1 + `"  }},
          {"key": "host.name", "value": { "stringValue": "` + w1 + `"  }},
          {"key": "mw.account_key", "value": { "stringValue": "` + key + `"  }}
        ]
      },
      "scopeLogs": [ 

{
          "scope": {
            "name": "loadgen.library",
            "version": "1.0.0",
            "attributes": [
              {  "key": "my.scope.attribute", "value": { "stringValue": "some scope attribute" } }
            ]
          },
          "logRecords": [
`))

		for i := 1; i <= logCount; i++ {
			str.Write([]byte(`
            {
              "timeUnixNano": "` + tim + `",
              "observedTimeUnixNano": "` + tim + `",
              "traceId": "` + tid.String() + `",
              "spanId": "` + sid1.String() + `",
              "body": {
                "stringValue": "` + words + `"
              },
              "attributes": [
			 	{"key": "event.id", "value": { "stringValue": "` + strconv.Itoa(i) + `"  }},
			 	{"key": "event.message", "value": { "stringValue": "` + w1 + `"  }},
          		{"key": "event.domain", "value": { "stringValue": "` + w2 + `"  }},
          		{"key": "event.reason", "value": { "stringValue": "` + w3 + `"  }}
				]
			}
          `))
			if i != logCount {
				str.Write([]byte(","))
			}
		}
		str.Write([]byte(`]}  ]}  ]}`))

		///log.Printf("logs str %v", str.String())
		return str.String()

	case "traces":
		sid0 := trace.SpanID{}
		_, _ = randSource.Read(sid0[:])

		sid1 := trace.SpanID{}
		_, _ = randSource.Read(sid1[:])

		tid := trace.TraceID{}
		_, _ = randSource.Read(tid[:])

		return `
{
  "resourceSpans": [
    {
      "resource": {
        "attributes": [
          {
            "key": "service.name",
            "value": {
              "stringValue": "app-load-gen"
            }
          },
          {
            "key": "telemetry.sdk.language",
            "value": {
              "stringValue": "webjs"
            }
          },
          {
            "key": "telemetry.sdk.name",
            "value": {
              "stringValue": "opentelemetry"
            }
          },
          {
            "key": "telemetry.sdk.version",
            "value": {
              "stringValue": "1.19.0"
            }
          },
          {
            "key": "mw_agent",
            "value": {
              "boolValue": false
            }
          },
          {
            "key": "project.name",
            "value": {
              "stringValue": "app-load-gen"
            }
          },
          {
            "key": "mw.account_key",
            "value": {
              "stringValue": "` + key + `"
            }
          },
          {
            "key": "env",
            "value": {
              "stringValue": "test"
            }
          },
          {
            "key": "browser.trace",
            "value": {
              "boolValue": true
            }
          },
          {
            "key": "type",
            "value": {
              "stringValue": "rum"
            }
          },
          {
            "key": "origin",
            "value": {
              "stringValue": "https://load-test.io"
            }
          },
          {
            "key": "rum_origin",
            "value": {
              "stringValue": "https://load-test.io"
            }
          },
          {
            "key": "os",
            "value": {
              "stringValue": "Windows"
            }
          },
          {
            "key": "browser.name",
            "value": {
              "stringValue": "Chrome"
            }
          },
          {
            "key": "navigator.userAgent",
            "value": {
              "stringValue": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
            }
          },
          {
            "key": "root.url",
            "value": {
              "stringValue": "/test-api-call"
            }
          }
        ],
        "droppedAttributesCount": 0
      },
      "scopeSpans": [
        {
          "scope": {
            "name": "app-instrumentation-xml-http-request",
            "version": "0.1"
          },
          "spans": [
            {
              "traceId": "` + tid.String() + `",
              "spanId":"` + sid0.String() + `",
              "name": "POST https://load-test.io/api/v1/builder/widget/data?req=resource=trace-filter-label-by-service.name",
              "kind": 3,
              "startTimeUnixNano": "` + strconv.FormatInt(time.Now().UnixNano(), 10) + `",
              "endTimeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15).UnixNano(), 10) + `",
              "attributes": [
                {
                  "key": "http.method",
                  "value": {
                    "stringValue": "POST"
                  }
                },
                {
                  "key": "http.url",
                  "value": {
                    "stringValue": "https://load-test.io/api/v1/builder/widget/data?req=resource=trace-filter-label-by-service.name"
                  }
                },
                {
                  "key": "http.response_content_length",
                  "value": {
                    "intValue": 229
                  }
                },
                {
                  "key": "http.response_content_length_uncompressed",
                  "value": {
                    "intValue": 390
                  }
                },
                {
                  "key": "http.status_code",
                  "value": {
                    "intValue": 200
                  }
                },
                {
                  "key": "http.status_text",
                  "value": {
                    "stringValue": ""
                  }
                },
                {
                  "key": "http.host",
                  "value": {
                    "stringValue": "load-test.io"
                  }
                },
                {
                  "key": "http.scheme",
                  "value": {
                    "stringValue": "https"
                  }
                },
                {
                  "key": "http.user_agent",
                  "value": {
                    "stringValue": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
                  }
                },
                {
                  "key": "session.id",
                  "value": {
                    "stringValue": "e6f59f9ae1a0d2ea9270e34a5dcabc4e"
                  }
                },
                {
                  "key": "os",
                  "value": {
                    "stringValue": "Windows"
                  }
                },
                {
                  "key": "browser.name",
                  "value": {
                    "stringValue": "Chrome"
                  }
                },
                {
                  "key": "navigator.userAgent",
                  "value": {
                    "stringValue": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
                  }
                },
                {
                  "key": "env",
                  "value": {
                    "stringValue": "test"
                  }
                },
                {
                  "key": "root.url",
                  "value": {
                    "stringValue": "/apm/list"
                  }
                },
                {
                  "key": "recording.available",
                  "value": {
                    "stringValue": "true"
                  }
                },
                {
                  "key": "name",
                  "value": {
                    "stringValue": "ApiLoadUser"
                  }
                },
                {
                  "key": "email",
                  "value": {
                    "stringValue": "test@middleware.io"
                  }
                }
              ],
              "droppedAttributesCount": 0,
              "events": [
                {
                  "attributes": [],
                  "name": "open",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*1).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "send",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*2).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "fetchStart",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*3).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "domainLookupStart",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*4).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "domainLookupEnd",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*5).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "connectStart",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*6).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "secureConnectionStart",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*7).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "connectEnd",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*8).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "requestStart",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*9).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "responseStart",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*10).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "responseEnd",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*11).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "loaded",
                  "timeUnixNano":  "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*12).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                }
              ],
              "droppedEventsCount": 0,
              "status": {
                "code": 0
              },
              "links": [],
              "droppedLinksCount": 0
            },
            {
              "traceId": "` + tid.String() + `",
              "spanId": "` + sid1.String() + `",
              "name": "POST https://load-test.io/api/v1/builder/widget/data?req=resource=trace-list",
              "kind": 3,
              "startTimeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*14).UnixNano(), 10) + `",
              "endTimeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*1).UnixNano(), 10) + `",
              "attributes": [
                {
                  "key": "http.method",
                  "value": {
                    "stringValue": "POST"
                  }
                },
                {
                  "key": "http.url",
                  "value": {
                    "stringValue": "https://load-test.io/api/v1/builder/widget/data?req=resource=trace-list"
                  }
                },
                {
                  "key": "http.response_content_length",
                  "value": {
                    "intValue": 3259
                  }
                },
                {
                  "key": "http.response_content_length_uncompressed",
                  "value": {
                    "intValue": 18265
                  }
                },
                {
                  "key": "http.status_code",
                  "value": {
                    "intValue": 200
                  }
                },
                {
                  "key": "http.status_text",
                  "value": {
                    "stringValue": ""
                  }
                },
                {
                  "key": "http.host",
                  "value": {
                    "stringValue": "load-test.io"
                  }
                },
                {
                  "key": "http.scheme",
                  "value": {
                    "stringValue": "https"
                  }
                },
                {
                  "key": "http.user_agent",
                  "value": {
                    "stringValue": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
                  }
                },
                {
                  "key": "session.id",
                  "value": {
                    "stringValue": "e6f59f9ae1a0d2ea9270e34a5dcabc4e"
                  }
                },
                {
                  "key": "os",
                  "value": {
                    "stringValue": "Windows"
                  }
                },
                {
                  "key": "browser.name",
                  "value": {
                    "stringValue": "Chrome"
                  }
                },
                {
                  "key": "navigator.userAgent",
                  "value": {
                    "stringValue": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
                  }
                },
                {
                  "key": "env",
                  "value": {
                    "stringValue": "prod"
                  }
                },
                {
                  "key": "root.url",
                  "value": {
                    "stringValue": "/apm/list"
                  }
                },
                {
                  "key": "recording.available",
                  "value": {
                    "stringValue": "true"
                  }
                },
                {
                  "key": "name",
                  "value": {
                    "stringValue": "ApiLoadUser"
                  }
                },
                {
                  "key": "email",
                  "value": {
                    "stringValue": "test@middleware.io"
                  }
                }
              ],
              "droppedAttributesCount": 0,
              "events": [
                {
                  "attributes": [],
                  "name": "open",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*2).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "send",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*3).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "fetchStart",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*4).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "domainLookupStart",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*5).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "domainLookupEnd",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*6).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "connectStart",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*7).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "secureConnectionStart",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*8).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "connectEnd",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*9).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "requestStart",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*10).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "responseStart",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*11).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "responseEnd",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*12).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                },
                {
                  "attributes": [],
                  "name": "loaded",
                  "timeUnixNano": "` + strconv.FormatInt(time.Now().Add(fakeSpanDuration*15*13).UnixNano(), 10) + `",
                  "droppedAttributesCount": 0
                }
              ],
              "droppedEventsCount": 0,
              "status": {
                "code": 0
              },
              "links": [],
              "droppedLinksCount": 0
            }
          ]
        }
      ]
    }
  ]
}
}
`
	case "metrics":
	default:
		panic(fmt.Errorf("invalid data type"))
	}
	return ""
}
