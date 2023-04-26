use cucumber::{given, World};
use sql_database::repl::{Database, PersistantDatabase, Table};
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

fn main() {
    futures::executor::block_on(StatementWorld::run("tests/features/statements.feature"));
}
