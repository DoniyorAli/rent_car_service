BEGIN;

CREATE TABLE IF NOT EXISTS "brand" (
	"id" CHAR(36) NOT NULL PRIMARY KEY,
	"name" VARCHAR(255) UNIQUE NOT NULL,
    "country" TEXT NOT NULL,
    "manufacturer" TEXT NOT NULL,
	"about_brand" TEXT NOT NULL,
	"created_at" TIMESTAMP DEFAULT now(),
	"updated_at" TIMESTAMP,
	"deleted_at" TIMESTAMP
);
CREATE TABLE IF NOT EXISTS "car" (
	"id" CHAR(36) PRIMARY KEY,
	"model" VARCHAR(255) NOT NULL,
	"color" VARCHAR(255) NOT NULL,
    "car_type" VARCHAR(255) NOT NULL,
    "mileage" VARCHAR(255) NOT NULL,
    "year" VARCHAR(255) NOT NULL,
    "price" VARCHAR(255) NOT NULL,
    "brand_id" CHAR(36) REFERENCES "brand" (id),
	"created_at" TIMESTAMP DEFAULT now(),
	"updated_at" TIMESTAMP,
	"deleted_at" TIMESTAMP
);

COMMIT;