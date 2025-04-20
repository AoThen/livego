package pool

import (
	"sync"
)

type Pool struct {
	mu  sync.Mutex // 添加互斥锁保证并发安全
	pos int
	buf []byte
}

const maxpoolsize = 1000 * 1024 // 1MB 池大小

func (p *Pool) Get(size int) []byte {
	if size <= 0 {
		return nil // 处理无效大小请求
	}

	// 处理超过池大小的请求
	if size > maxpoolsize {
		return make([]byte, size)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// 检查剩余空间是否足够
	if maxpoolsize-p.pos < size {
		p.pos = 0 // 重置位置复用现有缓冲区
	}

	// 分配字节切片并更新位置
	result := p.buf[p.pos : p.pos+size]
	p.pos += size

	return result
}

func NewPool() *Pool {
	return &Pool{
		buf: make([]byte, maxpoolsize), // 预分配内存
	}
}

// package pool

// type Pool struct {
// 	pos int
// 	buf []byte
// }

// const maxpoolsize = 1000 * 1024

// func (pool *Pool) Get(size int) []byte {
// 	if maxpoolsize-pool.pos < size {
// 		pool.pos = 0
// 		pool.buf = make([]byte, maxpoolsize)
// 	}
// 	b := pool.buf[pool.pos : pool.pos+size]
// 	pool.pos += size
// 	return b
// }

// func NewPool() *Pool {
// 	return &Pool{
// 		buf: make([]byte, maxpoolsize),
// 	}
// }
