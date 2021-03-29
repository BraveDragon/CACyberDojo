CREATE TABLE `users` (
	`ID` CHAR(64) NOT NULL DEFAULT '' COLLATE 'utf8mb4_0900_ai_ci',
	`name` CHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`privateKey` BLOB NOT NULL,
	`mailAddress` CHAR(255) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`password` CHAR(255) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	PRIMARY KEY (`ID`) USING BTREE,
	UNIQUE INDEX `password` (`password`) USING BTREE,
	UNIQUE INDEX `mailAddress` (`mailAddress`) USING BTREE,
	UNIQUE INDEX `name` (`name`) USING BTREE
)
COMMENT='ユーザー情報を管理'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;

CREATE TABLE `characters` (
	`id` INT(10) UNSIGNED NOT NULL,
	`name` CHAR(255) NOT NULL DEFAULT '' COLLATE 'utf8mb4_0900_ai_ci',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `name` (`name`) USING BTREE
)
COMMENT='キャラクターを管理するテーブル。'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;

CREATE TABLE `owncharacters` (
	`userId` CHAR(64) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`characterId` INT(10) NOT NULL
)
COMMENT='ユーザーが保持するキャラクターを管理'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;

CREATE TABLE `gachas` (
	`gachaId` INT(10) UNSIGNED NOT NULL,
	`characterId` INT(10) UNSIGNED NOT NULL,
	`dropRate` FLOAT NOT NULL,
	UNIQUE INDEX `characterId` (`characterId`) USING BTREE
)
COMMENT='ガチャを管理。同一のgachaIdを持つ項目のdropRateの合計が1になっているかどうか確認すること。'
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB;