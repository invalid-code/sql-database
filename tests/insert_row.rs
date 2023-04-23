use cucumber::{given, then, when, World};
use sql_database::repl::{ExecuteResult, Row, Statement, StatementType, Table};

#[derive(Debug, World, Default)]
struct RowWorld {
    statement: Statement,
    table: Table,
    execute_result: Option<ExecuteResult>,
}

#[given("A Table")]
fn create_table(world: &mut RowWorld) {
    world.table = Table::create_table();
}

#[given("A Execute Result")]
fn create_execute_result(world: &mut RowWorld) {
    world.execute_result = None;
}

#[given("A Statement")]
fn create_statement(world: &mut RowWorld) {
    let row = Some(Row {
        id: 1,
        username: String::from("admin"),
        email: String::from("admin@web.net"),
    });
    world.statement = Statement {
        stype: Some(StatementType::Insert),
        row,
    };
}

#[when("I execute A Statement")]
fn execute_statement(world: &mut RowWorld) {
    world.execute_result = Statement::execute_statement(&world.statement, &mut world.table);
}

#[then("A Statement should be executed")]
fn check_statement(world: &mut RowWorld) {
    if let Some(res) = &world.execute_result {
        assert!(res.to_owned() == ExecuteResult::Success);
    }
}

fn main() {
    futures::executor::block_on(RowWorld::run("tests/features/insert_row.feature"));
}
