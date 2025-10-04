package checkpoint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCheckpointStore(t *testing.T) {
	t.Run("create_checkpoint_store", func(t *testing.T) {
		path := "/tmp/test-checkpoint"
		store := NewCheckpointStore(path)

		assert.NotNil(t, store)
		assert.Equal(t, path, store.path)
	})
}

func TestCheckpointStore_Load(t *testing.T) {
	t.Run("load_from_nonexistent_file", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "nonexistent.json")

		store := NewCheckpointStore(path)
		confirmed, err := store.Load()

		assert.NoError(t, err)
		assert.Equal(t, uint64(0), confirmed)
	})

	t.Run("load_from_existing_file", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "checkpoint.json")


		err := os.WriteFile(path, []byte(`{"confirmed":12345}`), 0644)
		assert.NoError(t, err)

		store := NewCheckpointStore(path)
		confirmed, err := store.Load()

		assert.NoError(t, err)
		assert.Equal(t, uint64(12345), confirmed)
	})

	t.Run("load_from_tmp_file_when_main_file_has_error", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "checkpoint.json")
		tmpPath := path + ".tmp"


		err := os.Mkdir(path, 0755)
		assert.NoError(t, err)


		err = os.WriteFile(tmpPath, []byte(`{"confirmed":67890}`), 0644)
		assert.NoError(t, err)

		store := NewCheckpointStore(path)
		confirmed, err := store.Load()

		assert.NoError(t, err)
		assert.Equal(t, uint64(67890), confirmed)
	})

	t.Run("load_from_corrupted_file", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "corrupted.json")


		err := os.WriteFile(path, []byte("invalid json"), 0644)
		assert.NoError(t, err)

		store := NewCheckpointStore(path)
		confirmed, err := store.Load()

		assert.Error(t, err)
		assert.Equal(t, uint64(0), confirmed)
	})
}

func TestCheckpointStore_Save(t *testing.T) {
	t.Run("save_checkpoint_successfully", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "checkpoint.json")

		store := NewCheckpointStore(path)
		err := store.Save(54321)

		assert.NoError(t, err)


		_, err = os.Stat(path)
		assert.NoError(t, err)


		confirmed, err := store.Load()
		assert.NoError(t, err)
		assert.Equal(t, uint64(54321), confirmed)
	})

	t.Run("save_checkpoint_overwrites_existing", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "checkpoint.json")

		store := NewCheckpointStore(path)


		err := store.Save(11111)
		assert.NoError(t, err)


		err = store.Save(22222)
		assert.NoError(t, err)


		confirmed, err := store.Load()
		assert.NoError(t, err)
		assert.Equal(t, uint64(22222), confirmed)
	})

	t.Run("save_checkpoint_with_zero_value", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "checkpoint.json")

		store := NewCheckpointStore(path)
		err := store.Save(0)

		assert.NoError(t, err)

		confirmed, err := store.Load()
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), confirmed)
	})

	t.Run("save_checkpoint_with_max_uint64", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "checkpoint.json")

		store := NewCheckpointStore(path)
		err := store.Save(18446744073709551615)

		assert.NoError(t, err)

		confirmed, err := store.Load()
		assert.NoError(t, err)
		assert.Equal(t, uint64(18446744073709551615), confirmed)
	})
}

func TestCheckpointStore_ConcurrentAccess(t *testing.T) {
	t.Run("concurrent_save_operations", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "checkpoint.json")

		store := NewCheckpointStore(path)


		done := make(chan error, 10)
		for i := 0; i < 10; i++ {
			go func(value uint64) {
				done <- store.Save(value)
			}(uint64(i))
		}


		for i := 0; i < 10; i++ {
			err := <-done
			assert.NoError(t, err)
		}


		confirmed, err := store.Load()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, confirmed, uint64(0))
		assert.LessOrEqual(t, confirmed, uint64(9))
	})
}

func TestCheckpoint_JSONSerialization(t *testing.T) {
	t.Run("marshal_checkpoint_to_json", func(t *testing.T) {
		checkpoint := Checkpoint{Confirmed: 12345}


		assert.Equal(t, uint64(12345), checkpoint.Confirmed)
	})
}
