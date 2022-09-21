CREATE TABLE IF NOT EXISTS `user_login_mobile_otps` (
    `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `mobile` int NOT NULL,
    `otp` int NOT NULL,
    `verification_id` varchar(255) NOT NULL,
    `status` varchar(50) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `expiry_at` timestamp NOT NULL,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);