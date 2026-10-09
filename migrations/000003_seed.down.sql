BEGIN;

DELETE FROM stock_batches;
DELETE FROM recipe_components;
DELETE FROM recipes;
DELETE FROM items;

COMMIT;
