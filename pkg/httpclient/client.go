package httpclient

import (
	"errors"
	"github.com/valyala/fasthttp"
	"math/rand"
	"time"
)

var ErrNilResponse = errors.New("unexpected nil response")
var ErrNon200Response = errors.New("API responded with non-200 status code")

type Header struct {
	Key   string
	Value string
}

func MakeRequest(c *fasthttp.Client, url string, maxRetries int, headers ...Header) ([]byte, error) {
	var (
		req  *fasthttp.Request
		resp *fasthttp.Response
	)
	retries := maxRetries
	for i := retries; i >= 0; i-- {
		req = fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.Header.SetMethod(fasthttp.MethodGet)
		for _, header := range headers {
			req.Header.Set(header.Key, header.Value)
		}
		req.Header.Set(fasthttp.HeaderUserAgent, getUserAgent())
		req.SetRequestURI(url)

		resp = fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		if err := c.DoTimeout(req, resp, time.Second*45); err != nil {
			fasthttp.ReleaseRequest(req)
			if retries == 0 {
				return nil, err
			}
		}

		if resp.Body() == nil {
			if retries == 0 {
				return nil, ErrNilResponse
			}
		}
		// url responded with 503, so try again
		if resp.StatusCode() == 503 {
			continue
		}
	}

	if resp.StatusCode() != 200 {
		return nil, ErrNon200Response
	}
	return resp.Body(), nil
}

func getUserAgent() string {
	payload := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:66.0) Gecko/20100101 Firefox/66.0",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:67.0) Gecko/20100101 Firefox/67.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12H321 Safari/600.1.4",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (iPad; CPU OS 7_1_2 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D257 Safari/9537.53",
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)",
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(payload))

	pick := payload[randomIndex]

	return pick
}
