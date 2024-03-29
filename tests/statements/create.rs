use super::StatementWorld;
use cucumber::{given, then};

#[given("a create database command")]
fn create_create_db_cmd(world: &mut StatementWorld) {
    world.command = "create db first_db".to_owned();
}

#[then("the persistant database should have 1 database")]
fn check_db_one(world: &mut StatementWorld) {
    assert!(world.per_db.num_dbs > 0);
}

#[given("a create table command")]
fn create_create_table_cmd(world: &mut StatementWorld) {
    world.command = "create table first_db first_table".to_owned();
}

#[then("the database should have 1 table")]
fn check_table_one(world: &mut StatementWorld) {
    assert!(world.per_db.dbs[0].num_tables == 1);
}
