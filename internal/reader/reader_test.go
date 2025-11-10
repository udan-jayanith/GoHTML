package Reader_test

import (
	"io"
	"strings"
	"testing"

	Reader "github.com/udan-jayanith/GoHTML/internal/reader"
)

//Test Loop.
//io.reader returns a error early.
//byteReads exceed.
//Loop break early.
//Exceed maximum empty reads

func Test_Iter_Loop(t *testing.T) {
	input := "Hello world"
	rd := Reader.NewReader(strings.NewReader(input))
	iter := Reader.NewIter(rd)

	accumulatedStr := ""
	for byt := range iter.Loop() {
		accumulatedStr += string(byt)
	}

	if iter.Err != io.EOF {
		t.Fatal("Expected error io.EOF but got", iter.Err.Error())
	} else if accumulatedStr != input {
		t.Fatal("Expected", input, "but got", accumulatedStr)
	}
}

func Test_Iter_ByteReadsExceed(t *testing.T) {
	rd := Reader.NewReaderConfigured(strings.NewReader("Hello world"), 1, Reader.DefaultMaxConsecutiveEmptyReads, len("hello"))
	iter := rd.Iter()

	accumulatedStr := ""
	for byt := range iter.Loop() {
		accumulatedStr += string(byt)
	}

	if accumulatedStr != "Hello" {
		t.Fatal("Expected 'Hello' but got", accumulatedStr)
	} else if iter.Err != Reader.ReachedMaxBytesRead {
		t.Fatal("Expected", Reader.ReachedMaxBytesRead, "but got", iter.Err.Error())
	}
}

func Test_Iter_Loop_BreakEarly(t *testing.T) {
	rd := Reader.NewReader(strings.NewReader("Hello world"))
	iter := rd.Iter()

	accumulatedStr := ""
	for byt := range iter.Loop() {
		if accumulatedStr == "Hello" {
			break
		}
		accumulatedStr += string(byt)
	}

	if iter.Err != nil {
		t.Fatal("Expected no error but got", iter.Err.Error())
	} else if accumulatedStr != "Hello" {
		t.Fatal("Expected 'Hello' but got", accumulatedStr)
	}
}

type emptyReads struct {
}

func (rc *emptyReads) Read(p []byte) (int, error) {
	return 0, nil
}

func Test_Iter_ExceedMinimumEmptyReads(t *testing.T) {
	rd := Reader.NewReader(&emptyReads{})
	iter := rd.Iter()
	for byt := range iter.Loop() {
		t.Fatal("Byte read", byt, "but expected not byte reads")
	}

	if iter.Err != Reader.ReachedMaxConsecutiveEmptyReads {
		t.Fatal("Expected error", Reader.ReachedMaxConsecutiveEmptyReads, "but got", iter.Err)
	}
}
