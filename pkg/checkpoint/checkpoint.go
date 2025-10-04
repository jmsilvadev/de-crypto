package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
)

type Checkpoint struct {
	Confirmed uint64 `json:"confirmed"`
}

type CheckpointStore struct {
	path string
	mu   sync.Mutex
}

func NewCheckpointStore(path string) *CheckpointStore {
	return &CheckpointStore{path: path}
}

func (s *CheckpointStore) Load() (uint64, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		// lets try the tmp before send an error
		data, err = os.ReadFile(s.path + ".tmp")
		if err != nil {
			return 0, err
		}
	}

	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return 0, err
	}
	return cp.Confirmed, nil
}

func (s *CheckpointStore) Save(confirmed uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// lets be sure we will not corrupt the existent file
	tmp := s.path + ".tmp"

	data, err := json.Marshal(&Checkpoint{Confirmed: confirmed})
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, s.path)
}
