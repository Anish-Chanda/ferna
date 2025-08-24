-- sqlite3 compatible table defs
-- Enums are represented as TEXT with CHECK constraints because sqlite3 has no native enum type.
-- datetime/timestamp columns are stored as TEXT in ISO-8601 (DEFAULT CURRENT_TIMESTAMP).

PRAGMA foreign_keys = ON;

-- users
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL COLLATE NOCASE UNIQUE,
  password_hash TEXT,
  auth_provider TEXT NOT NULL CHECK(auth_provider IN ('local','google')),
  provider_user_id TEXT, -- e.g., Google sub (nullable)
  created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  CHECK (auth_provider != 'local' OR password_hash IS NOT NULL),
  UNIQUE(auth_provider, provider_user_id)
);

-- species
CREATE TABLE IF NOT EXISTS species (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  common_name TEXT NOT NULL,
  scientific_name TEXT,
  light_pref TEXT NOT NULL CHECK(light_pref IN ('low','medium','bright_indirect','bright_direct')),
  default_water_interval_days INTEGER NOT NULL CHECK(default_water_interval_days >= 1),
  default_fertilizer_interval_days INTEGER NOT NULL CHECK(default_fertilizer_interval_days >= 1),
  toxicity TEXT NOT NULL CHECK(toxicity IN ('toxic_to_pets','non_toxic')),
  care_notes TEXT,
  care_notes_source TEXT,
  created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

-- location (user-scoped)
CREATE TABLE IF NOT EXISTS location (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  UNIQUE(user_id, name),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- user_plants
CREATE TABLE IF NOT EXISTS user_plants (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  species_id INTEGER NOT NULL,
  nickname TEXT,
  image_url TEXT,
  notes TEXT,
  water_interval_days_override INTEGER,
  fertilizer_interval_days_override INTEGER,
  location_id INTEGER,
  created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (species_id) REFERENCES species(id) ON DELETE RESTRICT,
  FOREIGN KEY (location_id) REFERENCES location(id) ON DELETE SET NULL
);

-- plant_tasks
CREATE TABLE IF NOT EXISTS plant_tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  plant_id INTEGER NOT NULL,
  task_type TEXT NOT NULL CHECK(task_type IN ('watering','fertilizer')),
  snoozed_until TEXT,
  interval_days INTEGER NOT NULL CHECK(interval_days >= 1),
  tolerance_days INTEGER NOT NULL CHECK(tolerance_days >= 0),
  next_due_at TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  UNIQUE(plant_id, task_type),
  FOREIGN KEY (plant_id) REFERENCES user_plants(id) ON DELETE CASCADE
);

-- care_events
CREATE TABLE IF NOT EXISTS care_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  plant_id INTEGER NOT NULL,
  task_id INTEGER,
  event_type TEXT NOT NULL CHECK(event_type IN ('watering','fertilizer','repotting','pruning','other')),
  happened_at TEXT NOT NULL,
  notes TEXT,
  created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY (plant_id) REFERENCES user_plants(id) ON DELETE CASCADE,
  FOREIGN KEY (task_id) REFERENCES plant_tasks(id) ON DELETE SET NULL
);

-- updated_at triggers (SQLite pattern: AFTER UPDATE with guard)
CREATE TRIGGER IF NOT EXISTS trg_users_updated_at
AFTER UPDATE ON users FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
  UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_species_updated_at
AFTER UPDATE ON species FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
  UPDATE species SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_location_updated_at
AFTER UPDATE ON location FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
  UPDATE location SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_user_plants_updated_at
AFTER UPDATE ON user_plants FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
  UPDATE user_plants SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_plant_tasks_updated_at
AFTER UPDATE ON plant_tasks FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
  UPDATE plant_tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_care_events_updated_at
AFTER UPDATE ON care_events FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
  UPDATE care_events SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
