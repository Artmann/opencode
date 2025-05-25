# Message History Persistence Plan

## Current State Analysis

The `messagehistory` package currently provides an in-memory service for managing user input history with navigation capabilities (up/down arrows). The service includes:

- **In-memory storage**: History stored in a slice (`[]string`)
- **Navigation**: Up/down navigation through history
- **Deduplication**: Prevents duplicate consecutive entries
- **Thread-safe**: Uses `sync.RWMutex` for concurrent access
- **Global singleton**: Single instance via `InitService()`

## Problem Statement

Message history is lost between application runs, requiring users to re-enter previously used commands/messages.

## Proposed Solution

### 1. Database Integration

**Add new table for message history:**

```sql
CREATE TABLE IF NOT EXISTS message_history (
    id TEXT PRIMARY KEY,
    session_id TEXT REFERENCES sessions(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f000Z', 'now')),
    UNIQUE(session_id, message)
);

CREATE INDEX message_history_session_id_idx ON message_history(session_id);
CREATE INDEX message_history_created_at_idx ON message_history(created_at);
```

**Benefits:**

- Leverages existing database infrastructure
- Automatic cleanup when sessions are deleted
- Per-session history isolation
- Prevents duplicate messages per session

### 2. Service Architecture Changes

**Hybrid approach: In-memory + Database**

- Keep current in-memory implementation for performance
- Add database persistence layer
- Load history from database on initialization
- Save new entries to database asynchronously

**Interface Extensions:**

```go
type Service interface {
    // Existing methods
    Add(message string)
    Clear()
    GetCurrent() string
    NavigateDown() (string, bool)
    NavigateUp(currentMessage string) (string, bool)
    Reset()
    SetCurrent(message string)
    Size() int

    // New methods
    LoadFromSession(ctx context.Context, sessionID string) error
    SaveToSession(ctx context.Context, sessionID string) error
}
```

### 3. Implementation Strategy

**Phase 1: Database Schema**

1. Add migration for `message_history` table
2. Generate SQLC queries for CRUD operations

**Phase 2: Service Enhancement**

1. Add database dependency to service
2. Implement `LoadFromSession()` to populate in-memory history
3. Implement `SaveToSession()` for persistence

**Phase 3: Integration**

1. Modify `InitService()` to accept database connection
2. Load history when session starts
3. Save history on message add (async)
4. Handle session switching

### 4. Configuration Options

**Performance Considerations:**

- Batch database writes for better performance
- Lazy loading of history (load on first access)
- Option to disable persistence (fallback to in-memory only)

### 5. Migration Strategy

**Backward Compatibility:**

- Existing in-memory functionality remains unchanged
- Database persistence is additive
- Graceful degradation if database is unavailable

**Data Migration:**

- No existing data to migrate (currently in-memory only)
- New installations start with empty history

### 6. Testing Strategy

**Unit Tests:**

- Test in-memory functionality (existing)
- Test database persistence operations
- Test session isolation
- Test error handling (database unavailable)

**Integration Tests:**

- Test full workflow: save → restart → load
- Test session switching

### 7. Alternative Approaches Considered

**File-based storage:**

- Pros: Simple, no database dependency
- Cons: Manual file management, no session isolation, harder cleanup

**Global history (not per-session):**

- Pros: Simpler implementation
- Cons: No isolation, potential confusion across different contexts

**Selected approach (database + per-session) provides the best balance of functionality, performance, and maintainability.**

## Success Criteria

- Message history persists across application restarts
- History is isolated per session
- Performance impact is minimal (< 10ms for typical operations)
- Backward compatibility maintained
- Comprehensive test coverage (> 90%)
