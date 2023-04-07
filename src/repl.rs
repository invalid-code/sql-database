use std::io::{stdin, stdout, Write};

enum MetaCommandResult {
    Success,
    Unknown,
}

enum PrepareResult {
    Success,
    Unknown,
    SyntaxError,
}

enum StatementType {
    Insert,
    Select,
}

struct Statement {
    stype: Option<StatementType>,
    row: Option<Row>,
}

struct Row {
    id: u32,
    username: String,
    email: String,
}

struct Table {
    num_rows: i32,
    pages: Option<Vec<Row>>,
}

fn execute_meta_command(cmd: &String) -> Option<MetaCommandResult> {
    if cmd == ".exit\n" {
        println!("Goodbye!");
        std::process::exit(0);
    } else {
        Some(MetaCommandResult::Unknown)
    }
}

fn prepare_statement(cmd: &String, statement: &mut Statement) -> PrepareResult {
    if &cmd[..6] == "insert" {
        statement.stype = Some(StatementType::Insert);
        let args: Vec<&str> = cmd[6..].split(" ").collect();
        if args.len() < 3 {
            return PrepareResult::SyntaxError;
        }
        return PrepareResult::Success;
    }

    if cmd == "select\n" {
        statement.stype = Some(StatementType::Select);
        return PrepareResult::Success;
    }
    PrepareResult::Unknown
}

fn execute_statement(statement: &Statement, table: Table) {
    match statement.stype {
        Some(StatementType::Insert) => println!("Inserting some data"),
        Some(StatementType::Select) => println!("Selecting some data"),
        _ => (),
    }
}

fn serialize_row(source: &Row, destination: Option<Row>) -> Option<Row> {
    match destination {
        Some(_) => None,
        None => Some(Row {
            id: source.id.to_owned().clone(),
            username: source.username.to_owned().clone(),
            email: source.email.to_owned().clone(),
        }),
    }
}

fn deserialize_row(source: Option<Row>, destination: Row) {
    match source {
        Some(_) => (),
        None => drop(destination),
    }
}

fn read_input(buf: &mut String) {
    print!("cmd: ");
    stdout().flush().unwrap();
    stdin()
        .read_line(buf)
        .expect("error when reading from stdin");
}

pub fn cli() {
    let table = Table {
        num_rows: 0,
        pages: None,
    };
    loop {
        let mut command = String::new();

        read_input(&mut command);

        if command.chars().nth(0).unwrap() == '.' {
            if let Some(mcr) = execute_meta_command(&command) {
                match mcr {
                    MetaCommandResult::Success => {
                        continue;
                    }
                    MetaCommandResult::Unknown => {
                        println!("Pls provide a valid command");
                        continue;
                    }
                }
            }
        }

        let mut statement = Statement {
            stype: None,
            row: None,
        };
        match prepare_statement(&command, &mut statement) {
            PrepareResult::Success => execute_statement(&statement, table),
            PrepareResult::Unknown => {
                println!("Unknown statement found");
                continue;
            }
            PrepareResult::SyntaxError => {
                println!("syntax error");
                continue;
            }
        }
    }
}
