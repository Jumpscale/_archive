package rest

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRewriteURL(t *testing.T) {
	tt := []struct {
		src    string
		dest   string
		expect string
	}{
		{
			"http://localhost/store/ns/68b329da9893e34099c7d8ad5cb9c940",
			"https://stor.jumpscale.org/storx",
			"https://stor.jumpscale.org/storx/store/ns/68b329da9893e34099c7d8ad5cb9c940",
		},
		{
			"http://localhost/store/ns/68b329da9893e34099c7d8ad5cb9c940",
			"https://stor.jumpscale.org",
			"https://stor.jumpscale.org/store/ns/68b329da9893e34099c7d8ad5cb9c940",
		},
		{
			"http://localhost/path1/store/ns/68b329da9893e34099c7d8ad5cb9c940",
			"https://stor.jumpscale.org",
			"https://stor.jumpscale.org/store/ns/68b329da9893e34099c7d8ad5cb9c940",
		},
	}

	for _, test := range tt {
		dest, err := url.Parse(test.dest)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		src, err := url.Parse(test.src)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		newURL := rewriteURL(src, dest)
		assert.Equal(t, test.expect, newURL.String())
	}

}
