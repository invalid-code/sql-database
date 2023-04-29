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
        let mut per_db: Option<PersistantDatabase> = None;
        let mut per_db_name: Option<String> = None;

        loop {
            let mut command = String::new();

            read_input(&mut command);

            match MetaCommandType::execute_meta_command(&command, &mut per_db, &mut per_db_name) {
                Ok(_) => {
                    continue;
                }
                Err(err) => match err {
                    MetaCommandErr::Unknown => {
                        println!("unknown meta command found!");
                        continue;
                    }
                    MetaCommandErr::NotMetaCommand => (),
                },
            }

            match StatementType::execute_statement(&command, per_db.as_mut()) {
                Ok(_) => (),
                Err(err) => match err {
                    StatementErr::Prepare(prepare_err) => match prepare_err {
                        PrepareErr::Unknown => {
                            println!("unknown statement found!");
                            continue;
                        }
                        PrepareErr::SyntaxErr => {
                            println!("invalid statement found!");
                            continue;
                        }
                    },
                    StatementErr::Execute(execute_err) => match execute_err {
                        ExecuteErr::DatabaseDoesNotExist => {
                            println!("database does not exist!");
                            continue;
                        }
                        ExecuteErr::TableDoesNotExist => {
                            println!("table does not exist!");
                            continue;
                        }
                        ExecuteErr::NoOpenDatabase => {
                            println!("no open database");
                            continue;
                        }
                    },
                },
            }
            if let Some(name) = &per_db_name {
                if let Some(db) = &per_db {
                    PersistantDatabase::save_db(name, db);
                }
            }
        }
    }
}
