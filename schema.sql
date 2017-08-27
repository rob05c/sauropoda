create table species(name text, height_m numeric, length_m numeric, weight_kg numeric, popularity integer);

create table player(name text);

create table player_dinosaur(player_name text, dinosaur_id integer);

create table dinosaur(id integer, name text, power integer, health integer);

create table user(name text unique, salt blob, hash blob, updated integer);

create table config(token_key text);
