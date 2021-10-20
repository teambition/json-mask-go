package jsonmask

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"compress/gzip"
	"compress/zlib"
)

// JSONMask
type JSONMask struct {
	queryKey string
	next     http.Handler
}

func NewJSONMask(next http.Handler, queryKey string) (http.Handler, error) {
	if next == nil {
		return nil, fmt.Errorf("http.Handler required")
	}
	if queryKey == "" {
		return nil, fmt.Errorf("queryKey required")
	}
	return &JSONMask{queryKey, next}, nil
}

func (h *JSONMask) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fields := req.URL.Query().Get(h.queryKey)
	if fields == "" {
		h.next.ServeHTTP(rw, req)
		return
	}

	sl, err := Compile(fields)
	if err != nil {
		h.next.ServeHTTP(rw, req)
		return
	}
	jrw := NewResponseWriter(rw, sl)
	defer jrw.Close()

	h.next.ServeHTTP(jrw, req)
}

type ResponseWriter struct {
	buf    bytes.Buffer
	rw     http.ResponseWriter
	sl     Selection
	status JSONMaskResponseWriterStatus
}

type JSONMaskResponseWriterStatus struct {
	Status  int
	Size    int
	Trimmed int
}

type ResponseWriterWithCloseNotify struct {
	*ResponseWriter
}

type JSONMaskResponseWriter interface {
	http.ResponseWriter
	Close() error
	Status() JSONMaskResponseWriterStatus
}

func NewResponseWriter(rw http.ResponseWriter, sl Selection) JSONMaskResponseWriter {
	jrw := &ResponseWriter{rw: rw, sl: sl}
	if _, ok := rw.(http.CloseNotifier); !ok {
		return jrw
	}
	return &ResponseWriterWithCloseNotify{jrw}
}

func (jrw *ResponseWriterWithCloseNotify) CloseNotify() <-chan bool {
	return jrw.rw.(http.CloseNotifier).CloseNotify()
}

func (jrw *ResponseWriter) Header() http.Header {
	return jrw.rw.Header()
}

func (jrw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := jrw.buf.Write(b)
	jrw.status.Size += size
	return size, err
}

func (jrw *ResponseWriter) WriteHeader(s int) {
	jrw.status.Status = s
}

func (jrw *ResponseWriter) Close() error {
	data := jrw.buf.Bytes()
	if jrw.status.Status == http.StatusOK {
		encoding := jrw.rw.Header().Get("Content-Encoding")
		if d, err := uncompress(encoding, data); err == nil {
			if d, err = jrw.sl.Mask(d); err == nil {
				if d, err = compress(encoding, d); err == nil {
					data = d
					jrw.rw.Header().Set("Content-Length", strconv.Itoa(len(d)))
					jrw.status.Trimmed = jrw.status.Size - len(d)
				}
			}
		}
	}
	jrw.rw.WriteHeader(jrw.status.Status)
	_, err := jrw.rw.Write(data)
	return err
}

func (jrw *ResponseWriter) Status() JSONMaskResponseWriterStatus {
	return jrw.status
}

func (jrw *ResponseWriter) Flush() {
	if f, ok := jrw.rw.(http.Flusher); ok {
		f.Flush()
	}
}

func (jrw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := jrw.rw.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("not a hijacker: %T", jrw.rw)
}

func compress(encoding string, data []byte) ([]byte, error) {
	var err error
	var wc io.WriteCloser
	var b bytes.Buffer

	switch encoding {
	case "gzip":
		wc = gzip.NewWriter(&b)
	case "deflate":
		wc = zlib.NewWriter(&b)
	}
	if wc == nil {
		return data, nil
	}
	if _, err = wc.Write(data); err != nil {
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func uncompress(encoding string, data []byte) ([]byte, error) {
	var err error
	var rc io.ReadCloser

	switch encoding {
	case "gzip":
		rc, err = gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	case "deflate":
		rc, err = zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	}

	if rc == nil {
		return data, nil
	}
	var b bytes.Buffer
	if _, err := io.Copy(&b, rc); err != nil {
		return nil, err
	}
	if err := rc.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
