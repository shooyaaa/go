package library

import (
	"time"
)

type Op uint

const Op_Smaller Op = 1
const Op_Bigger Op = 2
const Op_Range Op = 3

type OpValue interface {
}

type Limitation struct {
	op    Op
	value OpValue
	unit  TimeUnit
}

type TimeUnit uint

const HourUnit TimeUnit = 1
const DayUnit TimeUnit = 2
const MonthUnit TimeUnit = 3

func timeUnitValue(unit TimeUnit) int {
	now := time.Now()
	switch unit {
	case HourUnit:
		return now.Hour()
	case DayUnit:
		return int(now.Weekday())
	case MonthUnit:
		return int(now.Month())
	}
	return -1
}

func (l *Limitation) satisfied() bool {
	target := timeUnitValue(l.unit)
	switch l.op {
	case Op_Smaller:
		return l.value.(int) > target
	case Op_Bigger:
		return l.value.(int) < target
	case Op_Range:
		for _, item := range l.value.([]int) {
			if item == target {
				return true
			}
		}
		return false
	}
	return false
}

type CronJob struct {
	interval    int
	count       int
	times       int
	callback    func()
	limitations []*Limitation
}

func (c *CronJob) AddLimitation(unit TimeUnit, op Op, value interface{}) {
	limitation := &Limitation{
		op:    op,
		value: value,
		unit:  unit,
	}
	c.limitations = append(c.limitations, limitation)
}

func (c *CronJob) valid() bool {
	for _, value := range c.limitations {
		if !value.satisfied() {
			return false
		}
	}
	return true
}

type Ticker struct {
	closeChan chan int
	jobs      []*CronJob
}

func (t *Ticker) Cron(interval int, callback func()) *CronJob {
	return t.CronTimes(-1, interval, callback)
}

func (t *Ticker) Once(interval int, callback func()) *CronJob {
	return t.CronTimes(1, interval, callback)
}

func (t *Ticker) CronTimes(times int, interval int, callback func()) *CronJob {
	cronJob := CronJob{
		interval: interval,
		count:    0,
		callback: callback,
		times:    times,
	}
	t.jobs = append(t.jobs, &cronJob)
	return &cronJob
}

func (t *Ticker) Start() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for _, job := range t.jobs {
				job.count++
				if job.count >= job.interval && (job.times > 0 || job.times < 0) && job.valid() {
					go job.callback()
					job.count = 0
					if job.times > 0 {
						job.times--
					}
				}
			}
		case <-t.closeChan:
			return
		default:
			time.Sleep(time.Millisecond * 40)
		}
	}
}

func (t *Ticker) Stop() {
	t.closeChan <- 1
}

func NewTicker() *Ticker {
	return &Ticker{
		closeChan: make(chan int),
	}
}
