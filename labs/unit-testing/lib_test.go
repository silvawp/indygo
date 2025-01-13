package lib

import (
	"math/rand"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestInitialize(t *testing.T) {
	var lib MyLib

	lib.Initialize(StatusIdle)

	if lib.stateInfo != StatusIdle {
		t.Logf("Unexpected stateInfo: %d, expected: %d", lib.stateInfo, StatusIdle)
	}

	if len(lib.memory) != DefaultMemorySizeInBytes {
		t.Logf("Unexpected default mem size: %d, expected: %d", len(lib.memory), DefaultMemorySizeInBytes)
	}

	if lib.Status() != StatusIdle {
		t.Log("unexpected status")
	}
}

func TestMemory(t *testing.T) {
	// Test setup
	var lib MyLib
	setup(&lib)

	lib.memory[8191] = byte(255)

	mem := lib.Memory()

	if !slices.Equal(mem, lib.memory) {
		t.Log("memories aren't equal!")
		t.Fail()
	}

	i, ok := lib.FindAndUpdateFirst(255, 0xa5)
	if ok != nil {
		t.Log("value 255 should found at index 8191")
		t.Fail()
	}

	if i != 8191 {
		t.Logf("unexpected index: %d", i)
		t.Fail()
	}

	i, ok = lib.FindAndUpdateFirst(255, 0xa5)
	if ok == nil {
		t.Log("FindAndUpdateFirst should fail to find 255 and replace it to 0xa5")
		t.Fail()
	}

	if i >= 0 {
		t.Logf("unexpected index: %d", i)
		t.Fail()
	}

	t.Log("Expected error:", ok)
}

func TestLock(t *testing.T) {
	var lib MyLib
	setup(&lib)
	m := make([]byte, 256*1024*1024)
	fillWithRandomBytes(m)
	lib.SetMemory(m)

	var wg sync.WaitGroup
	var N = 5

	var f = func(id int) {
		defer wg.Done()
		x, err := lib.FindAndUpdateFirst(255, 0)
		if err == nil {
			t.Log("255 found at ", x)
			t.Fail()
		}
		t.Log("f: done: ", id)
	}

	var g = func(id int, s byte, r byte) {
		defer wg.Done()

		nReplaces := lib.ReplaceAll(s, r)

		if nReplaces == 0 {
			t.Logf("g[%d]: no '%d' found", id, s)
		} else {
			t.Logf("g[%v]: value '%d' replaced by '%d' %d times", id, s, r, nReplaces)
		}
	}

	wg.Add(2 * N)
	for i := range N {
		go f(100 + i)
		var search = "0123456789ABCDEF"[rand.Int()%16]
		var replace = "0123456789ABCDEF"[rand.Int()%16]
		go g(200+i, search, replace)
	}

	wg.Wait()
}

func BenchmarkLock(b *testing.B) {
	var lib MyLib
	setup(&lib)
	m := make([]byte, 256*1024*1024)
	fillWithRandomBytes(m)
	lib.SetMemory(m)

	b.ResetTimer()
	for nTest := 0; nTest < b.N; nTest++ {
		var wg sync.WaitGroup
		var N = 5
		var f = func(id int) {
			defer wg.Done()
			x, err := lib.FindAndUpdateFirst(255, 0)
			if err == nil {
				b.Log("255 found at ", x)
				b.Fail()
			}
			b.Log("f: done: ", id)
		}

		var g = func(id int, s byte, r byte) {
			defer wg.Done()

			nReplaces := lib.ReplaceAll(s, r)

			if nReplaces == 0 {
				b.Logf("g[%d]: no '%d' found", id, s)
			} else {
				b.Logf("g[%v]: value '%d' replaced by '%d' %d times", id, s, r, nReplaces)
			}
		}

		wg.Add(2 * N)
		for i := range N {
			go f(100 + i)
			var search = "0123456789ABCDEF"[rand.Int()%16]
			var replace = "0123456789ABCDEF"[rand.Int()%16]
			go g(200+i, search, replace)
		}

		wg.Wait()
	}

	b.StopTimer()
}

func FuzzToUpper(f *testing.F) {
	f.Add("hello")

	f.Fuzz(func(t *testing.T, s string) {
		out := strings.ToUpper(s)
		if out != strings.ToUpper(s) {
			t.Fail()
		}
	})
}

func setup(lib *MyLib) {
	lib.memory = make([]byte, 8192)
	fillWithRandomBytes(lib.memory)
}

func fillWithRandomBytes(m []byte) {
	for i := range len(m) {
		m[i] = byte("0123456789ABCDEF"[rand.Int()%16])
	}
}
