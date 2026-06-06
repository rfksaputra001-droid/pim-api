-- ============================================================
-- PIM Database Schema
-- PostgreSQL
--
-- Jalankan file ini untuk setup database dari awal:
--   psql -U <user> -d <dbname> -f database/schema.sql
-- ============================================================

-- Enum types
CREATE TYPE "ContributionStatus" AS ENUM ('PENDING', 'VERIFIED', 'REJECTED');
CREATE TYPE "ExpenseCategory"    AS ENUM ('OPERASIONAL', 'KEGIATAN', 'PERALATAN', 'LAIN_LAIN');
CREATE TYPE "AdminRole"          AS ENUM ('ADMIN', 'SUPER_ADMIN');

-- Tabel Contribution
CREATE TABLE "Contribution" (
    id               TEXT                    NOT NULL,
    name             TEXT                    NOT NULL,
    phone            TEXT                    NOT NULL,
    amount           INTEGER                 NOT NULL,
    notes            TEXT,
    status           "ContributionStatus"    NOT NULL DEFAULT 'PENDING',
    "proofImageUrl"  TEXT                    NOT NULL,
    "rejectionReason" TEXT,
    "createdAt"      TIMESTAMP(3)            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "verifiedAt"     TIMESTAMP(3),
    "verifiedBy"     TEXT,

    CONSTRAINT "Contribution_pkey" PRIMARY KEY (id)
);

CREATE INDEX "Contribution_status_idx"    ON "Contribution" (status);
CREATE INDEX "Contribution_createdAt_idx" ON "Contribution" ("createdAt");
CREATE INDEX "Contribution_amount_idx"    ON "Contribution" (amount);

-- Tabel Expense
CREATE TABLE "Expense" (
    id               TEXT                NOT NULL,
    date             TIMESTAMP(3)        NOT NULL,
    description      TEXT                NOT NULL,
    category         "ExpenseCategory"   NOT NULL,
    amount           INTEGER             NOT NULL,
    "receiptImageUrl" TEXT               NOT NULL,
    "createdBy"      TEXT                NOT NULL,
    "createdAt"      TIMESTAMP(3)        NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Expense_pkey" PRIMARY KEY (id)
);

CREATE INDEX "Expense_date_idx" ON "Expense" (date);

-- Tabel Admin
CREATE TABLE "Admin" (
    id             TEXT         NOT NULL,
    username       TEXT         NOT NULL,
    "passwordHash" TEXT         NOT NULL,
    role           "AdminRole"  NOT NULL DEFAULT 'ADMIN',
    "createdAt"    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "lastLogin"    TIMESTAMP(3),

    CONSTRAINT "Admin_pkey"         PRIMARY KEY (id),
    CONSTRAINT "Admin_username_key" UNIQUE (username)
);

-- ============================================================
-- Seed: default superadmin
-- username : superadmin
-- password : @IRPK2206
-- ============================================================
INSERT INTO "Admin" (id, username, "passwordHash", role, "createdAt")
VALUES (
    'test-admin-001',
    'superadmin',
    '$2a$12$weWTMS0oyHNke7.gTKD1de8mcAgUpxJVJGvTJKWt6o2aHTslcQkze',
    'SUPER_ADMIN',
    NOW()
);
