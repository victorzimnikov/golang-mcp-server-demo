-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE RESTRICT,
  title TEXT NOT NULL CHECK (trim(title) <> ''),
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'todo' CHECK (
    status IN (
      'todo',
      'in_progress',
      'blocked',
      'done',
      'cancelled'
    )
  ),
  priority TEXT NOT NULL DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
  source TEXT NOT NULL CHECK (source IN ('claude', 'codex', 'human')),
  version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX tasks_project_id_idx ON tasks(project_id, status, created_at, id);

CREATE TABLE activity_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE RESTRICT,
  entity_type TEXT NOT NULL CHECK (
    entity_type IN ('project', 'task', 'decision', 'note')
  ),
  entity_id INTEGER NOT NULL CHECK (entity_id > 0),
  event_type TEXT NOT NULL CHECK (trim(event_type) <> ''),
  source TEXT NOT NULL CHECK (source IN ('claude', 'codex', 'human')),
  payload_json TEXT NOT NULL DEFAULT '{}' CHECK (json_valid(payload_json)),
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX activity_events_project_id_idx ON activity_events(project_id, created_at DESC, id DESC);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activity_events;

DROP TABLE IF EXISTS tasks;

-- +goose StatementEnd