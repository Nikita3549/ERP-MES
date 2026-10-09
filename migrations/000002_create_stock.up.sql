BEGIN;

CREATE TABLE stock_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL,
    item_type item_type NOT NULL,
    supplier_batch_number TEXT NOT NULL,
    initial_quantity NUMERIC(18, 6) NOT NULL CHECK (initial_quantity > 0),
    reserved_quantity NUMERIC(18, 6) NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    remaining_quantity NUMERIC(18, 6)
    GENERATED ALWAYS AS (initial_quantity - reserved_quantity) STORED,
    received_at DATE NOT NULL,
    expires_at DATE NOT NULL,

    CONSTRAINT stock_batches_item_fk
    FOREIGN KEY (item_id, item_type) REFERENCES items (id, type),
    CONSTRAINT stock_batches_raw_only
    CHECK (item_type = 'raw_material'),
    CONSTRAINT stock_batches_reserved_le_initial
    CHECK (reserved_quantity <= initial_quantity),
    CONSTRAINT stock_batches_expiry_after_receipt
    CHECK (expires_at >= received_at),
    CONSTRAINT stock_batches_supplier_number_key
    UNIQUE (item_id, supplier_batch_number)
);

CREATE INDEX stock_batches_fifo_idx
ON stock_batches (item_id, expires_at)
WHERE reserved_quantity < initial_quantity;

COMMIT;
