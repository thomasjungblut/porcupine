package rwlock

import (
	"github.com/anishathalye/porcupine"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
)

func TestHappyPath(t *testing.T) {
	db := NewDatabase()

	key := 42
	client := NewDatabaseRecorder(db, 0)
	for i := 0; i < 100; i++ {
		_, _ = client.Get(key)
		client.Put(key, i)
		if rand.Float32() < 0.25 {
			client.Del(key)
		}
	}

	verifyOperations(t, client.operations)
}

func TestMultiGoroutines(t *testing.T) {
	numGoRoutines := 4
	d := NewDatabase()

	var operations []porcupine.Operation
	var opsLock sync.Mutex
	wg := sync.WaitGroup{}

	for n := 0; n < numGoRoutines; n++ {
		wg.Add(1)
		go func(db *db, id int) {
			client := NewDatabaseRecorder(db, id)

			for j := 0; j < 100; j++ {
				for i := 0; i < 10; i++ {
					_, _ = client.Get(i)
					client.Put(i, i)
					if rand.Float32() < 0.25 {
						client.Del(i)
					}
				}
			}

			opsLock.Lock()
			defer opsLock.Unlock()

			operations = append(operations, client.operations...)

			wg.Done()
		}(d, n)
	}

	wg.Wait()

	verifyOperations(t, operations)
}

func verifyOperations(t *testing.T, operations []porcupine.Operation) {
	result, info := porcupine.CheckOperationsVerbose(Model, operations, 0)
	require.NoError(t, porcupine.VisualizePath(Model, info, t.Name()+"_porcupine.html"))
	require.Equal(t, porcupine.CheckResult(porcupine.Ok), result, "output was not linearizable")
}
