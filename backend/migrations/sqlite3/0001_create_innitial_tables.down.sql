PRAGMA foreign_keys = ON;
BEGIN TRANSACTION;

-- Drop triggers first
DROP TRIGGER IF EXISTS trg_care_events_updated_at;
DROP TRIGGER IF EXISTS trg_plant_tasks_updated_at;
DROP TRIGGER IF EXISTS trg_user_plants_updated_at;
DROP TRIGGER IF EXISTS trg_location_updated_at;
DROP TRIGGER IF EXISTS trg_species_updated_at;
DROP TRIGGER IF EXISTS trg_users_updated_at;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS care_events;
DROP TABLE IF EXISTS plant_tasks;
DROP TABLE IF EXISTS user_plants;
DROP TABLE IF EXISTS location;
DROP TABLE IF EXISTS species;
DROP TABLE IF EXISTS users;

COMMIT;
PRAGMA foreign_keys = ON;