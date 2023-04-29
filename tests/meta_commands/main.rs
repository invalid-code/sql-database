use cucumber::{given, then, when, World};
use sql_database::repl::{MetaCommandType, PersistantDatabase};

#[derive(Debug, World, Default)]
struct MetaCmdWorld {
    command: String,
    per_db: Option<PersistantDatabase>,
}

#[given("a open command")]
fn create_open_command(world: &mut MetaCmdWorld) {
    world.command = ".open database.db".to_string();
}

#[when("I execute the meta command")]
fn execute_meta_command(world: &mut MetaCmdWorld) {
    MetaCommandType::execute_meta_command(&world.command, &mut world.per_db, &mut None);
}

#[then("there should be a persistent database")]
fn check_per_db(world: &mut MetaCmdWorld) {
    assert!(world.per_db.is_some());
}

fn main() {
    futures::executor::block_on(MetaCmdWorld::run("tests/features/meta_command.feature"));
}
