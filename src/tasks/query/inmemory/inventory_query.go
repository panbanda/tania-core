package inmemory

import (
	"github.com/Tanibox/tania-server/src/assets/storage"
	"github.com/Tanibox/tania-server/src/tasks/query"
	uuid "github.com/satori/go.uuid"
)

type MaterialQueryInMemory struct {
	Storage *storage.MaterialStorage
}

func NewMaterialQueryInMemory(s *storage.MaterialStorage) query.MaterialQuery {
	return MaterialQueryInMemory{Storage: s}
}

func (s MaterialQueryInMemory) FindMaterialByID(inventoryUID uuid.UUID) <-chan query.QueryResult {
	result := make(chan query.QueryResult)

	go func() {
		s.Storage.Lock.RLock()
		defer s.Storage.Lock.RUnlock()

		ci := query.TaskMaterialQueryResult{}
		for _, val := range s.Storage.MaterialMap {
			// WARNING, domain leakage

			if val.UID == inventoryUID {
				ci.UID = val.UID
				ci.Name = val.Name
			}
		}

		result <- query.QueryResult{Result: ci}

		close(result)
	}()

	return result
}