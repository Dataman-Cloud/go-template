CREATE DATABASE IF NOT EXISTS notification;

CREATE TABLE IF NOT EXISTS message (
  id     bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  message_type     VARCHAR(64),
  resource_id     bigint(20),
  resource_type     VARCHAR(64),
  sink_name  VARCHAR(64),
  dump_time timestamp,
  created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  primary key (id)
);
