BEGIN;

INSERT INTO items (number, name, unit, type) VALUES
('RM-COCOA-P', 'Какао-порошок', 'kg', 'raw_material'),
('RM-COCOA-B', 'Какао-масло', 'kg', 'raw_material'),
('RM-SUGAR', 'Сахар', 'kg', 'raw_material'),
('RM-MILK-P', 'Сухое молоко', 'kg', 'raw_material'),
('RM-PEANUT', 'Арахис', 'kg', 'raw_material'),
('RM-HAZELNUT', 'Фундук', 'kg', 'raw_material'),
('SF-GLAZE', 'Шоколадная глазурь', 'kg', 'semi_finished'),
('SF-NOUGAT', 'Нуга', 'kg', 'semi_finished'),
('SF-CARAMEL', 'Карамельная масса', 'kg', 'semi_finished'),
('FP-BAR-PNT', 'Батончик с арахисом', 'pcs', 'finished'),
('FP-BAR-CRM', 'Батончик карамельный', 'pcs', 'finished'),
('FP-BOX-HZL', 'Конфеты с фундуком', 'pcs', 'finished');

INSERT INTO recipes (item_id, item_type, version, status, is_active)
SELECT
    i.id,
    i.type,
    v.version,
    'published',
    v.is_active
FROM items AS i
INNER JOIN (
    VALUES
    ('SF-GLAZE', 1, TRUE),
    ('SF-NOUGAT', 1, TRUE),
    ('SF-CARAMEL', 1, TRUE),
    ('FP-BAR-PNT', 1, FALSE), ('FP-BAR-PNT', 2, TRUE),
    ('FP-BAR-CRM', 1, TRUE),
    ('FP-BOX-HZL', 1, TRUE)
) AS v (number, version, is_active) ON i.number = v.number;

INSERT INTO recipe_components (recipe_id, item_id, component_item_id, quantity)
SELECT
    r.id,
    r.item_id,
    c.id,
    v.quantity
FROM (
    VALUES
    ('SF-GLAZE', 1, 'RM-COCOA-P', 0.400),
    ('SF-GLAZE', 1, 'RM-COCOA-B', 0.250),
    ('SF-GLAZE', 1, 'RM-SUGAR', 0.350),
    ('SF-CARAMEL', 1, 'RM-SUGAR', 0.700),
    ('SF-CARAMEL', 1, 'RM-MILK-P', 0.300),
    ('FP-BAR-PNT', 1, 'SF-GLAZE', 0.020),
    ('FP-BAR-PNT', 1, 'SF-NOUGAT', 0.030),
    ('FP-BAR-PNT', 1, 'RM-PEANUT', 0.008),
    ('FP-BAR-PNT', 2, 'SF-GLAZE', 0.020),
    ('FP-BAR-PNT', 2, 'SF-NOUGAT', 0.030),
    ('FP-BAR-PNT', 2, 'RM-PEANUT', 0.010),
    ('FP-BAR-CRM', 1, 'SF-GLAZE', 0.015),
    ('FP-BAR-CRM', 1, 'SF-CARAMEL', 0.035),
    ('FP-BAR-CRM', 1, 'RM-SUGAR', 0.005),
    ('FP-BOX-HZL', 1, 'SF-GLAZE', 0.050),
    ('FP-BOX-HZL', 1, 'SF-CARAMEL', 0.040),
    ('FP-BOX-HZL', 1, 'RM-HAZELNUT', 0.020)
) AS v (parent_number, version, component_number, quantity)
INNER JOIN items AS p ON v.parent_number = p.number
INNER JOIN items AS c ON v.component_number = c.number
INNER JOIN recipes AS r ON p.id = r.item_id AND v.version = r.version;

INSERT INTO stock_batches (
    item_id, item_type, supplier_batch_number,
    initial_quantity, received_at, expires_at
)
SELECT
    i.id,
    i.type,
    v.batch_number,
    v.quantity,
    v.received_at,
    v.expires_at
FROM (
    VALUES
    ('RM-COCOA-P', 'CP-2609-A', 500.000000, DATE '2026-09-01', DATE '2027-03-01'),
    ('RM-COCOA-P', 'CP-2608-A', 120.000000, DATE '2026-08-01', DATE '2026-12-01'),
    ('RM-COCOA-P', 'CP-2602-A', 80.000000, DATE '2026-02-01', DATE '2026-08-01'),
    ('RM-COCOA-B', 'CB-2609-A', 300.000000, DATE '2026-09-05', DATE '2027-06-01'),
    ('RM-SUGAR', 'SG-2607-A', 200.000000, DATE '2026-07-10', DATE '2027-01-10'),
    ('RM-SUGAR', 'SG-2609-A', 800.000000, DATE '2026-09-10', DATE '2027-09-10'),
    ('RM-MILK-P', 'MP-2608-A', 400.000000, DATE '2026-08-15', DATE '2027-02-15'),
    ('RM-PEANUT', 'PN-2609-A', 150.000000, DATE '2026-09-20', DATE '2027-03-20'),
    ('RM-HAZELNUT', 'HZ-2609-A', 5.000000, DATE '2026-09-22', DATE '2027-03-22'),
    ('RM-HAZELNUT', 'HZ-2604-A', 40.000000, DATE '2026-04-01', DATE '2026-09-01')
) AS v (number, batch_number, quantity, received_at, expires_at)
INNER JOIN items AS i ON v.number = i.number;

COMMIT;
