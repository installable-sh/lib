package fetch

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/installable-sh/lib/certs"
	"github.com/installable-sh/lib/log"
)

// Script represents a fetched shell script.
type Script struct {
	Content string
	Name    string
}

// Options configures how a script is fetched.
type Options struct {
	URL     string
	SendEnv bool
	NoCache bool
}

// NewClient creates an HTTP client with system and embedded CA certificates.
// Debug output is controlled by the logger's debug level.
func NewClient(logger log.DebugLogger) (*retryablehttp.Client, error) {
	certPool, err := certs.CertPool(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificates: %w", err)
	}

	logger.Debugf("Creating HTTP client with TLS config")
	client := retryablehttp.NewClient()
	client.RetryMax = 0 // Unlimited retries
	client.Logger = nil // Silence debug logging
	client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}

	return client, nil
}

// Fetch retrieves a script from a URL.
// Debug output is controlled by the logger's debug level.
func Fetch(ctx context.Context, client *retryablehttp.Client, opts Options, logger log.DebugLogger) (Script, error) {
	logger.Debugf("Creating request for %s", opts.URL)
	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", opts.URL, nil)
	if err != nil {
		return Script{}, err
	}

	userAgent := "run/1.0 (installable)"
	if ua := os.Getenv("USER_AGENT"); ua != "" {
		userAgent = ua
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/plain, text/x-shellscript, application/x-sh, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	logger.Debugf("Request headers: User-Agent=%s", userAgent)

	if opts.NoCache {
		req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		req.Header.Set("Pragma", "no-cache")
		logger.Debugf("Cache disabled via headers")
	}

	if opts.SendEnv {
		envCount := 0
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 && isValidHeaderName(parts[0]) {
				req.Header.Set("X-Env-"+parts[0], parts[1])
				envCount++
			}
		}
		logger.Debugf("Sending %d environment variables as X-Env-* headers", envCount)
	}

	logger.Debugf("Executing HTTP GET request")
	resp, err := client.Do(req)
	if err != nil {
		return Script{}, logger.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	logger.Debugf("Response: HTTP %d, Content-Type=%s, Content-Encoding=%s",
		resp.StatusCode,
		resp.Header.Get("Content-Type"),
		resp.Header.Get("Content-Encoding"))

	if resp.StatusCode != http.StatusOK {
		return Script{}, logger.Errorf("HTTP %d from %s", resp.StatusCode, opts.URL)
	}

	name := scriptName(resp, opts.URL)
	content, err := readBody(resp)
	if err != nil {
		return Script{}, err
	}

	logger.Debugf("Received script: name=%s, size=%d bytes", name, len(content))

	return Script{Content: content, Name: name}, nil
}

func isValidHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for _, c := range name {
		// HTTP header names must be tokens (RFC 7230)
		// Allow: A-Z a-z 0-9 ! # $ % & ' * + - . ^ _ ` | ~
		if c <= ' ' || c >= 127 || strings.ContainsRune("\"(),/:;<=>?@[\\]{}", c) {
			return false
		}
	}
	return true
}

func scriptName(resp *http.Response, url string) string {
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		_, params, err := mime.ParseMediaType(cd)
		if err == nil && params["filename"] != "" {
			return params["filename"]
		}
	}

	name := path.Base(url)
	if name == "" || name == "/" || name == "." {
		return "script.sh"
	}
	return name
}

func readBody(resp *http.Response) (string, error) {
	var reader io.Reader = resp.Body

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("gzip error: %w", err)
		}
		defer func() { _ = gzReader.Close() }()
		reader = gzReader
	case "deflate":
		reader = flate.NewReader(resp.Body)
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
