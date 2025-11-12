package library

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

// ConsistentHash 一致性哈希接口
// T 是节点类型，需要能够转换为字符串用于哈希计算
type ConsistentHash[T any] interface {
	// Add 添加节点到哈希环
	Add(node T)
	// Remove 从哈希环移除节点
	Remove(node T)
	// Get 根据 key 获取对应的节点
	Get(key any) (T, bool)
	// GetNodes 获取所有节点
	GetNodes() []T
	// GetNodeCount 获取节点数量
	GetNodeCount() int
}

// HashFunc 哈希函数类型
type HashFunc func(data []byte) uint32

// Stringer 将节点转换为字符串的接口
type Stringer interface {
	String() string
}

// consistentHashImpl 一致性哈希实现
type consistentHashImpl[T any] struct {
	mu           sync.RWMutex
	hashFunc     HashFunc
	replicas     int            // 虚拟节点数量
	ring         map[uint32]T   // 哈希环：hash值 -> 节点
	sortedHashes []uint32       // 排序后的哈希值，用于快速查找
	nodeMap      map[string]T   // 节点名称 -> 节点，用于去重
	toString     func(T) string // 节点转字符串函数
}

// NewConsistentHash 创建一致性哈希实例
// replicas: 每个节点的虚拟节点数量，建议值 150-200
// hashFunc: 哈希函数，如果为 nil 则使用默认的 CRC32
// toString: 将节点转换为字符串的函数，如果为 nil 且节点实现了 Stringer 接口则使用 String() 方法
func NewConsistentHash[T any](replicas int, hashFunc HashFunc, toString func(T) string) ConsistentHash[T] {
	if replicas <= 0 {
		replicas = 150 // 默认虚拟节点数
	}

	if hashFunc == nil {
		hashFunc = func(data []byte) uint32 {
			return crc32.ChecksumIEEE(data)
		}
	}

	impl := &consistentHashImpl[T]{
		hashFunc:     hashFunc,
		replicas:     replicas,
		ring:         make(map[uint32]T),
		sortedHashes: make([]uint32, 0),
		nodeMap:      make(map[string]T),
		toString:     toString,
	}

	// 如果没有提供 toString 函数，尝试使用 Stringer 接口
	if toString == nil {
		impl.toString = func(node T) string {
			if stringer, ok := any(node).(Stringer); ok {
				return stringer.String()
			}
			return fmt.Sprintf("%v", node)
		}
	}

	return impl
}

// Add 添加节点到哈希环
func (c *consistentHashImpl[T]) Add(node T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	nodeStr := c.toString(node)

	// 如果节点已存在，先移除
	if _, exists := c.nodeMap[nodeStr]; exists {
		c.removeNode(nodeStr)
	}

	// 添加节点到节点映射
	c.nodeMap[nodeStr] = node

	// 为节点创建虚拟节点
	for i := 0; i < c.replicas; i++ {
		hash := c.hashFunc([]byte(fmt.Sprintf("%s:%d", nodeStr, i)))
		c.ring[hash] = node
		c.sortedHashes = append(c.sortedHashes, hash)
	}

	// 重新排序
	sort.Slice(c.sortedHashes, func(i, j int) bool {
		return c.sortedHashes[i] < c.sortedHashes[j]
	})
}

// Remove 从哈希环移除节点
func (c *consistentHashImpl[T]) Remove(node T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	nodeStr := c.toString(node)
	c.removeNode(nodeStr)
}

// removeNode 内部方法：移除节点
func (c *consistentHashImpl[T]) removeNode(nodeStr string) {
	if _, exists := c.nodeMap[nodeStr]; !exists {
		return
	}

	// 删除所有虚拟节点
	newHashes := make([]uint32, 0, len(c.sortedHashes))
	for i := 0; i < c.replicas; i++ {
		hash := c.hashFunc([]byte(fmt.Sprintf("%s:%d", nodeStr, i)))
		delete(c.ring, hash)
	}

	// 重建排序数组
	for _, hash := range c.sortedHashes {
		if _, exists := c.ring[hash]; exists {
			newHashes = append(newHashes, hash)
		}
	}

	c.sortedHashes = newHashes
	delete(c.nodeMap, nodeStr)
}

// Get 根据 key 获取对应的节点
func (c *consistentHashImpl[T]) Get(key any) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.ring) == 0 {
		var zero T
		return zero, false
	}

	hash := c.hashFunc([]byte(fmt.Sprintf("%v", key)))

	// 使用二分查找找到第一个大于等于 hash 的节点
	idx := sort.Search(len(c.sortedHashes), func(i int) bool {
		return c.sortedHashes[i] >= hash
	})

	// 如果没找到（说明 hash 大于所有节点），则使用第一个节点（环形结构）
	if idx == len(c.sortedHashes) {
		idx = 0
	}

	node := c.ring[c.sortedHashes[idx]]
	return node, true
}

// GetNodes 获取所有节点（去重后）
func (c *consistentHashImpl[T]) GetNodes() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	nodes := make([]T, 0, len(c.nodeMap))
	for _, node := range c.nodeMap {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetNodeCount 获取节点数量
func (c *consistentHashImpl[T]) GetNodeCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.nodeMap)
}

// MD5HashFunc 使用 MD5 的哈希函数（用于更好的分布）
func MD5HashFunc(data []byte) uint32 {
	hash := md5.Sum(data)
	return uint32(hash[0])<<24 | uint32(hash[1])<<16 | uint32(hash[2])<<8 | uint32(hash[3])
}

// CRC32HashFunc 使用 CRC32 的哈希函数（默认，性能更好）
func CRC32HashFunc(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
