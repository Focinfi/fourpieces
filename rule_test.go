package fourpieces

import (
	"testing"
)

var inRangeTable = []struct {
	x, y    int
	inRange bool
}{
	{0, 0, true},
}

func TestInRange(t *testing.T) {
	for i, table := range inRangeTable {
		if inRange(table.x, table.y) != table.inRange {
			t.Errorf("test-%d failed", i)
		}
	}
}
