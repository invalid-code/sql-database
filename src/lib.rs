mod command_processor;
mod db_arch;

pub mod repl {
    pub use crate::command_processor::*;
    pub use crate::db_arch::*;
    use std::io::{stdin, stdout, Write};

    fn read_input(buf: &mut String) {
        print!("cmd: ");
        stdout().flush().unwrap();
        stdin()
            .read_line(buf)
            .expect("error when reading from stdin");
    }

    pub fn cli() {
        let path = String::from("database.db");
        let file = open_file(&path);

        loop {
            let mut command = String::new();
            let mut db: Option<Database> = None;

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

            match Statement::prepare_statement(&command, &mut statement, db.as_ref()) {
                PrepareResult::Success => {
                    if let Some(exec_res) =
                        Statement::execute_statement(&statement, &mut db.unwrap())
                    {
                        match exec_res {
                            ExecuteResult::Success => println!("Executed"),
                        }
                    }
                }
                PrepareResult::Unknown => {
                    println!("Unknown Statement found: {:?} statement found", statement);
                    continue;
                }
                PrepareResult::SyntaxError => {
                    println!("Invalid Syntax found: {:?}", statement);
                    continue;
                }
                PrepareResult::NoExistingTable => {
                    println!("Table does not exist");
                    continue;
                }
                PrepareResult::NoExistingDatabase => {
                    println!("Database does not exist");
                    continue;
                }
            }
            write_file(&path, file);
        }
    }
}