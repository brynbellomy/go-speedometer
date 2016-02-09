package speedometer

import (
	"sync"
	"time"
)

type Speedometer struct {
	mutex *sync.RWMutex

	startTime     time.Time
	lapTime       time.Time
	count         uint64
	countSinceLap uint64
}

func New() *Speedometer {
	return &Speedometer{
		mutex: &sync.RWMutex{},
	}
}

func (r *Speedometer) Start() {
	r.mutex.Lock()

	r.startTime = time.Now()
	r.lapTime = time.Now()
	r.count = 0
	r.countSinceLap = 0

	r.mutex.Unlock()
}

func (r *Speedometer) Incr(delta uint64) {
	r.mutex.Lock()
	r.count += delta
	r.countSinceLap += delta
	r.mutex.Unlock()
}

func (r *Speedometer) Speed() Speed {
	r.mutex.RLock()
	count := r.countSinceLap
	duration := time.Now().Sub(r.lapTime)
	r.mutex.RUnlock()
	return Speed{count, duration}
}

func (r *Speedometer) Lap() Speed {
	speed := r.Speed()

	r.mutex.Lock()

	r.lapTime = time.Now()
	r.countSinceLap = 0

	r.mutex.Unlock()
	return speed
}

func (r *Speedometer) GetCount() uint64 {
	var count uint64
	r.mutex.RLock()
	count = r.count
	r.mutex.RUnlock()
	return count
}

func (r *Speedometer) GetCountSinceLap() uint64 {
	var count uint64
	r.mutex.RLock()
	count = r.countSinceLap
	r.mutex.RUnlock()
	return count
}

type Speed struct {
	count    uint64
	duration time.Duration
}

func (s Speed) PerSecond() float64 {
	c := float64(s.count)
	d := s.duration.Seconds()
	return c / d
}

func (s Speed) PerNanosecond() float64 {
	c := float64(s.count)
	d := float64(int64(s.duration))
	return c / d
}
