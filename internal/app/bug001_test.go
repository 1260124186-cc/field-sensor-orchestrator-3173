package app

import (
 "fmt"
 "sync"
 "testing"
)

func TestConcurrentCycleScheduling(t *testing.T) {
 application := NewDemo()
 const workers = 24
 ids := make(chan string, workers)
 errs := make(chan error, workers)
 var wg sync.WaitGroup
 for index := 0; index < workers; index++ {
  wg.Add(1)
  go func(index int) {
   defer wg.Done()
   cycle, err := application.Schedule(fmt.Sprintf("coordinator-%02d", index))
   if err != nil { errs <- err; return }
   ids <- cycle.ID
  }(index)
 }
 wg.Wait(); close(ids); close(errs)
 for err := range errs { t.Fatalf("schedule failed: %v", err) }
 seen := map[string]bool{}
 for id := range ids { if seen[id] { t.Fatalf("duplicate cycle id %s", id) }; seen[id] = true }
 if len(seen) != workers { t.Fatalf("scheduled=%d want=%d", len(seen), workers) }
}

func TestSequentialCycleScheduling(t *testing.T) {
 application := NewDemo()
 first, err := application.Schedule("one"); if err != nil { t.Fatal(err) }
 second, err := application.Schedule("two"); if err != nil { t.Fatal(err) }
 if first.ID == second.ID { t.Fatalf("sequential ids must differ: %s", first.ID) }
}
