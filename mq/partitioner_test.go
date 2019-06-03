package mq

import "testing"

func BenchmarkUID2Index(b *testing.B) {
	str := "abcdefghijklmn"

	b.ResetTimer()

	for index := 0; index < b.N; index++ {
		UID2Index(str)
	}
}

func BenchmarkUID2IndexSDBMHash(b *testing.B) {
	str := "abcdefghijklmn"

	b.ResetTimer()

	for index := 0; index < b.N; index++ {
		UID2IndexSDBMHash(str)
	}
}
