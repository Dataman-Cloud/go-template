CREATE DATABASE IF NOT EXISTS notification;

CREATE TABLE IF NOT EXISTS message (
id     VARCHAR(64),
type     VARCHAR(64),
resource_id     VARCHAR(64),
resource_type     VARCHAR(64),
sink_name  VARCHAR(64),
dump_time timestamp,
primary key (id, sink_name)
);
