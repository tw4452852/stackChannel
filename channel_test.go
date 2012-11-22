package stackChannel

import (
	"bytes"
	"testing"
)

func TestChannel(t *testing.T) {
	stuffs := [][]byte{
		[]byte("hello"),
		[]byte("tw"),
		[]byte("\n"),
	}
	c := newStackChannel()
	//write to channel
	var (
		length   int
		expected []byte
	)
	for _, stuff := range stuffs {
		if _, err := c.Write(stuff); err != nil {
			t.Fatalf("write stuff(%q) to channel failed: %s\n", stuff, err)
		}
		length += len(stuff)
		expected = append(expected, stuff...)
	}
	//read from channel
	result := make([]byte, length)
	if n, err := c.Read(result); err != nil || n != length {
		t.Fatalf("read %d from channel failed: %s\n", n, err)
	}
	//verify the result
	if bytes.Compare(result, expected) != 0 {
		t.Fatalf("result(%v) != expected(%v)\n", result, expected)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("close channel failed: %s\n", err)
	}
}

func BenchmarkChannel(b *testing.B) {
	stuffs := [][]byte{
		[]byte("hello"),
		[]byte("tw"),
		[]byte("\n"),
	}
	c := newStackChannel()
	var (
		length   int
		expected []byte
	)
	//write to channel
	for i := 0; i < b.N; i++ {
		b.Log(i)
		length = 0
		expected = expected[:0]
		for _, stuff := range stuffs {
			if _, err := c.Write(stuff); err != nil {
				b.Fatalf("write stuff(%q) to channel failed: %s\n", stuff, err)
			}
			length += len(stuff)
			expected = append(expected, stuff...)
		}
		//read from channel
		result := make([]byte, length)
		if n, err := c.Read(result); err != nil || n != length {
			b.Fatalf("read %d from channel failed: %s\n", n, err)
		}
		//verify the result
		if bytes.Compare(result, expected) != 0 {
			b.Fatalf("result(%v) != expected(%v)\n", result, expected)
		}
	}
	if err := c.Close(); err != nil {
		b.Fatalf("close channel failed: %s\n", err)
	}
}
