CREATE SCHEMA dealls_bumble;

-- users
drop type if exists sex;
create type sex AS ENUM('female', 'male');

create table if not exists dealls_bumble.users
(
    id    SERIAL primary key,
    uid CHAR(16) NOT NULL UNIQUE,
    name varchar(50) NOT NULL,
    email varchar(50) UNIQUE NOT NULL,
    username varchar(20) UNIQUE NOT NULL,
    hashed_password BYTEA NOT NULL,
    sex sex NOT NULL,
    -- age  INT NOT NULL check (age >= 18),
    birthdate TIMESTAMP NOT NULL,
    verified bool NOT NULL DEFAULT false,
    max_swipes INT DEFAULT 10,
    created_at TIMESTAMP DEFAULT current_timestamp,
    is_deleted bool NOT NULL DEFAULT false
);

create index if not exists users_uid on dealls_bumble.users using hash (uid);
create index if not exists users_name on dealls_bumble.users using hash (name);
create index if not exists users_email on dealls_bumble.users using hash (email);
create index if not exists users_username on dealls_bumble.users using hash (username);
create index if not exists users_sex on dealls_bumble.users using hash (sex);

-- user_images
create table if not exists dealls_bumble.user_images
(
    id    SERIAL primary key,
    user_id BIGINT NOT NULL,
    title varchar(50),
    url VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    is_deleted bool NOT NULL DEFAULT false
);

alter table dealls_bumble.user_images
	add constraint fk_user_id foreign key (user_id) references users(id) on delete cascade;

-- user_matches
create table if not exists dealls_bumble.users_matches (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    match_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    is_deleted bool NOT NULL DEFAULT false
);

alter table dealls_bumble.users_matches
	add constraint fk_user_id foreign key (user_id) references users(id) on delete cascade;
alter table dealls_bumble.users_matches
	add constraint fk_match_id foreign key (match_id) references users(id) on delete cascade;

-- premium packages

create table if not exists dealls_bumble.premium_packages
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    perks_code VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    is_deleted bool NOT NULL DEFAULT false
);

-- perks
create table if not exists dealls_bumble.perks
(
    id SERIAL PRIMARY KEY,
    perks_code VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT current_timestamp,
    is_deleted bool NOT NULL DEFAULT false
);
