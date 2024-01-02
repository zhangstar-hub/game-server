package note

import "testing"

// func TestForTest(t *testing.T) {
// 	ForTest(3)
// 	t.Log("TestForTest")
// }

func BenchmarkForTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ForTest(3)
	}
}
