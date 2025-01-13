// Lib a package for a library (sharable code)
package lib

import (
	"fmt"
	"sync"
)

const (
	DefaultMemorySizeInBytes = 128
)

const (
	StatusIdle = iota
	StatusReceiving
	StatusProcessing
	StatusTransmitting
)

type MyLib struct {
	stateInfo int
	memory    []byte
	lock      sync.Mutex
}

func (l *MyLib) Initialize(s int) {
	l.memory = make([]byte, DefaultMemorySizeInBytes)
	l.stateInfo = s
}

func (l *MyLib) Status() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.stateInfo
}

func (l *MyLib) Memory() []byte {
	l.lock.Lock()
	mem := make([]byte, len(l.memory))
	copy(mem, l.memory)
	l.lock.Unlock()
	return mem
}

func (l *MyLib) SetMemory(m []byte) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.memory = make([]byte, 0)
	l.memory = m[:]

	// for _, v := range m {
	// 	l.memory = append(l.memory, v)
	// }
}

func (l *MyLib) FindAndUpdateFirst(b byte, r byte) (index int, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	for i, v := range l.memory {
		if v == b {
			index = i
			err = nil
			l.memory[i] = r
			return
		}
	}

	index = -1
	err = fmt.Errorf("value %v not found", b)
	return
}

func (l *MyLib) ReplaceAll(s byte, r byte) int {
	l.lock.Lock()
	defer l.lock.Unlock()

	m := make([]byte, len(l.memory))

	copy(m, l.memory)

	sum := 0
	for i, v := range m {
		if s == v {
			sum++
			m[i] = r
		}
	}

	l.memory = m

	return sum
}
