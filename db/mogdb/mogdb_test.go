package mogdb

import "testing"

func TestMogdbName(t *testing.T) {
	d := MogDBOpen("abc")
	if d.Name() != "mogdb" {
		t.Fatalf("mogdb Dialector name is not 'mogdb' but '%v'", d.Name())
	}
	t.Fatalf("mogdb Dialector name is '%v'", d.Name())
}
