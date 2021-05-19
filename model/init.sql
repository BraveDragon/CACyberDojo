CREATE TABLE `users` (
	`ID` VARCHAR(256) NOT NULL DEFAULT '0' COLLATE 'utf8mb4_0900_ai_ci',
	`name` CHAR(50) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`token` TEXT NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`mailAddress` TEXT NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`password` TEXT NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`score` INT(10) UNSIGNED NOT NULL DEFAULT '0',
	PRIMARY KEY (`ID`) USING BTREE,
	UNIQUE INDEX `name` (`name`) USING BTREE
)
COMMENT='ユーザー情報を管理'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;

CREATE TABLE `characters` (
	`id` BIGINT(20) UNSIGNED NOT NULL,
	`name` CHAR(255) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`strength` INT(10) UNSIGNED NOT NULL,
	`rarity` INT(10) NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `name` (`name`) USING BTREE
)
COMMENT='キャラクターを管理するテーブル。'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;

CREATE TABLE `owncharacters` (
	`userId` VARCHAR(256) NOT NULL DEFAULT '0' COLLATE 'utf8mb4_0900_ai_ci',
	`characterId` INT(10) NOT NULL
)
COMMENT='ユーザーが保持するキャラクターを管理'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;

CREATE TABLE `gachas` (
	`gachaId` INT UNSIGNED NOT NULL,
	`characterId` INT UNSIGNED NOT NULL,
	`dropRate` FLOAT NOT NULL
)
COMMENT='ガチャを管理。同一のgachaIdを持つ項目のdropRateの合計が1になっているかどうか確認すること。'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;