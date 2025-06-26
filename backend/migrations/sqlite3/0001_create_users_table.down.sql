-- Drop triggers
DROP TRIGGER IF EXISTS update_species_updated_at;
DROP TRIGGER IF EXISTS update_users_updated_at;
DROP TRIGGER IF EXISTS update_plants_updated_at;

-- Drop indexes
DROP INDEX IF EXISTS idx_species_common_name;
DROP INDEX IF EXISTS idx_species_scientific_name;

-- Drop tables
DROP TABLE IF EXISTS species;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS plants;