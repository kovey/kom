package meta

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	a := map[string]int{}
	vType := reflect.TypeOf(a)
	t.Log(vType.Name(), vType.Elem().Name(), vType.Key().Name())
}
