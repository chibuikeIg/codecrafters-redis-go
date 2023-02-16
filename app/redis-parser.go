package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Type represents a Value type
type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

type Value struct {
	typ   Type
	bytes []byte
	array []Value
}

// String converts Value to a string.
//
// If Value cannot be converted, an empty string is returned.
func (v Value) String() string {
	if v.typ == BulkString || v.typ == SimpleString {
		return string(v.bytes)
	}

	return ""
}

// Array converts Value to an array.
//
// If Value cannot be converted, an empty array is returned.
func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}

	return []Value{}
}

func readBytesUntilCRLF(bytestream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}

	for {
		b, err := bytestream.ReadBytes('\n')

		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)

		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}

	return readBytes[:len(readBytes)-2], nil
}

func decodeArray(bytestream *bufio.Reader) (Value, error) {

	readBytesForCount, err := readBytesUntilCRLF(bytestream)

	if err != nil {
		return Value{}, fmt.Errorf("Failed to read bulk string length %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))

	if err != nil {
		return Value{}, fmt.Errorf("Failed to parse bulk string length %s", err)
	}

	array := []Value{}

	for i := 1; i <= count; i++ {
		value, err := DecodeRESP(bytestream)

		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		typ:   Array,
		array: array,
	}, nil
}

func DecodeRESP(bytestream *bufio.Reader) (Value, error) {

	dataTypeByte, err := bytestream.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return decodeSimpleString(bytestream)
	case "$":
		return decodeBulkString(bytestream)
	case "*":
		return decodeArray(bytestream)

	}

	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

func decodeSimpleString(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readBytesUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

func decodeBulkString(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readBytesUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}

	return Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}
