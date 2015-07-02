
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE tasks (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '',
  `title` VARCHAR(255) NOT NULL COMMENT '',
  `description` VARCHAR(255) NOT NULL COMMENT '',
  `priority` INT NOT NULL COMMENT '',
  `created_at` BIGINT NOT NULL COMMENT '',
  `updated_at` BIGINT NOT NULL COMMENT '',
  `is_deleted` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '',
  `is_completed` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '',
  PRIMARY KEY (`id`)  COMMENT '',
  UNIQUE INDEX `id_UNIQUE` (`id` ASC)  COMMENT '');


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tasks
