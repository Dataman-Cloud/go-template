CREATE DATABASE IF NOT EXISTS notification;

CREATE TABLE IF NOT EXISTS message (

 id     VARCHAR(64),
 type     VARCHAR(64),
 resource_id     VARCHAR(64),
 resource_type     VARCHAR(64),

UNIQUE(id)
);
