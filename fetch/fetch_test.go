package fetch

import (
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		opts        Options
		wantContent string
		wantName    string
		wantErr     bool
	}{
		{
			name: "basic fetch",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("echo hello"))
			},
			opts:        Options{URL: ""},
			wantContent: "echo hello",
			wantName:    "script.sh",
		},
		{
			name: "with content-disposition",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Disposition", `attachment; filename="test.sh"`)
				_, _ = w.Write([]byte("echo test"))
			},
			opts:        Options{URL: ""},
			wantContent: "echo test",
			wantName:    "test.sh",
		},
		{
			name: "gzip encoded",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Encoding", "gzip")
				gz := gzip.NewWriter(w)
				_, _ = gz.Write([]byte("echo gzipped"))
				_ = gz.Close()
			},
			opts:        Options{URL: ""},
			wantContent: "echo gzipped",
			wantName:    "script.sh",
		},
		{
			name: "404 error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			opts:    Options{URL: ""},
			wantErr: true,
		},
		{
			name: "no-cache headers",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Cache-Control") != "no-cache, no-store, must-revalidate" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.Header.Get("Pragma") != "no-cache" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				_, _ = w.Write([]byte("echo nocache"))
			},
			opts:        Options{URL: "", NoCache: true},
			wantContent: "echo nocache",
			wantName:    "script.sh",
		},
		{
			name: "script name from URL path",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("echo path"))
			},
			opts:        Options{URL: ""}, // Will be set to server URL + /myscript.sh
			wantContent: "echo path",
			wantName:    "myscript.sh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			// Set URL to server URL with a script path
			if tt.name == "script name from URL path" {
				tt.opts.URL = server.URL + "/myscript.sh"
			} else {
				tt.opts.URL = server.URL + "/script.sh"
			}

			client, err := NewClient()
			if err != nil {
				t.Fatalf("NewClient() error: %v", err)
			}

			script, err := Fetch(context.Background(), client, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if script.Content != tt.wantContent {
				t.Errorf("Fetch() content = %q, want %q", script.Content, tt.wantContent)
			}

			if script.Name != tt.wantName {
				t.Errorf("Fetch() name = %q, want %q", script.Name, tt.wantName)
			}
		})
	}
}

func TestIsValidHeaderName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid simple", "FOO", true},
		{"valid with underscore", "FOO_BAR", true},
		{"valid with numbers", "FOO123", true},
		{"valid lowercase", "foo", true},
		{"empty", "", false},
		{"with space", "FOO BAR", false},
		{"with colon", "FOO:BAR", false},
		{"with slash", "FOO/BAR", false},
		{"with quotes", "FOO\"BAR", false},
		{"with brackets", "FOO[BAR]", false},
		{"non-ascii", "FOO\x80BAR", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidHeaderName(tt.input); got != tt.want {
				t.Errorf("isValidHeaderName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error: %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}
