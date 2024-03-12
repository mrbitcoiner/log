package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestCreate(t *testing.T){
	w := bytes.NewBuffer([]byte{})

	_, err := NewLog(WithWriter(w), WithLevel(LOGERR))

	if err != nil { t.Fatal(err) }
}

func TestLogLevel(t *testing.T) {
	w := bytes.NewBuffer([]byte{})
	log, err := NewLog(WithWriter(w), WithLevel(LOGERR))
	if err != nil { t.Fatal(err) }

	res := log.toLog(LOGWARN)
	res1 := log.toLog(LOGERR)
	res2 := log.toLog(LOGFATAL)
	res3 := log.toLog(LOGTRACE + 1)

	if res || !res1 || !res2 || res3 { 
		t.Fatal("unexpected result:", res, res1, res2, res3) 
	}
}

func TestErrOnInvalidLogLevel(t *testing.T) {
	_, err := NewLog(WithLevel(7))

	if !errors.As(err, interface{}(&ErrInvalidLogLevel)) {
		t.Fatal("expected invalid log level")
	}

}

func TestLogWritesToBuffer(t *testing.T) {
	w := bytes.NewBuffer([]byte{})
	log, err := NewLog(WithWriter(w), WithLevel(LOGERR))
	if err != nil { t.Fatal(err) }

	str := "testing 12344321"
	log.Err(str)

	if !strings.Contains(string(w.Bytes()), str) {
		t.Fatal("unexpected buffer:", string(w.Bytes()))
	}
}

func TestLogFatalPanics(t *testing.T){
	w := bytes.NewBuffer([]byte{})
	log, err := NewLog(WithWriter(w), WithLevel(LOGERR))
	if err != nil { t.Fatal(err) }
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("expecting panic")
		} 
		perror, ok := err.(error)
		if !ok {
			t.Fatal("unexpected panic:", err)
		}
		if !errors.Is(ErrLogFatal, perror) {
			t.Fatal("unexpected error type:", err)
		}
	}()

	// expected panic
	log.Fatal("testing log fatal")
}

func TestLogFatalFPanics(t *testing.T){
	w := bytes.NewBuffer([]byte{})
	log, err := NewLog(WithWriter(w), WithLevel(LOGERR))
	if err != nil { t.Fatal(err) }
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("expecting panic")
		} 
		perror, ok := err.(error)
		if !ok {
			t.Fatal("unexpected panic:", err)
		}
		if !errors.Is(ErrLogFatal, perror) {
			t.Fatal("unexpected error type:", err)
		}
	}()

	// expected panic
	log.Fatalf("testing log fatal")
}

