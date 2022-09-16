package stopwatch

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var outputs = make(map[string]*os.File)
var measurements = make(map[string]chan *measurement)
var DefaultOutput *os.File

func SetOutput(label string, f *os.File) {
	outputs[label] = f
}

func Measure(label string, f func()) {
	MeasureWithComment(label, "", f)
}
func MeasureWithComment(label string, comment string, f func()) {
	newMeasurement := prepareMeasurement(label)
	newMeasurement.comment = comment

	defer newMeasurement.Stop()
	newMeasurement.Start()
	f()
}

var now *measurement
var nowLock sync.Mutex

func Now(label string) {
	nowLock.Lock()
	defer nowLock.Unlock()
	if now == nil {
		now = prepareMeasurement(label)
		now.Start()
	}
	now.Stop()
	now = prepareMeasurement(label)
	now.Start()
}

func prepareMeasurement(label string) *measurement {
	m := &measurement{}

	_, ok := measurements[label]
	if !ok {
		measurements[label] = make(chan *measurement, 100000000)
	}
	measurements[label] <- m
	return m
}
func FlushSingle(label string, series chan *measurement) {
	f, ok := outputs[label]
	if !ok {
		f = DefaultOutput
	}
	fmt.Println("Flushing measurements to output", label)
	for m := range series {
		if f != nil && !m.end.IsZero() {
			_, _ = fmt.Fprintln(f, m.Duration().Nanoseconds(), "\t", m.comment)
		}
	}
}

func Flush() {
	fmt.Println("Start flushing measurements")
	wg := &sync.WaitGroup{}
	for label, data := range measurements {
		wg.Add(1)
		go func(l string, s chan *measurement) {
			defer wg.Done()
			FlushSingle(l, s)
		}(label, data)
		measurements[label] = nil
		close(data)
	}
	wg.Wait()
}

type measurement struct {
	start, end time.Time
	comment    string
}

func (m *measurement) Duration() time.Duration {
	return m.end.Sub(m.start)
}

func (m *measurement) Start() {
	if m.start.IsZero() {
		m.start = time.Now()
	}
}
func (m *measurement) Stop() {
	if !m.start.IsZero() && m.end.IsZero() {
		m.end = time.Now()
	}
}
