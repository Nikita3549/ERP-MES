BEGIN;

DROP TRIGGER IF EXISTS recipe_components_no_cycle ON recipe_components;
DROP FUNCTION IF EXISTS recipe_components_forbid_cycle();

DROP TABLE IF EXISTS recipe_components;
DROP TABLE IF EXISTS recipes;
DROP TYPE IF EXISTS recipe_status;

DROP TABLE IF EXISTS items;
DROP TYPE IF EXISTS item_type;
DROP TYPE IF EXISTS item_unit;

COMMIT;
