package database

import (
	"sync"

	"github.com/fabjan/robotoscope/core"
)

// RobotMap is an in-memory RobotStore.
type RobotMap struct {
	lock sync.RWMutex
	data map[string]int
}

// NewRobotMap creates an fresh RobotMap.
func NewRobotMap() RobotMap {
	return RobotMap{
		data: make(map[string]int),
	}
}

// Count increases the seen count for the given bot.
func (m *RobotMap) Count(name string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	count := m.data[name]
	m.data[name] = count + 1
	return nil
}

// List returns a list showing how many times each robot has been seen.
func (m *RobotMap) List() ([]core.RobotInfo, error) {
	info := []core.RobotInfo{}

	m.lock.RLock()
	defer m.lock.RUnlock()

	for robot, count := range m.data {
		i := core.RobotInfo{
			Seen:      count,
			UserAgent: robot,
		}
		info = append(info, i)
	}

	return info, nil
}
