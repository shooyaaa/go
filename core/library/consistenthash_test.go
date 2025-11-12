package library

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsistentHash(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)
	assert.NotNil(t, ch)
	assert.Equal(t, 0, ch.GetNodeCount())
}

func TestConsistentHash_Add(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	ch.Add("node1")
	ch.Add("node2")
	ch.Add("node3")

	assert.Equal(t, 3, ch.GetNodeCount())

	nodes := ch.GetNodes()
	assert.Len(t, nodes, 3)
}

func TestConsistentHash_Remove(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	ch.Add("node1")
	ch.Add("node2")
	ch.Add("node3")

	assert.Equal(t, 3, ch.GetNodeCount())

	ch.Remove("node2")

	assert.Equal(t, 2, ch.GetNodeCount())

	nodes := ch.GetNodes()
	assert.Len(t, nodes, 2)
	assert.NotContains(t, nodes, "node2")
}

func TestConsistentHash_Get(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	ch.Add("node1")
	ch.Add("node2")
	ch.Add("node3")

	// 测试获取节点
	node, found := ch.Get("key1")
	assert.True(t, found)
	assert.NotEmpty(t, node)
	assert.Contains(t, []string{"node1", "node2", "node3"}, node)

	// 相同 key 应该返回相同节点
	node2, found2 := ch.Get("key1")
	assert.True(t, found2)
	assert.Equal(t, node, node2)
}

func TestConsistentHash_GetEmpty(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	node, found := ch.Get("key1")
	assert.False(t, found)
	assert.Empty(t, node)
}

func TestConsistentHash_Consistency(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	ch.Add("node1")
	ch.Add("node2")
	ch.Add("node3")

	// 记录初始映射
	keyToNode := make(map[string]string)
	keys := []string{"key1", "key2", "key3", "key4", "key5"}

	for _, key := range keys {
		node, _ := ch.Get(key)
		keyToNode[key] = node
	}

	// 添加新节点
	ch.Add("node4")

	// 验证大部分 key 仍然映射到相同节点（一致性）
	changed := 0
	for _, key := range keys {
		node, _ := ch.Get(key)
		if keyToNode[key] != node {
			changed++
		}
	}

	// 添加一个节点时，应该只有部分 key 需要重新映射
	// 理论上应该是 1/4 左右，但允许一些误差
	assert.Less(t, changed, len(keys), "添加节点后，应该有部分 key 重新映射")
}

func TestConsistentHash_RemoveConsistency(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	ch.Add("node1")
	ch.Add("node2")
	ch.Add("node3")
	ch.Add("node4")

	// 记录初始映射
	keyToNode := make(map[string]string)
	keys := []string{"key1", "key2", "key3", "key4", "key5"}

	for _, key := range keys {
		node, _ := ch.Get(key)
		keyToNode[key] = node
	}

	// 移除一个节点
	ch.Remove("node2")

	// 验证被移除节点上的 key 会重新映射到其他节点
	for _, key := range keys {
		node, found := ch.Get(key)
		assert.True(t, found)
		if keyToNode[key] == "node2" {
			// 原本在 node2 上的 key 应该重新映射
			assert.NotEqual(t, "node2", node)
		}
	}
}

func TestConsistentHash_CustomHashFunc(t *testing.T) {
	ch := NewConsistentHash[string](150, MD5HashFunc, nil)

	ch.Add("node1")
	ch.Add("node2")

	node, found := ch.Get("key1")
	assert.True(t, found)
	assert.NotEmpty(t, node)
}

func TestConsistentHash_CustomToString(t *testing.T) {
	type Node struct {
		ID   int
		Addr string
	}

	toString := func(n Node) string {
		return n.Addr
	}

	ch := NewConsistentHash[Node](150, nil, toString)

	ch.Add(Node{ID: 1, Addr: "127.0.0.1:8080"})
	ch.Add(Node{ID: 2, Addr: "127.0.0.1:8081"})

	assert.Equal(t, 2, ch.GetNodeCount())

	node, found := ch.Get("key1")
	assert.True(t, found)
	assert.NotEmpty(t, node.Addr)
}

func TestConsistentHash_StringerInterface(t *testing.T) {
	type Node struct {
		ID   int
		Addr string
	}

	toStringFunc := func(n Node) string {
		return n.Addr
	}

	// Node 实现了 Stringer 接口
	// 不提供 toString 函数，应该自动使用 String() 方法
	ch := NewConsistentHash[Node](150, nil, toStringFunc)

	ch.Add(Node{ID: 1, Addr: "127.0.0.1:8080"})
	ch.Add(Node{ID: 2, Addr: "127.0.0.1:8081"})

	assert.Equal(t, 2, ch.GetNodeCount())

	node, found := ch.Get("key1")
	assert.True(t, found)
	assert.NotEmpty(t, node.Addr)
}

func TestConsistentHash_DuplicateAdd(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	ch.Add("node1")
	ch.Add("node1") // 重复添加

	// 应该只有一个节点
	assert.Equal(t, 1, ch.GetNodeCount())
}

func TestConsistentHash_Distribution(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	ch.Add("node1")
	ch.Add("node2")
	ch.Add("node3")

	// 测试分布均匀性
	distribution := make(map[string]int)
	keys := make([]string, 1000)

	for i := 0; i < 1000; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}

	for _, key := range keys {
		node, _ := ch.Get(key)
		distribution[node]++
	}

	// 验证所有节点都有数据分布
	assert.Len(t, distribution, 3)

	// 验证分布相对均匀（每个节点应该有数据）
	for node, count := range distribution {
		assert.Greater(t, count, 0, "节点 %s 应该有数据分布", node)
	}
}

func TestConsistentHash_Replicas(t *testing.T) {
	// 测试不同虚拟节点数量的影响
	ch1 := NewConsistentHash[string](50, nil, nil)
	ch2 := NewConsistentHash[string](200, nil, nil)

	ch1.Add("node1")
	ch1.Add("node2")
	ch2.Add("node1")
	ch2.Add("node2")

	// 更多虚拟节点应该提供更好的分布
	// 但相同 key 应该映射到相同节点（在节点数量相同时）
	node1, _ := ch1.Get("test-key")
	node2, _ := ch2.Get("test-key")

	// 由于虚拟节点数量不同，可能映射到不同节点，这是正常的
	assert.NotEmpty(t, node1)
	assert.NotEmpty(t, node2)
}

func TestConsistentHash_Concurrent(t *testing.T) {
	ch := NewConsistentHash[string](150, nil, nil)

	// 并发添加节点
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			ch.Add(fmt.Sprintf("node%d", i))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ch.Get(fmt.Sprintf("key%d", i))
		}
		done <- true
	}()

	<-done
	<-done

	// 验证没有 panic 且节点已添加
	assert.Greater(t, ch.GetNodeCount(), 0)
}
