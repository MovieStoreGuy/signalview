package client_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/MovieStoreGuy/signalview/pkg/client"
	"github.com/stretchr/testify/require"
)

const (
	testToken = "butts"
)

func TestCachedRequests(t *testing.T) {
	t.Parallel()
	generator := client.NewCachedRequest(testToken)
	req, err := generator(context.TODO(), http.MethodGet, "http://localhost", nil)
	require.NoError(t, err)
	require.Equal(t, []string{testToken}, req.Header["X-Sf-Token"])
	require.Equal(t, []string{"application/json"}, req.Header["Content-Type"])
	require.Equal(t, []string{"signalview"}, req.Header["User-Agent"])

	_, err = generator(nil, "", "", nil)
	require.Error(t, err)

	req, err = generator(context.TODO(), http.MethodGet, "http://localhost", nil)
	require.NoError(t, err)
	require.Equal(t, []string{testToken}, req.Header["X-Sf-Token"])
	require.Equal(t, []string{"application/json"}, req.Header["Content-Type"])
	require.Equal(t, []string{"signalview"}, req.Header["User-Agent"])
}

func TestConfiguredClient(t *testing.T) {
	t.Parallel()
	c := client.NewConfiguredClient()
	require.Equal(t, c.Timeout, 16*time.Second)

	c = client.NewConfiguredClient(func(cl *http.Client) {
		cl.Timeout = 0
	})
	require.Equal(t, c.Timeout, 0)
}
