package main

import (
	"testing"
)

func Benchmark_listTaskWithRedis(b *testing.B) {
	s := NewService()
	for i := 0; i < b.N; i++ {
		s.GetTaskWithRedis(1)
	}
}

func Benchmark_listTask(b *testing.B) {
	s := NewService()
	for i := 0; i < b.N; i++ {
		s.GetTask(1)
	}
}
