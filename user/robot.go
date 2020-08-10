package user

import (
	"container/list"
	"sync"
)

//机器人容器，并发安全，功能较简单，复杂机器人策略可参考自己实现
type RobotQueue struct {
	robots *list.List
	sync.Mutex
}

func NewRobotQueue(robots []Account) *RobotQueue {
	pool := &RobotQueue{
		robots: list.New(),
	}
	for i := 0; i < len(robots); i++ {
		pool.Enqueue(&robots[i])
	}
	return pool
}

func (rc *RobotQueue) Dequeue() *Account {
	rc.Lock()
	defer rc.Unlock()
	elem := rc.robots.Front()
	if elem == nil {
		return nil
	}
	inst := elem.Value.(*Account)
	rc.robots.Remove(elem)
	return inst
}

func (rc *RobotQueue) Enqueue(rbs ...*Account) {
	rc.Lock()
	defer rc.Unlock()
	for _, rb := range rbs {
		rc.robots.PushBack(rb)
	}
}

func (rc *RobotQueue) Size() int {
	rc.Lock()
	defer rc.Unlock()
	return rc.robots.Len()
}
