//go:build !goexperiment.jsonv2

package omap

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

// MarshalJSON handles JSON marshaling for the Map.
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('{')

	first := true
	for k, v := range m.All() {
		// marshal key
		key, err := json.Marshal(map[K]uint8{k: 0})
		if err != nil {
			return nil, errors.New("unsupported key type")
		}
		// extract the actual key from `{"key":0}`
		key = bytes.TrimPrefix(key, []byte{'{'})
		key = bytes.TrimSuffix(key, []byte{':', '0', '}'})

		// marshal value
		value, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		if first {
			first = false
		} else {
			buf.WriteByte(',')
		}
		buf.Write(key)
		buf.WriteByte(':')
		buf.Write(value)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON handles JSON unmarshaling for the Map.
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	if !bytes.HasPrefix(data, []byte{'{'}) {
		return errors.New("expected JSON object")
	}

	dec := json.NewDecoder(bytes.NewReader(data))

	// skip '{'
	if _, err := dec.Token(); err != nil {
		return err
	}

	for dec.More() {
		// unmarshal key
		kt, err := dec.Token()
		if err != nil {
			return err
		}
		var key K
		kv := reflect.ValueOf(&key).Elem()
		switch kv.Kind() {
		case reflect.String:
			kv.SetString(kt.(string))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(kt.(string), 10, 64)
			if err != nil {
				return errors.New("invalid integer key: " + kt.(string))
			}
			kv.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u, err := strconv.ParseUint(kt.(string), 10, 64)
			if err != nil {
				return errors.New("invalid unsigned integer key: " + kt.(string))
			}
			kv.SetUint(u)
		default:
			return errors.New("unsupported key type: " + kv.Type().String())
		}

		// unmarshal value
		var value V
		if err = dec.Decode(&value); err != nil {
			return err
		}

		m.Set(key, value)
	}

	return nil
}
