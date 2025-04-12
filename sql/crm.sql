INSERT INTO category(name, is_enabled, created_at, updated_at)
VALUES 
    ('Одежда', FALSE, NOW(), NOW()),
    ('Бумажные изделия', FALSE, NOW(), NOW()),
	('Аксессуары', FALSE, NOW(), NOW());

INSERT INTO goods(id, name, category_id, is_enabled, goods_type, product_type, goods_count, defective_count, code, created_at, updated_at)
VALUES
	(gen_random_uuid(), 'Ступор мозговины', 31, FALSE, 'Футболка', 'не хрупоке', 0, 0, 'TSHIRT001', NOW(), NOW());	

DELETE FROM collection WHERE name = 'Футболка'

INSERT INTO collection (name, is_enabled, created_at, updated_at)
VALUES ('Футболка', FALSE, NOW(), NOW());

INSERT INTO collection_goods (collection_id, goods_id, created_at)
SELECT 2, id, NOW()
FROM goods WHERE name = 'Ступор мозговины';