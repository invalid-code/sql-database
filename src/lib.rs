mod command_processor;

pub mod repl {
    pub use crate::command_processor::*;
    use std::io::{stdin, stdout, Write};

    fn read_input(buf: &mut String) {
        print!("cmd: ");
        stdout().flush().unwrap();
        stdin()
            .read_line(buf)
            .expect("error when reading from stdin");
    }

    pub fn cli() {
        let mut table = Table::create_table();

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

            match Statement::prepare_statement(&command, &mut statement) {
                PrepareResult::Success => {
                    if let Some(exec_res) = Statement::execute_statement(&statement, &mut table) {
                        match exec_res {
                            ExecuteResult::Success => println!("Executed"),
                            ExecuteResult::TableFull => println!("table full"),
                        }
                    }
                }
                PrepareResult::Unknown => {
                    println!("Unknown statement found");
                    continue;
                }
                PrepareResult::SyntaxError => {
                    println!("syntax error");
                    continue;
                }
                PrepareResult::InputTooLong => {
                    println!("input too long");
                    continue;
                }
            }
        }
    }
}
