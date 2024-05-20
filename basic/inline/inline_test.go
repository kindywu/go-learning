package inline

import "testing"

func BenchmarkNoInline(b *testing.B) {
	x, y := 1, 2
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = maxNoinline(x, y)
	}
}

func BenchmarkInline(b *testing.B) {
	x, y := 1, 2
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = maxInline(x, y)
	}
}
