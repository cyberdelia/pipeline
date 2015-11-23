// Pipeline allows you to compose read and write pipeline from
// io.Reader and io.Writer.
package pipeline

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"testing"
)

func encodePipeline(w io.WriteCloser) (io.WriteCloser, error) {
	return base64.NewEncoder(base64.StdEncoding, w), nil
}

func decodePipeline(r io.ReadCloser) (io.ReadCloser, error) {
	return ioutil.NopCloser(base64.NewDecoder(base64.StdEncoding, r)), nil
}

type nopCloser struct {
	io.Writer
}

func (nc nopCloser) Close() error { return nil }

func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}

func TestPipeWriter(t *testing.T) {
	var output bytes.Buffer
	input := bytes.NewBufferString("aloha")
	w, err := PipeWrite(NopCloser(&output), encodePipeline)
	if err != nil {
		t.Error(err)
	}
	io.Copy(w, input)
	w.Close()
	if output.String() != "YWxvaGE=" {
		t.Errorf("unexpected output, wants %s, got %s", "YWxvaGE=", output.String())
	}
}

func TestPipeRead(t *testing.T) {
	var output bytes.Buffer
	var input bytes.Buffer
	input.WriteString("YWxvaGE=")
	r, err := PipeRead(ioutil.NopCloser(&input), decodePipeline)
	if err != nil {
		t.Error(err)
	}
	io.Copy(&output, r)
	r.Close()
	if output.String() != "aloha" {
		t.Errorf("unexpected output, wants %s, got %s", "aloha", output.String())
	}
}
