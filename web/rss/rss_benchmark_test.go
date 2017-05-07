package rss

import "testing"

func BenchmarkUpdate(b *testing.B) {
	Update()
}