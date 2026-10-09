# erp-mes

## Run

Docker (db → migrations → api):

    docker compose --profile api up --build

Local:

    docker compose up -d db
    go build -o bin/api ./cmd && ./bin/api

## Config

Env vars or `.env` (see `.env.example`). All optional:

| Var | Default |
|---|---|
| DB_HOST | localhost |
| DB_PORT | 5432 |
| DB_USER | postgres |
| DB_PASSWORD | postgres |
| DB_NAME | erp-mes |
| HTTP_PORT | 8080 |
| HTTP_READ_HEADER_TIMEOUT | 5s |
| HTTP_READ_TIMEOUT | 15s |
| HTTP_WRITE_TIMEOUT | 20s |
| HTTP_IDLE_TIMEOUT | 60s |
| HTTP_SHUTDOWN_TIMEOUT | 20s |

On SIGINT/SIGTERM the server drains requests for up to `HTTP_SHUTDOWN_TIMEOUT`, then closes the DB.
Keep it `>= HTTP_WRITE_TIMEOUT` and below compose `stop_grace_period` (25s).

## Endpoints

- `GET /health` — liveness, `{"status":"ok"}`
- `GET /ready` — readiness, pings DB; `503 {"status":"fail"}` on failure

## Schema

Postgres. Three concerns: a catalog of all items, versioned recipes (BOM) over that catalog, and
batch-level stock of raw materials. Invariants live in the schema — foreign keys, `CHECK`, partial
unique indexes, one trigger — not in the service layer, so they hold no matter how many API
consumers or transports appear later.

### Enums

| Type | Values | Purpose |
|---|---|---|
| `item_unit` | `kg`, `l`, `pcs` | unit of measure |
| `item_type` | `raw_material`, `semi_finished`, `finished` | role in production; also carried into child tables as part of composite FKs |
| `recipe_status` | `draft`, `published` | lifecycle of a recipe version |

### `items`

Single catalog for everything: raw materials, semi-finished goods, finished products. One table, not
three, because a recipe component may be of any type, and a semi-finished item is simultaneously a
product (it has a recipe) and a component (it appears in other recipes).

- `number` is the business key; `items_number_lower_key` is a unique index on `lower(number)`, so
  articles are unique case-insensitively.
- `UNIQUE (id, type)` is redundant next to the primary key on `id`, and exists only to be the target
  of composite foreign keys. Child tables keep their own `item_type` column and constrain it with a
  plain `CHECK`. That is how "no recipe for a raw material" and "no batch for a semi-finished item"
  are enforced in the database without triggers and without subqueries inside `CHECK`.

### `recipes`

Header of one BOM version: a row is one version of the recipe for one item.

- `UNIQUE (item_id, version)` with `CHECK (version > 0)` — versions are never reused, so history stays
  intact.
- `CHECK (item_type <> 'raw_material')` — raw materials cannot have recipes.
- `recipes_one_active_per_item` is a partial unique index `ON (item_id) WHERE is_active`: at most one
  active recipe per item, while every other version keeps living next to it. No service-side
  "deactivate the old one first" dance can produce two active rows.
- `CHECK (NOT is_active OR status = 'published')` — a draft cannot be made active.
- `UNIQUE (id, item_id)` is again a composite-FK target: recipe lines carry the parent's `item_id`.
- `created_at` for ordering versions and for audit.

### `recipe_components`

BOM lines: how much of a component goes into one unit of the parent item.

- `quantity NUMERIC(18, 6) CHECK (quantity > 0)` — numeric rather than float, because quantities are
  multiplied and summed while exploding a BOM and accumulated binary error is not acceptable.
- Composite FK `(recipe_id, item_id) → recipes (id, item_id)` with `ON DELETE CASCADE`: lines exist
  only together with their recipe version, and the denormalized parent `item_id` cannot drift away
  from the recipe it belongs to.
- `UNIQUE (recipe_id, component_item_id)` — a component appears at most once per recipe, so readers
  never have to sum duplicate lines.
- `CHECK (component_item_id <> item_id)` — an item is not a direct component of itself.
- Trigger `recipe_components_no_cycle` calls `recipe_components_forbid_cycle()`: a recursive
  reachability walk from the new component, raising if the parent is reachable. The
  `pg_advisory_xact_lock` serializes those checks — otherwise two concurrent transactions each insert
  one half of a cycle, and each half is valid on its own. Two things to know about it: the graph is
  built at *item* level, not per recipe version, so edges from all versions (including `draft`) count
  as existing — conservative, it may reject an edge that only closes a cycle through a draft, but it
  never lets a real cycle through; and it fires only on `INSERT` and
  `UPDATE OF item_id, component_item_id`.
- Indexes `(item_id, component_item_id)` and `(component_item_id)` serve the two directions of the
  graph: explosion downward, and where-used upward (which recipes contain this item).

### `stock_batches`

Batch stock of raw materials: receipt, expiry, reservation.

- `CHECK (item_type = 'raw_material')` — only raw materials have batches; receipt and output of
  semi-finished goods is not modelled yet.
- `initial_quantity > 0`, `reserved_quantity >= 0`, `CHECK (reserved_quantity <= initial_quantity)`.
- `remaining_quantity` is `GENERATED ALWAYS AS (initial_quantity - reserved_quantity) STORED`, so the
  remainder cannot drift out of sync — writes only ever touch `reserved_quantity`.
- `CHECK (expires_at >= received_at)`.
- `UNIQUE (item_id, supplier_batch_number)` makes receiving idempotent: re-posting the same supplier
  delivery note cannot create a second batch.
- `stock_batches_fifo_idx` is a partial index `ON (item_id, expires_at) WHERE reserved_quantity <
  initial_quantity` — picking by FEFO (earliest expiry first). Exhausted batches drop out of the
  index, so it does not grow with history.

### Relations

| FK | Target | Why this shape |
|---|---|---|
| `recipes (item_id, item_type)` | `items (id, type)` | type is carried into the child so a `CHECK` can forbid recipes for raw materials |
| `recipe_components (recipe_id, item_id)` | `recipes (id, item_id)`, `ON DELETE CASCADE` | a line cannot reference a recipe of a different item; lines die with their version |
| `recipe_components (component_item_id)` | `items (id)` | a component is any catalog item |
| `stock_batches (item_id, item_type)` | `items (id, type)` | type is carried into the batch so a `CHECK` leaves raw materials only |

Foreign keys to `items` have no `ON DELETE` clause: an item referenced by a recipe or a batch cannot
be deleted (`NO ACTION`). That is deliberate — catalog rows are retired, not removed.

### Seed (`000003_seed`)

A chocolate-bar plant: 12 items (6 raw, 3 semi-finished, 3 finished), 7 recipes, 10 raw batches. The
data is picked to exercise edge cases:

- `FP-BAR-PNT` has version 1 (inactive) and version 2 (active) — covers `one_active_per_item` and
  "give me the active recipe" queries.
- Two-level BOM (`finished → semi_finished → raw_material`) — covers recursive explosion.
- Expired batches (`CP-2602-A`, `HZ-2604-A`) and a deliberately tiny one (`HZ-2609-A`, 5 kg) — cover
  FEFO picking and insufficient stock.

Known gap: `SF-NOUGAT` has a `published`, active recipe with no rows in `recipe_components`. The
schema permits an empty recipe, so exploding `FP-BAR-PNT` stops at nougat.

## Migrations

`migrations/NNNNNN_<name>.{up,down}.sql`, applied by `migrate/migrate` before the API starts.
