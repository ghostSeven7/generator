package generator

import (
	"fmt"
	"math"
	"errors"
	"sync"
	"time"
)

var (
	//最大机器数量
	maxMachine uint8 = math.MaxUint8

	machineShift uint8 = 3
	stepShift    uint8 = 4
	timeShift          = stepShift + machineShift

	//每次增长步长
	step = 1

	//基础时间 毫秒 时间协调
	epoch int64 = 1559551969000
)

type node struct {
	sync.Mutex
	//机器号
	machine uint8
	//当前步数
	step int64

	time time.Time
	last int64

	machineShift,
	stepShift,
	timeShift uint8

	stepMask int64
}

/*
@param m uint8 机器号 支持256台
 */

func New(m uint8) (n *node, err error) {
	if m > maxMachine {
		return nil, errors.New(fmt.Sprintf("machine number must be between 0 and %d", math.MaxInt8))
	}
	n = new(node)
	n.machine = m
	n.machineShift = machineShift
	n.stepShift = stepShift
	n.timeShift = timeShift
	n.stepMask = -1 ^ (-1 << stepShift)
	var cutTime = time.Now()
	n.time = cutTime.Add(time.Unix(epoch/1000, (epoch%1000)*1000000).Sub(cutTime))
	return
}

func (n *node) generator() (id int64) {
	n.Lock()
	now := time.Since(n.time).Nanoseconds() / 1000000

	//同一毫秒
	if now == n.last {
		if n.step = (n.step + 1) & n.stepMask; n.step == 0 {
			//等待到下一毫秒
			for now <= n.last {
				now = time.Since(n.time).Nanoseconds() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.last = now

	id = (int64(n.machine) << n.machineShift) | (now << n.timeShift) | (n.step)
	n.Unlock()
	return
}
