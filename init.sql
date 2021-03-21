CREATE TABLE `users` (
	`id` TEXT NOT NULL,
	`mailAddress` TEXT NOT NULL,
	`passWord` TEXT NOT NULL,
	`name` TEXT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`privateKey` BLOB NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE INDEX `privateKey` (`privateKey`),
	UNIQUE INDEX `name` (`name`),
	UNIQUE INDEX `mailAddress` (`mailAddress`),
	UNIQUE INDEX `passWord` (`passWord`)
)
COMMENT='ユーザー情報を管理'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;