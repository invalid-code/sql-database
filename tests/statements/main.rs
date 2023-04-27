use cucumber::{given, when, World};
use sql_database::repl::{Database, PersistantDatabase, StatementType, Table};
mod create;
mod insert;

#[derive(Debug, World, Default)]
struct StatementWorld {
    per_db: PersistantDatabase,
    command: String,
}

#[given("a persistant database")]
fn create_persistant_database(world: &mut StatementWorld) {
    world.per_db = PersistantDatabase::create_persistant_database();
}

#[given("a database")]
fn create_database(world: &mut StatementWorld) {
    let db = Database::create_database("first_db".to_owned());
    world.per_db.push_db(&db);
}

#[given("a table")]
fn create_table(world: &mut StatementWorld) {
    let table = Table::create_table("first_table".to_owned());
    match world.per_db.push_table("first_db", table) {
        Ok(_) => (),
        Err(err) => panic!("{:?}", err),
    }
}

#[when("I execute the command")]
fn execute_insert_one(world: &mut StatementWorld) {
    StatementType::execute_statement(&world.command, &mut world.per_db);
}

fn main() {
    futures::executor::block_on(StatementWorld::run("tests/features/statements.feature"));
}
