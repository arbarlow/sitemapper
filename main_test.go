package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrape(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("./test")))
	defer ts.Close()

	parse(ts.URL)

	assert.Equal(t, len(results), 5)
}
