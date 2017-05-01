package drainchecker

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type eofDetectRoundTripper struct {
	http.RoundTripper
}

// RoundTripper returns a new http.RoundTripper which prints a message to
// stderr when the HTTP response body is not properly drained.
func RoundTripper(upstream http.RoundTripper) http.RoundTripper {
	return eofDetectRoundTripper{upstream}
}

type eofDetectReader struct {
	eofSeen bool
	rd      io.ReadCloser
}

func (rd *eofDetectReader) Read(p []byte) (n int, err error) {
	n, err = rd.rd.Read(p)
	if err == io.EOF {
		rd.eofSeen = true
	}
	return n, err
}

func (rd *eofDetectReader) Close() error {
	if !rd.eofSeen {
		buf, err := ioutil.ReadAll(rd)
		msg := fmt.Sprintf("body not drained, %d bytes not read", len(buf))
		if err != nil {
			msg += fmt.Sprintf(", error: %v", err)
		}

		if len(buf) > 0 {
			if len(buf) > 20 {
				buf = append(buf[:20], []byte("...")...)
			}
			msg += fmt.Sprintf(", body: %q", buf)
		}

		fmt.Fprintln(os.Stderr, msg)
	}
	return rd.rd.Close()
}

func (tr eofDetectRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	res, err = tr.RoundTripper.RoundTrip(req)
	res.Body = &eofDetectReader{rd: res.Body}
	return res, err
}
