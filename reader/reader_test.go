package Reader_test

import (
	"io"
	"strings"
	"testing"

	Reader "github.com/udan-jayanith/GoHTML/reader"
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
func Test_Iter_NoneEOF_error_early(t *testing.T) {}

func Test_Iter_ByteReadsExceed(t *testing.T) {}

func Test_Iter_Loop_BreakEarly(t *testing.T) {}

func Test_Iter_ExceedMinimumEmptyReads(t *testing.T) {}
