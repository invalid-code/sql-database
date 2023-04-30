use super::StatementWorld;
use cucumber::{given, then};

#[given("a insert command")]
fn create_insert_cmd(world: &mut StatementWorld) {
    world.command = "insert first_db first_table 1 admin@web.net admin".to_string();
}

#[then("the table should have 1 row")]
fn check_insert_one(world: &mut StatementWorld) {
    assert!(world.per_db.dbs[0].tables[0].num_rows == 1);
}
