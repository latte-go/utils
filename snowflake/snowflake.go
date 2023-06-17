package snowflake

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	epoch int64 = 1669852800 * 1000 // 开始时间截 (2022-12-01) 单位：毫秒

	stepBits uint8 = 12 // 以随机数开始的序列号占用位数
	nodeBits uint8 = 10 // 节点ID占用位数

	stepMask int64 = -1 ^ (-1 << stepBits) // 序列号最大值, 即每个时间单位(毫秒)可生成的最大数量
	nodeMask int64 = -1 ^ (-1 << nodeBits) // 节点ID最大值

	timeShift = stepBits + nodeBits // 时间戳左移位数
	stepShift = nodeBits            // 序列号左移位数
)

type generateType uint8

const (
	generateLong generateType = iota
	generateShort
)

// ID 雪花算法分布式唯一ID生成器
// 世界上没有两片完全相同的树叶(雪花) ——（德国）莱布尼茨
// 使用分布式节点时，最高支持1023个节点，若初始化时不设置节点或设置为0，将在生成ID时使用随机数代替
// 最高支持每毫秒4095个ID, 当序列号不足时等待至下一毫秒
//
//	+----------------------------------------------------------------------------+
//
// ｜01111111｜11111111｜11111111｜11111111｜11111111｜11111111｜11111111｜11111111｜
// ｜ 1bit代表正数 + 41bit毫秒级时间戳 + 12bit以随机数开始的自增序列号 + 10bit节点位   ｜
//
//	+----------------------------------------------------------------------------+
type ID int64

type Short int64

type Node struct {
	mu    sync.Mutex
	time  int64 // 时间戳毫秒计数
	node  int64 // 节点
	step  int64 // 序列号，毫秒内递增
	first int64 // 毫秒内序列号的开始,随机数
}

var (
	snowNode *Node
	nodeOnce sync.Once
)

// NewNode 返回 Node Struct 用于生成唯一 ID
// node 支持最大值1023，大于1023时node会与1023取模
// node 为0时，会以随机数代替
func NewNode(nodeId int64) *Node {
	nodeOnce.Do(func() {
		snowNode = &Node{node: nodeId & nodeMask}
	})
	return snowNode
}

// Long 生成一个64bit长度的、唯一的 ID
func (n *Node) Long() ID {
	return ID(n.generate(generateLong))
}

// Short 生成一个53bit长度的、唯一的 ID
func (n *Node) Short() Short {
	return Short(n.generate(generateShort))
}

// generate 生成唯一ID
func (n *Node) generate(t generateType) int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	var (
		node = n.node
		now  = time.Now().UnixMilli()
	)

	if node <= 0 && t == generateLong {
		node = rand.New(rand.NewSource(now)).Int63n(nodeMask)
	}

	if now <= n.time {
		// 时间单位内的序列号自增，并在序列号到达最大数时延迟到下一个时间单位
		n.step = (n.step + 1) & stepMask
		if n.step == n.first {
			for now <= n.time {
				now = time.Now().UnixMilli()
			}
			n.first = 0
		}
	} else {
		n.first = 0
	}

	if n.first == 0 {
		n.first = rand.New(rand.NewSource(now)).Int63n(stepMask-1) + 1
		n.step = n.first
	}

	n.time = now

	if t == generateShort {
		return (now-epoch)<<stepBits | n.step
	} else {
		return (now-epoch)<<timeShift |
			(n.step << stepShift) |
			(node)
	}
}

// ShortID 返回精度53bit、去除掉节点信息的唯一ID，可用于JS的交互
func (i ID) ShortID() Short {
	return Short(i >> nodeBits)
}

// Time 返回 ID 中存储的时间戳
func (i ID) Time() int64 {
	return (i.Int64() >> timeShift) + epoch
}

// Step 返回 ID 中存储的序列号
func (i ID) Step() int64 {
	return i.Int64() >> stepShift & stepMask
}

// Node 返回 ID 中存储的节点ID
func (i ID) Node() int64 {
	return i.Int64() & nodeMask
}

// Base2 返回 ID 的二进制字符串
func (i ID) Base2() string {
	return strconv.FormatInt(i.Int64(), 2)
}

func (i ID) Int64() int64 {
	return int64(i)
}

// String 返回 ID 的十进制字符串
func (i ID) String() string {
	return strconv.FormatInt(i.Int64(), 10)
}

func (i ID) Hex() string {
	return strconv.FormatInt(i.Int64(), 16)
}

func (i ID) Base64() string {
	return base64.StdEncoding.EncodeToString(i.Bytes())
}

func (i ID) Bytes() []byte {
	return []byte(i.String())
}

func (i ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(i))
	return b
}

func ParseInt64(id int64) ID {
	return ID(id)
}

// ParseBase2 转换二进制字符串到 ID
func ParseBase2(id string) (ID, error) {
	return convert(id, 2)
}

// ParseString 转换十进制字符串到 ID
func ParseString(id string) (ID, error) {
	return convert(id, 10)
}

func ParseHex(id string) (ID, error) {
	return convert(id, 16)
}

func ParseBase64(id string) (ID, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}
	return ParseBytes(b)
}

func ParseBytes(id []byte) (ID, error) {
	return convert(string(id), 10)
}

func ParseIntBytes(id [8]byte) ID {
	return ID(int64(binary.BigEndian.Uint64(id[:])))
}

func convert(id string, base int) (ID, error) {
	i, err := strconv.ParseInt(id, base, 64)
	return ID(i), err
}

func (s Short) ID() ID {
	return ID(s << nodeBits)
}

// Time 返回 Short 中存储的时间戳
func (s Short) Time() int64 {
	return (s.Int64() >> stepBits) + epoch
}

// Step 返回 Short 中存储的序列号
func (s Short) Step() int64 {
	return s.Int64() & stepMask
}

func (s Short) Int64() int64 {
	return int64(s)
}

func (s Short) String() string {
	return strconv.FormatInt(s.Int64(), 10)
}
