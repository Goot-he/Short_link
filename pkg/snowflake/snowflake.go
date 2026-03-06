package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	// 部分长度
	workerIDBits   = 7  // 机器ID位数
	sequenceIDBits = 3  // 时钟序列ID位数
	sequenceBits   = 12 // 序列号位数

	// 最大值
	maxWorkerID   = (1 << workerIDBits) - 1
	maxSequenceID = (1 << sequenceIDBits) - 1 // 最大的时钟序列ID（3位）
	maxSequence   = (1 << sequenceBits) - 1   // 最大的序列号（4095）

	// 位移量
	workerIDShift   = sequenceBits + sequenceIDBits
	sequenceIDShift = sequenceBits
	timestampShift  = workerIDBits + sequenceIDBits + sequenceBits
	epoch           = int64(1735660800000)
)

// 起始时间戳（纪元时间）

type SnowflakeIDGenerator struct {
	workerID      int64
	sequence      int64
	sequenceID    int64
	lastTimestamp int64
	mu            sync.Mutex
}

func NewSnowflakeIDGenerator(workerID int64) (*SnowflakeIDGenerator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("workerID out of range")
	}
	return &SnowflakeIDGenerator{
		workerID:      workerID,
		sequence:      0,
		sequenceID:    0, // 时钟序列ID初始值
		lastTimestamp: -1,
	}, nil
}

func (s *SnowflakeIDGenerator) currentTimestamp() int64 {
	return time.Now().UnixNano()/int64(time.Millisecond) - epoch
}

func (s *SnowflakeIDGenerator) waitForNextMillis(lastTimestamp int64) int64 {
	timestamp := s.currentTimestamp()
	for timestamp <= lastTimestamp {
		timestamp = s.currentTimestamp()
	}
	return timestamp
}

// 生成ID
func (s *SnowflakeIDGenerator) GenerateID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := s.currentTimestamp()

	// 如果当前时间戳与上次时间戳不同，重置时钟序列ID和序列号
	if timestamp != s.lastTimestamp {
		s.sequence = 0
		s.sequenceID = 0 // 时间戳不同，重置时钟序列ID
		s.lastTimestamp = timestamp
	} else if s.sequence < maxSequence {
		s.sequence++
	} else if s.sequenceID < maxSequenceID {
		// 如果序列号达到最大值，使用时钟序列ID增加，防止冲突
		s.sequenceID++
		s.sequence = 0
	} else {
		// 如果时钟序列ID也达到最大值，等待下一个毫秒
		timestamp = s.waitForNextMillis(s.lastTimestamp)
		s.lastTimestamp = timestamp
		s.sequence = 0
		s.sequenceID = 0
	}

	// 生成ID: 时间戳 + 工作机器ID + 时钟序列ID + 序列号
	id := (timestamp << timestampShift) |
		(s.workerID << workerIDShift) |
		(s.sequenceID << sequenceIDShift) |
		s.sequence

	return id
}
