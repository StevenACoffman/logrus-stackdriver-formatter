package logadapter

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/StevenACoffman/logrus-stackdriver-formatter/test"
	"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
)

func TestStackSkip(t *testing.T) {
	var out bytes.Buffer

	logger := logrus.New()
	logger.Out = &out
	logger.Formatter = NewFormatter(
		WithService("test"),
		WithVersion("0.1"),
		WithStackSkip("github.com/StevenACoffman/logrus-stackdriver-formatter"),
		WithSkipTimestamp(),
	)

	mylog := test.LogWrapper{
		Logger: logger,
	}
	mylog.Error("my log entry")

	want := map[string]interface{}{
		"severity": "ERROR",
		"message":  "my log entry",
		"logName":  "projects//logs/test",
		"trace":    "projects//traces/1",
		"serviceContext": map[string]interface{}{
			"service": "test",
			"version": "0.1",
		},
		"context": map[string]interface{}{
			"data": map[string]interface{}{
				"trace": "1",
			},
			"reportLocation": map[string]interface{}{
				"file":     "testing/testing.go",
				"line":     865.0,
				"function": "tRunner",
			},
		},
		"sourceLocation": map[string]interface{}{
			"file":     "testing/testing.go",
			"line":     865.0,
			"function": "tRunner",
		},
	}
	var got map[string]interface{}
	err := json.Unmarshal(out.Bytes(), &got)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"unexpected output = %# v; want = %# v",
			pretty.Formatter(got),
			pretty.Formatter(want))
	}
}
