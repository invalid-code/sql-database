use super::StatementWorld;
use cucumber::{given, then, when};
use sql_database::repl::StatementType;

#[given("a insert command")]
fn create_insert_cmd(world: &mut StatementWorld) {
    world.command = "insert 1 admin@web.net admin first_db first_table".to_owned();
}

#[when("I execute the command")]
fn execute_insert_one(world: &mut StatementWorld) {
    StatementType::execute_statement(&world.command, &mut world.per_db);
}

#[then("the table should have 1 row")]
fn check_insert_one(world: &mut StatementWorld) {
    assert!(world.per_db.dbs[0].tables[0].rows.len() == 1);
}
