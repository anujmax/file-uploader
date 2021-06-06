CREATE SCHEMA IF NOT EXISTS file_uploader;

use file_uploader;

create table file_metadata
(
    file_identifier varchar(36)   NOT NULL,
    file_name       VARCHAR(1000) NOT NULL,
    file_size       INT           NOT NULL,
    file_type       varchar(36)   NOT NULL,
    created_date    DATE,
    PRIMARY KEY (file_identifier)
);