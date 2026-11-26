// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncore

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

// These cover the length checks backported from v2.9.2. Each input panicked
// with a slice bounds or index error on v2.3.1.

func TestMalformedLengthsDoNotPanic(t *testing.T) {
	t.Run("string length that overflows int32", func(t *testing.T) {
		// A string element whose length is MaxInt32: adding the 4 length bytes wraps negative.
		elem := append([]byte{byte(TypeString), 'k', 0x00}, 0xff, 0xff, 0xff, 0x7f, 'x')
		if _, _, ok := ReadElement(elem); ok {
			t.Fatal("ReadElement accepted an element whose length overflows")
		}

		doc := binary.LittleEndian.AppendUint32(nil, uint32(4+len(elem)+1))
		doc = append(append(doc, elem...), 0x00)
		if err := Document(doc).Validate(); err == nil {
			t.Fatal("Validate accepted a document whose element length overflows")
		}
	})

	t.Run("document shorter than the minimum", func(t *testing.T) {
		// Length 0, with bytes still following: d[length-1] indexed before the slice.
		doc := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		if err := Document(doc).Validate(); !errors.Is(err, ErrInvalidLength) {
			t.Fatalf("Document.Validate() = %v, want ErrInvalidLength", err)
		}
		if err := Array(doc).Validate(); !errors.Is(err, ErrInvalidLength) {
			t.Fatalf("Array.Validate() = %v, want ErrInvalidLength", err)
		}
	})

	t.Run("reader with a length below the minimum", func(t *testing.T) {
		// Length 3: the 4 length bytes don't fit in a 3-byte buffer.
		if _, err := NewDocumentFromReader(bytes.NewReader([]byte{0x03, 0x00, 0x00, 0x00})); !errors.Is(err, ErrInvalidLength) {
			t.Fatalf("NewDocumentFromReader() error = %v, want ErrInvalidLength", err)
		}
	})
}
