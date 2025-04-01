package api

import (
	"reflect"
	"testing"
)

func TestIdentity4x4(t *testing.T) {
	m1 := Identity4x4()
	m2 := Matrix4x4{}
	m2[0][0] = 1
	m2[1][1] = 1
	m2[2][2] = 1
	m2[3][3] = 1

	if !reflect.DeepEqual(m1, m2) {
		t.Errorf("Identity4x4(): \n%#v\n%#v", m1, m2)
	}
}
