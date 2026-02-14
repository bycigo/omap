//go:build goexperiment.jsonv2

package omap

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

// MarshalJSON handles JSON marshaling for the Map.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := jsontext.NewEncoder(&buf)
	if err := m.MarshalJSONTo(enc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalJSONTo encodes the Map into JSON using the provided encoder.
func (m *Map[K, V]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	for k, v := range m.All() {
		// write key
		if err := json.MarshalEncode(enc, k, json.StringifyNumbers(true)); err != nil {
			return err
		}

		// write value
		if err := json.MarshalEncode(enc, v); err != nil {
			return err
		}
	}

	if err := enc.WriteToken(jsontext.EndObject); err != nil {
		return err
	}

	return nil
}

// UnmarshalJSON handles JSON unmarshaling for the Map.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	dec := jsontext.NewDecoder(bytes.NewReader(data))
	return m.UnmarshalJSONFrom(dec)
}

// UnmarshalJSONFrom decodes JSON data into the Map using the provided decoder.
func (m *Map[K, V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if kind := dec.PeekKind(); kind != '{' {
		return fmt.Errorf("expected object")
	}

	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	for {
		kind := dec.PeekKind()
		if kind == '}' {
			if _, err := dec.ReadToken(); err != nil {
				return err
			}
			return nil
		}

		var k K
		if err := json.UnmarshalDecode(dec, &k); err != nil {
			return err
		}

		var v V
		if err := json.UnmarshalDecode(dec, &v); err != nil {
			return err
		}

		m.Set(k, v)
	}
}
