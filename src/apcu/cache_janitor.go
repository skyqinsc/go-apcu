package apcu

import (
	"container/heap"
	"sync"
	"time"
)

const (
	MaxWatcherSize int64 = 1000
)

type janitorNode struct {
	key         string
	expiredTime int64
}

type priorityQueue []*janitorNode

type janitorQueue struct {
	priorityQueue priorityQueue
	mu            sync.RWMutex
	interval      time.Duration
	watcher       chan *janitorNode
	stop          chan bool
}

//Len ...
func (pq priorityQueue) Len() int { return len(pq) }

//Less ...
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].expiredTime != 0 || pq[i].expiredTime < pq[j].expiredTime
}

//Swap ...
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq priorityQueue) Push(node interface{}) {
	pq = append(pq, node.(*janitorNode))
}

func (pq priorityQueue) Pop() interface{} {
	n := len(pq)
	node := pq[n-1]
	pq = pq[0 : n-1]

	return node
}

func (pq priorityQueue) Top() interface{} {
	n := len(pq)
	node := pq[n-1]

	return node
}

func (pq priorityQueue) Empty() bool {
	if len(pq) == 0 {
		return true
	}

	return false
}

func newJanitor(interval time.Duration, maxSize int64) *janitorQueue {
	queue := &janitorQueue{
		interval:      interval,
		priorityQueue: make([]*janitorNode, 0),
		watcher:       make(chan *janitorNode, maxSize),
		stop:          make(chan bool),
	}
	heap.Init(&queue.priorityQueue)
	return queue
}

func (j *janitorQueue) notify(k string, expiredTime int64) {
	node := &janitorNode{
		key:         k,
		expiredTime: expiredTime,
	}
	j.watcher <- node
}

func (j *janitorQueue) insert(node *janitorNode) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.priorityQueue.Push(node)
}

func (j *janitorQueue) flush(c *cache) {
	j.mu.Lock()
	defer j.mu.Unlock()
	for !j.priorityQueue.Empty() {
		node := j.priorityQueue.Top().(*janitorNode)
		if node.expiredTime > time.Now().UnixNano() {
			return
		}
		delete(c.items, node.key)
		heap.Pop(j.priorityQueue)
	}
}

func (j *janitorQueue) run(c *cache) {
	go func() {
		ticker := time.NewTicker(j.interval)

		for {
			select {
			case <-ticker.C:
				j.flush(c)
			case node := <-j.watcher:
				j.insert(node)
			}
		}
	}()
}
