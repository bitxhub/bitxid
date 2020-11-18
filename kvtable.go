package bitxid

import (
	"fmt"

	"github.com/meshplus/bitxhub-kit/storage"
)

// KVTable .
type KVTable struct {
	store storage.Storage
}

// var tablelogger = log.NewWithModule("registry.Table")

var _ RegistryTable = (*KVTable)(nil)

// NewKVTable .
func NewKVTable(S storage.Storage) (*KVTable, error) {
	return &KVTable{
		store: S,
	}, nil
}

// HasItem whether table has the item(by key)
func (r *KVTable) HasItem(did DID) bool {
	exists := r.store.Has([]byte(did))
	return exists
}

// SetItem sets without any checks
func (r *KVTable) setItem(did DID, item TableItem) error {
	bitem, err := item.Marshal()
	if err != nil {
		return fmt.Errorf("kvtable marshal: %w", err)
	}
	r.store.Put([]byte(did), bitem)
	return nil
}

// CreateItem checks and sets
func (r *KVTable) CreateItem(item TableItem) error {
	did := item.GetID()
	if did == DID("") {
		return fmt.Errorf("kvtable create item id is null")
	}
	exist := r.HasItem(did)
	if exist == true {
		return fmt.Errorf("Key %s already existed in kvtable", did)
	}
	return r.setItem(did, item)
}

// UpdateItem checks and sets
func (r *KVTable) UpdateItem(item TableItem) error {
	did := item.GetID()
	if did == DID("") {
		return fmt.Errorf("kvtable create item id is null")
	}
	exist := r.HasItem(did)
	if exist == false {
		return fmt.Errorf("Key %s not existed in kvtable", did)
	}
	r.setItem(did, item)
	return nil
}

// GetItem checks ang gets
func (r *KVTable) GetItem(did DID, typ TableType) (TableItem, error) {
	exist := r.HasItem(did)
	if exist == false {
		return nil, fmt.Errorf("Key %s not existed in kvtable", did)
	}
	itemBytes := r.store.Get([]byte(did))
	switch typ {
	case DIDTableType:
		di := &DIDItem{}
		err := di.Unmarshal(itemBytes)
		if err != nil {
			return nil, fmt.Errorf("kvtable unmarshal did item: %w", err)
		}
		return di, nil
	case MethodTableType:
		mi := &MethodItem{}
		err := mi.Unmarshal(itemBytes)
		if err != nil {
			return nil, fmt.Errorf("kvtable unmarshal method item: %w", err)
		}
		return mi, nil
	default:
		return nil, fmt.Errorf("kvtable unknown table type: %d", typ)
	}
}

// DeleteItem without any checks
func (r *KVTable) DeleteItem(did DID) {
	r.store.Delete([]byte(did))
}

// Close .
func (r *KVTable) Close() error {
	err := r.store.Close()
	if err != nil {
		return fmt.Errorf("kvtable store: %w", err)
	}
	return nil
}
