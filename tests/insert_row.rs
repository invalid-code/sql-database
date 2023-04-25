use cucumber::{given, then, when, World};
use sql_database::repl::*;

#[derive(Debug, World, Default)]
struct RowWorld {
    command: String,
    table: Table,
}

#[given("a command")]
fn create_command(world: &mut RowWorld) {
    world.command = String::from("insert 1 admin@web.net admin first_db first_table");
}

#[given("a table")]
fn create_table(world: &mut RowWorld) {
    world.table = Table::create_table(String::from("first_table"));
}

#[when("I execute the command")]
fn execute_(world: &mut RowWorld) {}

fn main() {
    futures::executor::block_on(RowWorld::run("tests/features/insert_row.feature"));
}
