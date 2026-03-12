package object

import (
	"marabu/internal/crypto"
	"marabu/internal/messages"

	"encoding/json"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type ObjectManager struct {
	db           *leveldb.DB
	pendingFinds map[string][]chan interface{}
	mutex        sync.Mutex
}

func NewObjectManager(path string) (*ObjectManager, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &ObjectManager{
		db:           db,
		pendingFinds: make(map[string][]chan interface{}),
	}, nil
}

func (om *ObjectManager) Exists(id string) (bool, error) {
	return om.db.Has([]byte(id), nil)
}

func (om *ObjectManager) Get(id string) (interface{}, error) {
	data, err := om.db.Get([]byte(id), nil)
	if err != nil {
		return nil, err
	}

	var obj messages.ObjectSchema
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (om *ObjectManager) Put(object interface{}) (string, error) {
	canon, err := messages.Canonicalize(object)
	if err != nil {
		return "", err
	}
	id, err := crypto.HashString(canon)
	if err != nil {
		return "", err
	}

	// Marshal and store
	data, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	if err := om.db.Put([]byte(id), data, nil); err != nil {
		return "", err
	}
	return id, nil
}

// Implement FindObject with channels for pending requests
func (om *ObjectManager) FindObject(id string) (interface{}, error) {
	// First, try to get the object immediately
	obj, err := om.Get(id)
	if err == nil {
		return obj, nil
	}

	// If not found, set up a pending channel
	om.mutex.Lock()
	ch := make(chan interface{}, 1)
	om.pendingFinds[id] = append(om.pendingFinds[id], ch)
	om.mutex.Unlock()

	// Wait for the object to be provided by someone else (e.g., after a network fetch)
	result, ok := <-ch
	if !ok {
		return nil, fmt.Errorf("find for object %s was cancelled", id)
	}
	return result, nil
}

// When you later receive the object (e.g., after a network fetch and Put):
func (om *ObjectManager) notifyWaiters(id string, obj interface{}) {
	om.mutex.Lock()
	defer om.mutex.Unlock()
	for _, ch := range om.pendingFinds[id] {
		ch <- obj
		close(ch)
	}
	delete(om.pendingFinds, id)
}
