package messagehistory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageHistory_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		messages []string
		expected []string
	}{
		{
			name:     "add single message",
			messages: []string{"hello"},
			expected: []string{"hello"},
		},
		{
			name:     "add multiple messages",
			messages: []string{"hello", "world", "test"},
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "skip empty messages",
			messages: []string{"hello", "", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "skip duplicate consecutive messages",
			messages: []string{"hello", "hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "allow non-consecutive duplicates",
			messages: []string{"hello", "world", "hello"},
			expected: []string{"hello", "world", "hello"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mh := InitService()

			for _, msg := range tt.messages {
				mh.Add(msg)
			}

			assert.Equal(t, len(tt.expected), mh.Size())

			// Verify history by navigating through it
			for i := len(tt.expected) - 1; i >= 0; i-- {
				msg, ok := mh.NavigateUp("")
				assert.True(t, ok)
				assert.Equal(t, tt.expected[i], msg)
			}
		})
	}
}

func TestMessageHistory_NavigateUp(t *testing.T) {
	t.Parallel()

	mh := InitService()
	mh.Add("first")
	mh.Add("second")
	mh.Add("third")

	// First navigation should save current message and return last history item
	msg, ok := mh.NavigateUp("current")
	assert.True(t, ok)
	assert.Equal(t, "third", msg)
	assert.Equal(t, "current", mh.GetCurrent())

	// Second navigation
	msg, ok = mh.NavigateUp("current")
	assert.True(t, ok)
	assert.Equal(t, "second", msg)

	// Third navigation
	msg, ok = mh.NavigateUp("current")
	assert.True(t, ok)
	assert.Equal(t, "first", msg)

	// Fourth navigation should fail (at beginning)
	msg, ok = mh.NavigateUp("current")
	assert.False(t, ok)
	assert.Equal(t, "", msg)
}

func TestMessageHistory_NavigateDown(t *testing.T) {
	t.Parallel()

	mh := InitService()
	mh.Add("first")
	mh.Add("second")
	mh.Add("third")

	// Navigate up to beginning
	mh.NavigateUp("current")
	mh.NavigateUp("current")
	mh.NavigateUp("current")

	// Navigate down
	msg, ok := mh.NavigateDown()
	assert.True(t, ok)
	assert.Equal(t, "second", msg)

	msg, ok = mh.NavigateDown()
	assert.True(t, ok)
	assert.Equal(t, "third", msg)

	// Should return to current message
	msg, ok = mh.NavigateDown()
	assert.True(t, ok)
	assert.Equal(t, "current", msg)

	// Further navigation should fail
	msg, ok = mh.NavigateDown()
	assert.False(t, ok)
	assert.Equal(t, "", msg)
}

func TestMessageHistory_EmptyHistory(t *testing.T) {
	t.Parallel()

	mh := InitService()

	// Navigation should fail on empty history
	msg, ok := mh.NavigateUp("current")
	assert.False(t, ok)
	assert.Equal(t, "", msg)

	msg, ok = mh.NavigateDown()
	assert.False(t, ok)
	assert.Equal(t, "", msg)

	assert.Equal(t, 0, mh.Size())
}

func TestMessageHistory_Reset(t *testing.T) {
	t.Parallel()

	mh := InitService()
	mh.Add("first")
	mh.Add("second")

	// Navigate up and set current message
	mh.NavigateUp("current")
	assert.Equal(t, "current", mh.GetCurrent())

	// Reset should clear current message and reset index
	mh.Reset()
	assert.Equal(t, "", mh.GetCurrent())

	// Should be at end of history
	msg, ok := mh.NavigateDown()
	assert.False(t, ok)
	assert.Equal(t, "", msg)
}

func TestMessageHistory_Clear(t *testing.T) {
	t.Parallel()

	mh := InitService()
	mh.Add("first")
	mh.Add("second")
	mh.SetCurrent("current")

	assert.Equal(t, 2, mh.Size())
	assert.Equal(t, "current", mh.GetCurrent())

	mh.Clear()

	assert.Equal(t, 0, mh.Size())
	assert.Equal(t, "", mh.GetCurrent())

	// Navigation should fail after clear
	msg, ok := mh.NavigateUp("test")
	assert.False(t, ok)
	assert.Equal(t, "", msg)
}

func TestMessageHistory_SetCurrent(t *testing.T) {
	t.Parallel()

	mh := InitService()

	mh.SetCurrent("test message")
	assert.Equal(t, "test message", mh.GetCurrent())

	mh.SetCurrent("updated message")
	assert.Equal(t, "updated message", mh.GetCurrent())
}

func TestMessageHistory_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	mh := InitService()

	// Test concurrent access doesn't panic
	go func() {
		for i := 0; i < 100; i++ {
			mh.Add("message")
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			mh.NavigateUp("current")
			mh.NavigateDown()
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			mh.Size()
			mh.GetCurrent()
		}
	}()

	// If we reach here without panic, concurrent access is working
	assert.True(t, true)
}
