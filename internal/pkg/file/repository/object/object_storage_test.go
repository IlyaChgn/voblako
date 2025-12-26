package repository

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewTestClient(rt http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: rt,
	}
}

func TestObjectStorage_UploadFile(t *testing.T) {
	client := NewTestClient(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		header := make(http.Header)

		// HEAD запрос для проверки bucket
		if req.Method == http.MethodHead {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     header,
			}, nil
		}

		// PUT запрос для загрузки файла
		if req.Method == http.MethodPut {
			header.Set("ETag", "\"d41d8cd98f00b204e9800998ecf8427e\"")
			header.Set("Content-Length", "0")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     header,
				Request:    req,
			}, nil
		}

		// GET запрос для получения региона bucket
		if req.Method == http.MethodGet && !strings.Contains(req.URL.Path, "test-key") {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">us-east-1</LocationConstraint>`)),
				Header:     header,
			}, nil
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     header,
		}, nil
	}))

	mc, err := minio.New("localhost:9000", &minio.Options{
		Creds:     nil,
		Secure:    false,
		Transport: client.Transport,
	})
	assert.NoError(t, err)

	storage := NewObjectStorage(mc, "test-bucket")

	tests := []struct {
		name        string
		key         string
		contentType string
		file        []byte
		size        int64
		wantErr     bool
	}{
		{
			name:        "OK",
			key:         "test-key",
			contentType: "text/plain",
			file:        []byte("test data"),
			size:        9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.UploadFile(context.Background(), tt.key, tt.contentType, tt.file, tt.size)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObjectStorage_GetFile(t *testing.T) {
	testData := "test data"

	client := NewTestClient(roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		header := make(http.Header)

		// HEAD запрос для проверки bucket
		if req.Method == http.MethodHead {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     header,
			}, nil
		}

		// GET запрос для получения файла
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "test-key") {
			header.Set("Content-Length", "9")
			header.Set("Content-Type", "text/plain")
			header.Set("Accept-Ranges", "bytes")
			header.Set("ETag", "\"d41d8cd98f00b204e9800998ecf8427e\"")
			header.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			return &http.Response{
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(strings.NewReader(testData)),
				Header:        header,
				ContentLength: 9,
				Request:       req,
			}, nil
		}

		// GET запрос для получения региона bucket
		if req.Method == http.MethodGet {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">us-east-1</LocationConstraint>`)),
				Header:     header,
			}, nil
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
			Header:     header,
		}, nil
	}))

	mc, err := minio.New("localhost:9000", &minio.Options{
		Creds:     nil,
		Secure:    false,
		Transport: client.Transport,
	})
	assert.NoError(t, err)

	storage := NewObjectStorage(mc, "test-bucket")

	tests := []struct {
		name    string
		key     string
		want    []byte
		wantErr bool
	}{
		{
			name: "OK",
			key:  "test-key",
			want: []byte("test data"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storage.GetFile(context.Background(), tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
