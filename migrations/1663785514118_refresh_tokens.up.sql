CREATE TABLE IF NOT EXISTS `refresh_tokens` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `token` text NOT NULL,
    `status` varchar(50) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `expiry_at` timestamp NOT NULL,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);