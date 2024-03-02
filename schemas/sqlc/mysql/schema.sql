CREATE TABLE `users` (
    `id` CHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `is_admin` TINYINT(1) NOT NULL DEFAULT 0,
    `created_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),

    PRIMARY KEY (`id`)
);



-- -- CreateTable
-- CREATE TABLE "users" (
--     "id" UUID NOT NULL,
--     "name" VARCHAR NOT NULL,
--     "email" VARCHAR NOT NULL,
--     "password" VARCHAR NOT NULL,
--     "is_admin" BOOLEAN NOT NULL DEFAULT false,
--     "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

--     PRIMARY KEY ("id")
-- );