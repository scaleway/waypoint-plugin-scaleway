package plugintesting

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

var UpdateCassettes = os.Getenv("WAYPOINT_UPDATE_CASSETTES") == "true"

type cleanUpFunction func() error

type TestTools struct {
	t                *testing.T
	HttpClient       *http.Client
	scwClient        *scw.Client
	cleanupFunctions []cleanUpFunction
	cleanupRecorder  func() error
}

// Init new TestTools for a test
// TestTools.HttpClient should be used for all requests during test
func Init(t *testing.T) *TestTools {
	cfg, err := scw.LoadConfig()
	if err != nil {
		cfg = &scw.Config{}
	}
	httpClient, cleanupRecorder, err := getHTTPRecorder(t, UpdateCassettes)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to init http recorder: %w", err))
	}

	client, err := scw.NewClient(
		scw.WithHTTPClient(httpClient),
		scw.WithProfile(&cfg.Profile),
		scw.WithEnv(),
	)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to init scaleway client: %w", err))
	}
	tt := &TestTools{
		t:                t,
		HttpClient:       httpClient,
		scwClient:        client,
		cleanupFunctions: []cleanUpFunction{cleanupRecorder},
		cleanupRecorder:  cleanupRecorder,
	}
	return tt
}

// CleanUp cleans resources created with TestTools
func (tt *TestTools) CleanUp() {
	tt.t.Log("Cleaning up resources")
	errs := []error(nil)
	for _, cleanUpFunc := range tt.cleanupFunctions {
		err := cleanUpFunc()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		tt.t.Fatal(errors.Join(errs...))
	}
	err := tt.cleanupRecorder()
	if err != nil {
		tt.t.Fatal(err)
	}
}

func (tt *TestTools) genName(prefix string) string {
	return fmt.Sprintf("waypoint-%s-%s", prefix, strings.ToLower(tt.t.Name()))
}

// isHTTPCodeError returns true if err is an http error with code statusCode
func isHTTPCodeError(err error, statusCode int) bool {
	if err == nil {
		return false
	}

	responseError := &scw.ResponseError{}
	if errors.As(err, &responseError) && responseError.StatusCode == statusCode {
		return true
	}
	return false
}

// is404Error returns true if err is an HTTP 404 error
func is404Error(err error) bool {
	notFoundError := &scw.ResourceNotFoundError{}
	return isHTTPCodeError(err, http.StatusNotFound) || errors.As(err, &notFoundError)
}
