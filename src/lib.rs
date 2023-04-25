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
        let mut per_db = PersistantDatabase::create_persistant_database();

        loop {
            let mut command = String::new();

            read_input(&mut command);

            match StatementType::parse_statement(&command) {
                Ok(statement) => match statement {
                    StatementType::Insert(id, email, username, dname, tname) => {
                        StatementType::execute_insert(
                            id,
                            email,
                            username,
                            dname,
                            tname,
                            &mut per_db,
                        );
                    }
                    StatementType::Select(dname, tname) => {
                        StatementType::execute_select(dname, tname, &per_db);
                    }
                    StatementType::Create(dstruct, dstructn, dname) => match dname {
                        Some(dname) => {
                            StatementType::execute_create(
                                dstruct,
                                dstructn,
                                per_db.get_db(&dname),
                                &mut per_db,
                            );
                        }
                        None => {
                            StatementType::execute_create(dstruct, dstructn, None, &mut per_db);
                        }
                    },
                },
                Err(err) => match err {
                    StatementErr::Unknown => println!("unknown statement found"),
                    StatementErr::SyntaxErr => println!("invalid statement found"),
                },
                // Statement::Result(val) => match val {
                //     StatementResult::Err(err) => match err {
                //         StatementErr::Unknown => println!("unknown statement found"),
                //         StatementErr::SyntaxErr => println!("invalid statement found"),
                //     },
                //     StatementResult::Success => (),
                // },
                // Statement::Statement(statement) => match statement {
                //     StatementType::Insert(id, email, username, dname, tname) => {
                //         StatementType::execute_insert(
                //             id,
                //             email,
                //             username,
                //             dname,
                //             tname,
                //             &mut per_db,
                //         );
                //     }
                //     StatementType::Select(dname, tname) => {
                //         StatementType::execute_select(dname, tname, &per_db);
                //     }
                //     StatementType::Create(dstruct, dstructn, dname) => match dname {
                //         Some(dname) => {
                //             StatementType::execute_create(
                //                 dstruct,
                //                 dstructn,
                //                 per_db.get_db(&dname),
                //                 &mut per_db,
                //             );
                //         }
                //         None => {
                //             StatementType::execute_create(dstruct, dstructn, None, &mut per_db);
                //         }
                //     },
                // },
            }
        }
    }
}
