package database

import (
	"sync"

	"github.com/fabjan/robotoscope/core"
)

type robotMap struct {
	lock sync.RWMutex
	data map[string]int
}

func NewRobotMap() robotMap {
	return robotMap{
		data: make(map[string]int),
	}
}

func (m *robotMap) Count(name string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	count := m.data[name]
	m.data[name] = count + 1
	return nil
}

func (m *robotMap) List() ([]core.RobotInfo, error) {
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
