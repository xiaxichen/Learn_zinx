package znet

import (
	"errors"
	"learn_zinx/zinx/logger"
	"learn_zinx/zinx/ziface"
	"sync"
)

/*
链接管理模块
*/
type ConnectionManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                  // 读写锁

}

func (cmr *ConnectionManager) Add(connection ziface.IConnection) {
	// 包含共享资源map 加写锁
	cmr.connLock.Lock()
	defer cmr.connLock.Unlock()

	// 将conn加入到ConnManager
	cmr.connections[connection.GetConnID()] = connection
	logger.Log.Debugf("connection add to connections ConnId=%d", connection.GetConnID())
}

func (cmr *ConnectionManager) Remove(connection ziface.IConnection) error {
	// 包含共享资源map 加写锁
	cmr.connLock.Lock()
	defer cmr.connLock.Unlock()
	_, ok := cmr.connections[connection.GetConnID()]
	if ok {
		// 将conn从ConnManager删除
		delete(cmr.connections, connection.GetConnID())
		logger.Log.Debugf("remove connection add to ConnManager successFuly:conn Id=%d ;connection num=%d", connection.GetConnID(), cmr.Len())
		return nil

	} else {
		logger.Log.Warnf("connection Id not in Connections ConnId=%d", connection.GetConnID())
		return errors.New("Connection undefined！")
	}
}

func (cmr *ConnectionManager) Get(connID uint32) (ziface.IConnection, error) {
	// 包含共享资源map 加读锁
	cmr.connLock.RLock()
	defer cmr.connLock.RUnlock()
	connection, ok := cmr.connections[connID]
	if ok {
		return connection, nil
	} else {
		return nil, errors.New("connection ID undefined")
	}
}

func (cmr *ConnectionManager) Len() int {
	return len(cmr.connections)
}

func (cmr *ConnectionManager) Clear() {
	// 包含共享资源map 加写锁
	cmr.connLock.Lock()
	defer cmr.connLock.Unlock()
	// 删除并停止Conn的工作
	for ConnId, connection := range cmr.connections {
		connection.Stop()
		logger.Log.Debugf("connection stop ConnID=%d", ConnId)
		delete(cmr.connections, ConnId)
		logger.Log.Debugf("connection remove ConnID=%d", ConnId)
	}
	logger.Log.Infof("Clear All connections success! conn num=%d", cmr.Len())
}

func MewConnManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    sync.RWMutex{},
	}
}
