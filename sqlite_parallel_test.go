package go_sqlite3_locked_tests_test

import (
	"runtime"
	"github.com/elgs/gostrgen"
	"github.com/comwrg/go-sqlite3-locked-tests"
	"testing"
	"sync"
)

const TestNum = 10000

func Select(s *go_sqlite3_locked_tests.Sqlite, t *testing.T) {
	_, err := s.Select(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsertChanAndSelectParallel(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := make(chan string)
	s := go_sqlite3_locked_tests.Sqlite{}
	wg := sync.WaitGroup{}

	err := s.Init("TestInsertChanAndSelectParallel")
	if err != nil {
		t.Fatal("sqlite open failed.", err)
	}
	defer s.Close()
	go func() {
		for {
			select {
			case m := <-c:
				err := s.Insert(m)
				if err != nil {
					t.Fail()
				}
				wg.Done()
			default:
				runtime.Gosched()
			}
		}
	}()
	wg.Add(TestNum)
	for i := 0; i < TestNum; i++ {
		go func() {
			str, _ := gostrgen.RandGen(50, gostrgen.All, "", "")
			c <- str
		}()
		//////////////////
		//go Select(&s, t)
		//////////////////
	}
	wg.Wait()
}

func TestInsertMutexAndSelectParallel(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	mux := sync.Mutex{}
	s := go_sqlite3_locked_tests.Sqlite{}
	wg := sync.WaitGroup{}

	err := s.Init("TestInsertMutexAndSelectParallel")
	if err != nil {
		t.Fatal("sqlite init failed.")
	}
	defer s.Close()
	wg.Add(TestNum)
	for i := 0; i < TestNum; i++ {
		go func() {
			str, _ := gostrgen.RandGen(50, gostrgen.All, "", "")
			mux.Lock()
			err = s.Insert(str)
			mux.Unlock()
			if err != nil {
			}
			wg.Done()
		}()
		//////////////////
		//go Select(&s, t)
		//////////////////
	}
	wg.Wait()
}

