package util

import (
	"reflect"
	"testing"
)

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	if !reflect.DeepEqual(logger, GetLogger()) {
		t.Fatal("Logger is not singleton")
	}
}
