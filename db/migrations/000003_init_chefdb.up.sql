START TRANSACTION;

-- DO $$ BEGIN IF NOT EXISTS(

--         SELECT 1

--         FROM pg_namespace

--         WHERE

--             nspname = 'chef'

--     ) THEN CREATE SCHEMA "chef";

-- END IF;

-- END $$;

CREATE SCHEMA IF NOT EXISTS "chef";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    chef.chef_orders (
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
            CONSTRAINT pk_chef_orders PRIMARY KEY (id)
    );

CREATE UNIQUE INDEX ix_chef_orders_id ON chef.chef_orders (id);

COMMIT;
