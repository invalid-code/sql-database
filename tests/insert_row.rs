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
fn execute_command(world: &mut RowWorld) {
    match StatementType::parse_statement(&world.command) {
        Ok(statement) => {
            if let StatementType::Insert(id, email, username, dname, tname) = statement {
                match StatementType::execute_insert(id, email, username, dname, tname, &mut per_db)
                {
                    Ok(_) => (),
                    Err(err) => match err {
                        ExecuteErr::DatabaseDoesNotExist => {
                            panic!("database does not exist");
                        }
                        ExecuteErr::TableDoesNotExist => {
                            panic!("table does not exist");
                        }
                    },
                }
            }
        }
        Err(err) => match err {
            StatementErr::Unknown => panic!("unknown statement found"),
            StatementErr::SyntaxErr => panic!("invalid syntax found"),
        },
    }
}

#[then("the command should be executed")]
fn check_command(world: &mut RowWorld) {}

fn main() {
    futures::executor::block_on(RowWorld::run("tests/features/insert_row.feature"));
}
