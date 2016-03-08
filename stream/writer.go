package stream

import (
	"encoding/json"
	"io"
)

// newline separator
var newLine = []byte{'\n'}

// Writer implements a stream writer, that serializes
// build streams in a structured format
type Writer struct {
	JSONWriter
}

// NewWriter returns a new stream writer used for writing
// build stream packets to writer w.
func NewWriter(w io.Writer) *Writer {
	return NewWriterJSON(&jsonWriter{w, json.NewEncoder(w)})
}

// NewWriterJSON returns a new stream writer used for
// writing build stream packets to JSONWriter w. This is
// used primarily for writing to Websockets.
func NewWriterJSON(w JSONWriter) *Writer {
	return &Writer{w}
}

// Write writes the bytes to a packet which is in turn
// written to the writer in JSON format.
func (w *Writer) Write(d []byte) (int, error) {
	p := &Packet{
		Type: PacketLine,
		Data: string(d),
	}
	if err := w.WriteJSON(p); err != nil {
		return 0, err
	}
	return len(d), nil
}

// JSONWriter defines a Writer that is capable of writing
// JSON encoded objects.
type JSONWriter interface {
	WriteJSON(v interface{}) error
}

type jsonWriter struct {
	writer  io.Writer
	encoder *json.Encoder
}

func (w *jsonWriter) WriteJSON(v interface{}) error {
	if err := w.encoder.Encode(v); err != nil {
		return err
	}
	return nil
	//_, err := w.writer.Write(newLine)
	// return err
}
