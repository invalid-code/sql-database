use cucumber::{given, then, when, World};
use sql_database::repl;

#[derive(Debug, World, Default)]
struct MetaCmdWorld {
    command: String,
}

fn main() {
    futures::executor::block_on(MetaCmdWorld::run("tests/features/meta_command.feature"));
}
