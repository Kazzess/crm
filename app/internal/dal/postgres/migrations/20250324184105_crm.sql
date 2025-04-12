-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE category (
	id SERIAL,
	-- unique identifier is raised automatically.
	name TEXT NOT NULL,
	-- name of category, required field.
	is_enabled BOOLEAN NOT NULL,
	-- category availability, true - available, false unavailable.
	created_at TIMESTAMPTZ NOT NULL,
	-- date and time when category was created.
	updated_at TIMESTAMPTZ NOT NULL,
	-- date and time when category was updated.
	CONSTRAINT category_id_pk PRIMARY KEY (id) -- rule that made the id field the primary key.
);
CREATE TABLE goods (
	id UUID,
	-- unique identifier is generated automatically.
	name TEXT NOT NULL,
	--  name of goods, required field.
	category_id INT NOT NULL,
	-- its id field from table category.
	is_enabled BOOLEAN NOT NULL,
	-- amount of goods is bigger than 0.
	goods_type TEXT NOT NULL,
	-- what type is it, socks or T Shirt.
	product_type TEXT NOT NULL,
	-- breakable/unbreakable goods.
	goods_count INT NOT NULL,
	-- amount of goods.
	defective_count INT NOT NULL,
	-- amount of defective goods minimum 0 <= goods count.
	code TEXT NOT NULL,
	-- TODO DOCUMENTARY
	created_at TIMESTAMPTZ NOT NULL,
	--  date and time when category was created. 
	updated_at TIMESTAMPTZ NOT NULL,
	--  date and time when category was updated.
	CONSTRAINT goods_id_pk PRIMARY KEY (id),
	-- rule that made the id field the primary  key.
	CONSTRAINT category_id_fk FOREIGN KEY (category_id) REFERENCES category (id),
	CONSTRAINT name_goods_type_unq UNIQUE (name, goods_type)
);
CREATE TABLE collection (
	id SERIAL,
	-- unique identifier is raised automatically.
	-- category_id INT NOT NULL
	name TEXT NOT NULL,
	-- name of category, required field.
	is_enabled BOOLEAN NOT NULL,
	-- true - available, false - unavailable.
	created_at TIMESTAMPTZ NOT NULL,
	-- date and time when collection was created.
	updated_at TIMESTAMPTZ NOT NULL,
	-- date and time when collection was updated. 
	CONSTRAINT collection_id_pk PRIMARY KEY (id),
	--CONSTRAINT category_id_fk FOREIGN KEY (category_id) REFERENCES category(id)
	CONSTRAINT collection_name_unq UNIQUE (name)
);
CREATE TABLE collection_goods (
	collection_id INT NOT NULL,
	goods_id UUID NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	CONSTRAINT collection_id_fk FOREIGN KEY (collection_id) REFERENCES collection (id) ON DELETE CASCADE,
	CONSTRAINT goods_id_fk FOREIGN KEY (goods_id) REFERENCES goods (id),
	CONSTRAINT collection_goods_unq UNIQUE (collection_id, goods_id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd