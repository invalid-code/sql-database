use cucumber::{given, then, when, World};
use sql_database::repl::{
    MetaCommandType, PersistantDatabase, PersistantDatabaseErr, StatementType,
};

#[derive(Debug, World, Default)]
struct PersistantWorld {
    open_cmd: String,
    create_command: String,
    per_db: Option<PersistantDatabase>,
    per_db_name: Option<String>,
}

#[given("a open command")]
fn create_per_db_name(world: &mut PersistantWorld) {
    world.open_cmd = ".open database.db".to_string()
}

#[given("a create command")]
fn create_create_meta_cmd(world: &mut PersistantWorld) {
    world.create_command = "create first_db".to_string()
}

#[given("a persistant database")]
fn create_per_db(world: &mut PersistantWorld) {
    world.per_db = None;
}

#[when("I execute all the commands")]
fn execute_commands(world: &mut PersistantWorld) {
    MetaCommandType::execute_meta_command(
        &world.open_cmd,
        &mut world.per_db,
        &mut world.per_db_name,
    );

    StatementType::execute_statement(&world.create_command, world.per_db.as_mut());
}

#[when("I save the database")]
fn execute_save_db(world: &mut PersistantWorld) {
    if let Some(name) = &world.per_db_name {
        if let Some(db) = &world.per_db {
            PersistantDatabase::save_db(name, db);
        }
    }
}

#[then("the database should have been saved")]
fn check_per_db(world: &mut PersistantWorld) {
    match PersistantDatabase::open_db(world.per_db_name.as_ref().unwrap()) {
        Ok(_) => assert!(world.per_db.as_ref().unwrap().num_dbs == 1),
        Err(_) => (),
    }
}

fn main() {
    futures::executor::block_on(PersistantWorld::run("tests/features/persistant.feature"))
}
