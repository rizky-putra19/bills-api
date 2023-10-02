-- -------------------------------------------------------------
-- TablePlus 5.3.6(495)
--
-- https://tableplus.com/
--
-- Database: digital-goods-stg
-- Generation Time: 2023-05-15 12:14:05.4350
-- -------------------------------------------------------------


-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS categories_id_seq;

-- Table Definition
CREATE TABLE "public"."categories" (
    "id" int4 NOT NULL DEFAULT nextval('categories_id_seq'::regclass),
    "name" varchar(50),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS customers_id_seq;

-- Table Definition
CREATE TABLE "public"."customers" (
    "id" int4 NOT NULL DEFAULT nextval('customers_id_seq'::regclass),
    "name" varchar(50) DEFAULT ''::character varying,
    "email" varchar(50) DEFAULT ''::character varying,
    "phone_number" varchar(15) DEFAULT ''::character varying,
    "password" text DEFAULT ''::text,
    "is_verified" bool DEFAULT false,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS order_items_id_seq;

-- Table Definition
CREATE TABLE "public"."order_items" (
    "id" int4 NOT NULL DEFAULT nextval('order_items_id_seq'::regclass),
    "order_id" int4,
    "product_id" int4,
    "price" numeric(10,2),
    "cogs" numeric(10,2),
    "fulfillment_status" varchar(10),
    "serial_code" varchar(100),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    "fulfillment_reference_number" varchar(50),
    "client_number" varchar(50),
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS orders_id_seq;

-- Table Definition
CREATE TABLE "public"."orders" (
    "id" int4 NOT NULL DEFAULT nextval('orders_id_seq'::regclass),
    "customer_id" int4,
    "order_date" timestamp,
    "total_price" numeric(10,2),
    "status" varchar(10),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS payments_id_seq;

-- Table Definition
CREATE TABLE "public"."payments" (
    "id" int4 NOT NULL DEFAULT nextval('payments_id_seq'::regclass),
    "order_id" int4,
    "payment_date" timestamp,
    "payment_method" varchar(20),
    "payment_amount" numeric(10,2),
    "payment_expiry_time" timestamp,
    "admin_fee" numeric(10,2),
    "status" varchar(10),
    "reference_number" varchar(100),
    "account_number" varchar(100),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS products_id_seq;

-- Table Definition
CREATE TABLE "public"."products" (
    "id" int4 NOT NULL DEFAULT nextval('products_id_seq'::regclass),
    "product_name" varchar(50),
    "category_id" int4,
    "description" text,
    "partner_fee" numeric(10,2),
    "cogs" numeric(10,2),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    PRIMARY KEY ("id")
);

