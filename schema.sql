create table species(name text, height_m numeric, length_m numeric, weight_kg numeric, popularity integer);

create table dinosaur(id integer, player text, positioned_id unique, latitude numeric, longitude numeric, catch_time integer, name text, power integer, health integer);

create table user(name text unique, salt blob, hash blob, updated integer);

create table config(token_key text);
