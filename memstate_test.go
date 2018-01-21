package oauth2state

import (
	"strconv"
	"sync"
	"testing"
)

type StubValueGenerator struct {
	count int
}

func (s StubValueGenerator) String() string {
	curr := s.count
	s.count++
	return strconv.Itoa(curr)
}

func TestNewMemStateStore(t *testing.T) {
	got := NewMemStateStore()

	if got == nil {
		t.Error("NewMemStateStore() returned nil")
	}
}

func TestMemStateStore_NewState(t *testing.T) {
	store := NewMemStateStore()
	store.valueGenerator = StubValueGenerator{256}

	url := "http://test.local"
	got, err := store.NewState(url)
	if err != nil {
		t.Errorf(`NewState(%q) failed on memory store with err: %v`, url, err)
	}
	if got != "256" {
		t.Errorf(`NewState(%q) = %q, want "256"`, url, got)
	}
}

func TestMemStateStore_Add(t *testing.T) {
	store := NewMemStateStore()
	store.valueGenerator = StubValueGenerator{}

	value := "A123456789012345678901234567890Z"
	wantURL := "http://test.local"
	err := store.Add(value, wantURL)
	if err != nil {
		t.Errorf(`Add(%q, %q) failed on memory store with err: %v`, value, wantURL, err)
	}
	got, gotOk := store.states[value]
	if !gotOk {
		t.Errorf(`Add(%q, %q) did not add the value to the map`, value, wantURL)
	}
	if got != wantURL {
		t.Errorf(`Add(%q, %q) did not add the url to the map. got %q, want %q`, value, wantURL, got, wantURL)
	}

	err = store.Add("", wantURL)
	if err == nil {
		t.Errorf(`Add("", %q) with empty value should have failed`, wantURL)
	}
}

func TestMemStateStore_Contains(t *testing.T) {
	store := NewMemStateStore()
	store.valueGenerator = StubValueGenerator{}
	store.Add("42", "http://test.local")

	ok, err := store.Contains("42")
	if err != nil {
		t.Errorf(`Contains("42") failed on memory store with err: %v`, err)
	}
	if !ok {
		t.Errorf(`Contains("42") = %v, want %v`, ok, true)
	}

	_, err = store.Contains("")
	if err == nil {
		t.Error(`Contains("") with empty value should have failed`)
	}
}

func TestMemStateStore_URL(t *testing.T) {
	store := NewMemStateStore()
	store.valueGenerator = StubValueGenerator{42}
	want := "http://test.local"

	got, err := store.URL("NotExistingValue")
	if err == nil {
		t.Errorf(`URL(%q) found an URL: %q`, "NotExistingValue", got)
	}

	value, _ := store.NewState(want)

	got, err = store.URL(value)
	if err != nil {
		t.Errorf(`URL(%q) failed on memory store with err: %v`, value, err)
	}
	if got != want {
		t.Errorf(`URL(%q) = %q, want %v`, value, got, want)
	}

	_, err = store.URL("")
	if err == nil {
		t.Error(`URL("") with empty value should have failed`)
	}
}

func TestMemStateStore_Delete(t *testing.T) {
	store := NewMemStateStore()
	store.valueGenerator = StubValueGenerator{42}
	url := "http://test.local"

	ok, _ := store.Contains("NotExistingValue")
	if ok {
		t.Errorf(`Contains(%q) found value`, "NotExistingValue")
	}

	value, _ := store.NewState(url)

	ok, _ = store.Contains(value)
	if !ok {
		t.Errorf(`Contains(%q) could not find value`, value)
	}

	store.Delete(value)

	ok, _ = store.Contains(value)
	if ok {
		t.Errorf(`Contains(%q) found deleted value`, value)
	}

	err := store.Delete("")
	if err == nil {
		t.Error(`Delete("") with empty value should have failed`)
	}
}

func BenchmarkRace(b *testing.B) {
	store := NewMemStateStore()

	var wg sync.WaitGroup

	loops := 100
	for i := 1; i <= loops; i++ {
		wg.Add(1)
		go func() {
			value, _ := store.NewState("http://test.local")
			store.Contains(value)
			store.URL(value)
			store.Delete(value)
			defer wg.Done()
		}()
	}
	wg.Wait()
}
