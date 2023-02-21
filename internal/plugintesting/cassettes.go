package plugintesting

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// getTestCassettePath returns a cassette path for the current test
func getTestCassettePath(t *testing.T) string {
	specialChars := regexp.MustCompile(`[\\?%*:|"<>. ]`)

	// Replace nested tests separators.
	fileName := strings.ReplaceAll(t.Name(), "/", "-")

	fileName = strcase.ToBashArg(fileName)

	// Replace special characters.
	fileName = specialChars.ReplaceAllLiteralString(fileName, "")

	return filepath.Join(".", "testdata", fileName)
}

// getHTTPRecorder return a http client with a recorder for requests
// cleanup function should be called for the requests to be written to the test file
func getHTTPRecorder(t *testing.T, update bool) (client *http.Client, cleanup func() error, err error) {
	t.Helper()
	recorderMode := recorder.ModeReplaying
	if update {
		recorderMode = recorder.ModeRecording
	}

	r, err := recorder.NewAsMode(getTestCassettePath(t), recorderMode, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add a filter which removes Authorization headers from all requests:
	r.AddFilter(func(i *cassette.Interaction) error {
		i.Request.Headers = i.Request.Headers.Clone()
		delete(i.Request.Headers, "x-auth-token")
		delete(i.Request.Headers, "X-Auth-Token")
		delete(i.Request.Headers, "Authorization")
		return nil
	})

	return &http.Client{Transport: r}, func() error {
		err := r.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop http recorder: %w", err)
		}
		return nil
	}, nil
}
