START TRANSACTION;

-- DO $$ BEGIN IF NOT EXISTS(

--         SELECT 1

--         FROM pg_namespace

--         WHERE

--             nspname = 'cook'

--     ) THEN CREATE SCHEMA "cook";

-- END IF;

-- END $$;

CREATE SCHEMA IF NOT EXISTS "cook";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    cook.cook_orders (
        id uuid NOT NULL DEFAULT (uuid_generate_v4()),
        order_id uuid NOT NULL,
        item_type integer NOT NULL,
        item_name text NOT NULL,
        time_up timestamp
        with
            time zone NOT NULL,
            created timestamp
        with
            time zone NOT NULL DEFAULT (now()),
            updated timestamp
        with
            time zone NULL,
            CONSTRAINT pk_cook_orders PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX ix_cook_orders_id ON cook.cook_orders (id);

COMMIT;
