package address

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type addrRecord struct {
	UserID  string `json:"userId"`
	Address string `json:"address"`
}

type MemoryAddressIndex struct {
	data map[string]string
}

func NewMemoryAddressIndexFromJSON(path string) (*MemoryAddressIndex, error) {
	log.Println("Loading address index from", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var recs []addrRecord
	if err := json.NewDecoder(f).Decode(&recs); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	log.Println("File loaded")

	data := make(map[string]string, len(recs))
	for i, r := range recs {
		a := strings.ToLower(strings.TrimSpace(r.Address))
		if a == "" {
			return nil, fmt.Errorf("empty address at index %d", i)
		}

		if len(a) != 42 || !strings.HasPrefix(a, "0x") {
			return nil, fmt.Errorf("invalid address at index %d: %q", i, r.Address)
		}
		data[a] = r.UserID
	}

	return &MemoryAddressIndex{data: data}, nil
}

func (m *MemoryAddressIndex) Lookup(addr string) (string, bool) {
	uid, ok := m.data[strings.ToLower(addr)]
	return uid, ok
}
