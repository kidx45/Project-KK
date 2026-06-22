CREATE TABLE `users` (
  `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(255) UNIQUE NOT NULL,
  `hashed_password` VARCHAR(255) NOT NULL,
  `full_name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `phone_number` VARCHAR(255) UNIQUE NOT NULL,
  `is_email_verified` BOOLEAN NOT NULL DEFAULT false,
  `is_phone_number_verified` BOOLEAN NOT NULL DEFAULT false,
  `password_changed_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX (`full_name`),
  INDEX (`email`)
);

CREATE TABLE `accounts` (
  `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(255) NOT NULL,
  `balance` BIGINT NOT NULL,
  `currency` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX (`username`),
  CONSTRAINT `fk_accounts_username` FOREIGN KEY (`username`) REFERENCES `users` (`username`)
);

CREATE TABLE `entries` (
  `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `account_id` BIGINT UNSIGNED NOT NULL,
  `amount` BIGINT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX (`account_id`),
  CONSTRAINT `fk_entries_account_id` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`)
);

CREATE TABLE `transfers` (
  `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `from_account_id` BIGINT UNSIGNED NOT NULL,
  `to_account_id` BIGINT UNSIGNED NOT NULL,
  `amount` BIGINT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX (`from_account_id`),
  INDEX (`to_account_id`),
  INDEX (`from_account_id`, `to_account_id`),
  CONSTRAINT `fk_transfers_from_account` FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`),
  CONSTRAINT `fk_transfers_to_account` FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`)
);