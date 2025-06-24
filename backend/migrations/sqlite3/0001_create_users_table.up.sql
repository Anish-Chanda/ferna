CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    pass_hash TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS species (
    id                             INTEGER PRIMARY KEY AUTOINCREMENT,
    common_name                    TEXT    NOT NULL,
    scientific_name                TEXT    UNIQUE,
    default_watering_frequency_days INTEGER NOT NULL DEFAULT 7,
    created_at                     DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at                     DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- triggers, updates updated_at for users and species tables 
CREATE TRIGGER IF NOT EXISTS update_species_updated_at
AFTER UPDATE ON species
FOR EACH ROW
BEGIN
    UPDATE species
      SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS update_users_updated_at
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
  UPDATE users
     SET updated_at = CURRENT_TIMESTAMP
   WHERE id = OLD.id;
END;

-- indexes
CREATE INDEX IF NOT EXISTS idx_species_common_name     ON species(common_name);
CREATE INDEX IF NOT EXISTS idx_species_scientific_name ON species(scientific_name);
