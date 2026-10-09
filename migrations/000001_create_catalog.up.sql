BEGIN;

CREATE TYPE item_unit AS ENUM ('kg', 'l', 'pcs');
CREATE TYPE item_type AS ENUM ('raw_material', 'semi_finished', 'finished');

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number TEXT NOT NULL,
    name TEXT NOT NULL,
    unit item_unit NOT NULL,
    type item_type NOT NULL,

    CONSTRAINT items_id_type_key UNIQUE (id, type)
);

CREATE UNIQUE INDEX items_number_lower_key ON items (lower(number));

CREATE TYPE recipe_status AS ENUM ('draft', 'published');

CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL,
    item_type item_type NOT NULL,
    version INT NOT NULL CHECK (version > 0),
    status recipe_status NOT NULL DEFAULT 'draft',
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT recipes_items_fk FOREIGN KEY (item_id, item_type) REFERENCES items (id, type),
    CONSTRAINT recipes_not_raw CHECK (item_type <> 'raw_material'),
    CONSTRAINT recipes_item_version_key UNIQUE (item_id, version),
    CONSTRAINT recipes_active_is_published CHECK (NOT is_active OR status = 'published'),
    CONSTRAINT recipes_id_item_key UNIQUE (id, item_id)
);

CREATE UNIQUE INDEX recipes_one_active_per_item ON recipes (item_id) WHERE is_active;

CREATE TABLE recipe_components (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id UUID NOT NULL,
    item_id UUID NOT NULL,
    component_item_id UUID NOT NULL REFERENCES items (id),
    quantity NUMERIC(18, 6) NOT NULL CHECK (quantity > 0),

    CONSTRAINT recipe_components_recipe_fk FOREIGN KEY (recipe_id, item_id) REFERENCES recipes (
        id, item_id
    ) ON DELETE CASCADE,
    CONSTRAINT recipe_components_recipe_item_key UNIQUE (recipe_id, component_item_id),
    CONSTRAINT recipe_components_no_self_ref CHECK (component_item_id <> item_id)
);

CREATE INDEX recipe_components_graph_idx
ON recipe_components (item_id, component_item_id);

CREATE INDEX recipe_components_component_idx
ON recipe_components (component_item_id);

CREATE FUNCTION recipe_components_forbid_cycle() RETURNS TRIGGER AS $$
DECLARE
    has_cycle BOOLEAN;
BEGIN
    PERFORM pg_advisory_xact_lock(hashtext('recipe_graph'));

    WITH RECURSIVE reachable (item_id) AS (
        SELECT NEW.component_item_id
        UNION
        SELECT rc.component_item_id
        FROM recipe_components AS rc
        INNER JOIN reachable AS r ON rc.item_id = r.item_id
    )
    SELECT EXISTS (
        SELECT 1 FROM reachable WHERE reachable.item_id = NEW.item_id
    ) INTO has_cycle;

    IF has_cycle THEN
        RAISE EXCEPTION 'recipe cycle: item % is reachable from component %',
            NEW.item_id, NEW.component_item_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER recipe_components_no_cycle
BEFORE INSERT OR UPDATE OF item_id, component_item_id ON recipe_components
FOR EACH ROW EXECUTE FUNCTION recipe_components_forbid_cycle();

COMMIT;
