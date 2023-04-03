package server

import (
	"github.com/panjf2000/gnet"
	"math/rand"
	"sync"
)

// 跳表

// SkipNode 链表节点
type SkipNode struct {
	id   uint64
	ip   uint32
	link gnet.Conn
	next []*SkipNode //各层的后向节点指针数组，数组的长度为层高，level
}

// SkipList 跳跃表
type SkipList struct {
	SkipNode              //跳表的header
	mutex    sync.RWMutex //锁
	update   []*SkipNode  //查询过程中的链式变量
	maxL     int          //最大层数，32
	skip     int          //层之间的比例，skip=4，1/4节点出现再上层
	level    int          //跳表当前层数
	length   int32        //调表的节点数
}

// NewSkipList 初始化跳跃表
func NewSkipList() *SkipList {
	l := &SkipList{}
	l.maxL = 32
	l.skip = 4
	l.SkipNode.next = make([]*SkipNode, l.maxL)
	l.update = make([]*SkipNode, l.maxL)
	return l
}

// Get 节点查找
func (l *SkipList) Get(id uint64, ip uint32) gnet.Conn {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	var prev = &l.SkipNode
	var next *SkipNode
	//1.从底层链表开始查询
	for i := l.level - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.id < id {
			prev = next
			next = prev.next[i]
		}
	}
	if next != nil && next.id == id && next.ip == ip {
		return next.link
	} else {
		return nil
	}
}

/*
*
1. 查找待插入的位置，需要获取每层的前驱节点
2.构造新节点，通过概率函数计算节点的层数level
3.讲新节点插入到第0层到第level-1层的链表中
*/
func (l *SkipList) Set(id uint64, ip uint32, link gnet.Conn) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	var prev = &l.SkipNode
	var next *SkipNode
	for i := l.level - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.id < id {
			prev = next
			next = prev.next[i]
		}
		l.update[i] = prev
	}

	//如果key已经存在
	if next != nil && next.id == id && next.ip == ip {
		next.link = link
		return
	}

	//随机生成新结点的层数
	level := l.randomLevel()
	if level > l.level {
		level = l.level + 1
		l.level = level
		l.update[l.level-1] = &l.SkipNode
	}

	//申请新的结点
	node := &SkipNode{}
	node.id = id
	node.ip = ip
	node.link = link
	node.next = make([]*SkipNode, level)

	//调整next指向
	for i := 0; i < level; i++ {
		node.next[i] = l.update[i].next[i]
		l.update[i].next[i] = node
	}
	l.length++
}

func (l *SkipList) Remove(id uint64, ip uint32) gnet.Conn {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	//获取每层的前驱节点=>list.update
	var prev = &l.SkipNode
	var next *SkipNode
	for i := l.level - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.id < id {
			prev = next
			next = prev.next[i]
		}
		l.update[i] = prev
	}

	//结点不存在
	node := next
	if next == nil || next.id != id || next.ip != ip {
		return nil
	}

	//调整next指向
	for i, v := range node.next {
		if l.update[i].next[i] == node {
			l.update[i].next[i] = v
			if l.SkipNode.next[i] == nil {
				l.level -= 1
			}
		}
		l.update[i] = nil
	}

	l.length--
	return node.link
}

func (l *SkipList) randomLevel() int {
	i := 1
	for ; i < l.maxL; i++ {
		if rand.Int()%l.skip != 0 {
			break
		}
	}
	return i
}
