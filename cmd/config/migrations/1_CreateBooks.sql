-- +migrate Up
CREATE TABLE `books` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `isbn` VARCHAR(255) NOT NULL UNIQUE,
    `title` VARCHAR(255) NOT NULL,
    `author` VARCHAR(255) NOT NULL,
    `publisher` VARCHAR(255) NOT NULL,
    `published` DATE NOT NULL,
    `genre` VARCHAR(255) NOT NULL,
    `language` VARCHAR(255) NOT NULL,
    `pages` INT NOT NULL,
    `availability` VARCHAR(30) NOT NULL,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    INDEX `idx_isbn` (`isbn`),
    INDEX `idx_title` (`title`),
    INDEX `idx_author` (`author`),
    INDEX `idx_published` (`published`),
    INDEX `idx_genre` (`genre`),
    INDEX `idx_language` (`language`),
    INDEX `idx_availability` (`availability`),
    INDEX `idx_deleted_at` (`deleted_at`),
    INDEX `idx_title_created_at_deleted_at` (`title`, `created_at`, `deleted_at`) USING BTREE,
    INDEX `idx_author_created_at_deleted_at` (`author`, `created_at`, `deleted_at`) USING BTREE,
    UNIQUE KEY `uk__title__author` (`title`, `author`)
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;
-- +migrate Down
DROP TABLE books;
