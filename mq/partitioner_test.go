package mq

import "testing"

const testKey = "abcdefghijklmn"

func BenchmarkUID2Index(b *testing.B) {
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		UID2Index(testKey)
	}
	b.StopTimer()
}

func BenchmarkUID2IndexSDBMHash(b *testing.B) {
	b.ResetTimer()

	for index := 0; index < b.N; index++ {
		UID2IndexSDBMHash(testKey)
	}
	b.StopTimer()
}

func BenchmarkUID2IndexCRC32(b *testing.B) {
	b.ResetTimer()

	for index := 0; index < b.N; index++ {
		UID2IndexCRC32(testKey)
	}
	b.StopTimer()
}
