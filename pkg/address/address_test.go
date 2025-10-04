package address

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryAddressIndexFromJSON_ValidFile(t *testing.T) {
	t.Run("create_memory_address_index_from_valid_json", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "addresses.json")

		addresses := []addrRecord{
			{UserID: "user1", Address: "0x1234567890123456789012345678901234567890"},
			{UserID: "user2", Address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"},
		}

		data, err := json.Marshal(addresses)
		assert.NoError(t, err)

		err = os.WriteFile(jsonFile, data, 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.NoError(t, err)
		assert.NotNil(t, index)

		userID, ok := index.Lookup("0x1234567890123456789012345678901234567890")
		assert.True(t, ok)
		assert.Equal(t, "user1", userID)

		userID, ok = index.Lookup("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
		assert.True(t, ok)
		assert.Equal(t, "user2", userID)
	})

	t.Run("create_memory_address_index_with_case_insensitive_lookup", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "addresses.json")

		addresses := []addrRecord{
			{UserID: "user1", Address: "0x1234567890123456789012345678901234567890"},
		}

		data, err := json.Marshal(addresses)
		assert.NoError(t, err)

		err = os.WriteFile(jsonFile, data, 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.NoError(t, err)

		userID, ok := index.Lookup("0X1234567890123456789012345678901234567890")
		assert.True(t, ok)
		assert.Equal(t, "user1", userID)
	})
}

func TestNewMemoryAddressIndexFromJSON_InvalidFile(t *testing.T) {
	t.Run("create_memory_address_index_from_nonexistent_file", func(t *testing.T) {
		index, err := NewMemoryAddressIndexFromJSON("nonexistent.json")
		assert.Error(t, err)
		assert.Nil(t, index)
	})

	t.Run("create_memory_address_index_from_invalid_json", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "invalid.json")

		err := os.WriteFile(jsonFile, []byte("invalid json"), 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.Error(t, err)
		assert.Nil(t, index)
	})

	t.Run("create_memory_address_index_with_empty_address", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "addresses.json")

		addresses := []addrRecord{
			{UserID: "user1", Address: ""},
		}

		data, err := json.Marshal(addresses)
		assert.NoError(t, err)

		err = os.WriteFile(jsonFile, data, 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.Error(t, err)
		assert.Nil(t, index)
	})

	t.Run("create_memory_address_index_with_invalid_address_length", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "addresses.json")

		addresses := []addrRecord{
			{UserID: "user1", Address: "0x123"}, // Too short.
		}

		data, err := json.Marshal(addresses)
		assert.NoError(t, err)

		err = os.WriteFile(jsonFile, data, 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.Error(t, err)
		assert.Nil(t, index)
	})

	t.Run("create_memory_address_index_with_invalid_address_prefix", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "addresses.json")

		addresses := []addrRecord{
			{UserID: "user1", Address: "1234567890123456789012345678901234567890"}, // No 0x prefix.
		}

		data, err := json.Marshal(addresses)
		assert.NoError(t, err)

		err = os.WriteFile(jsonFile, data, 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.Error(t, err)
		assert.Nil(t, index)
	})
}

func TestMemoryAddressIndex_Lookup(t *testing.T) {
	t.Run("lookup_existing_address", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "addresses.json")

		addresses := []addrRecord{
			{UserID: "user1", Address: "0x1234567890123456789012345678901234567890"},
		}

		data, err := json.Marshal(addresses)
		assert.NoError(t, err)

		err = os.WriteFile(jsonFile, data, 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.NoError(t, err)

		userID, ok := index.Lookup("0x1234567890123456789012345678901234567890")
		assert.True(t, ok)
		assert.Equal(t, "user1", userID)
	})

	t.Run("lookup_nonexistent_address", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "addresses.json")

		addresses := []addrRecord{
			{UserID: "user1", Address: "0x1234567890123456789012345678901234567890"},
		}

		data, err := json.Marshal(addresses)
		assert.NoError(t, err)

		err = os.WriteFile(jsonFile, data, 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.NoError(t, err)

		userID, ok := index.Lookup("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
		assert.False(t, ok)
		assert.Empty(t, userID)
	})

	t.Run("lookup_with_whitespace_trimming", func(t *testing.T) {

		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "addresses.json")

		addresses := []addrRecord{
			{UserID: "user1", Address: "  0x1234567890123456789012345678901234567890  "},
		}

		data, err := json.Marshal(addresses)
		assert.NoError(t, err)

		err = os.WriteFile(jsonFile, data, 0644)
		assert.NoError(t, err)

		index, err := NewMemoryAddressIndexFromJSON(jsonFile)
		assert.NoError(t, err)

		userID, ok := index.Lookup("0x1234567890123456789012345678901234567890")
		assert.True(t, ok)
		assert.Equal(t, "user1", userID)
	})
}
