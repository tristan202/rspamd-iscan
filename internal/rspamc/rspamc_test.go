package rspamc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fho/rspamd-iscan/internal/log"
	"github.com/fho/rspamd-iscan/internal/testutils/assert"
)

func TestLearnSucceedsOnJSONResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer srv.Close()

	c := New(log.SlogTestLogger(t), srv.URL, "")

	assert.NoError(t, c.Ham(context.Background(), strings.NewReader("mail")))
	assert.NoError(t, c.Spam(context.Background(), strings.NewReader("mail")))
}

func TestCheckReturnsErrorWhenRspamdIsUnreachable(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()

	c := New(log.SlogTestLogger(t), url, "")

	res, err := c.Check(context.Background(), strings.NewReader("mail"))
	assert.Error(t, err)
	assert.Equal(t, nil, res)
}
