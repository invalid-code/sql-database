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
        // let db = Database::create_database();
        let mut per_db = PersistantDatabase::create_persistant_database();

        loop {
            let mut command = String::new();

            read_input(&mut command);

            match StatementType::parse_statement(&command) {
                StatementType::Insert(_, _, _) => todo!(),
                StatementType::Select(_, _, _) => todo!(),
                StatementType::Create(dstruct, dstructn, dname) => match dname {
                    Some(dname) => StatementType::execute_create(
                        dstruct,
                        dstructn,
                        per_db.get_db(&dname),
                        &mut per_db,
                    ),
                    None => StatementType::execute_create(dstruct, dstructn, None, &mut per_db),
                },
                StatementType::Err(err) => match err {
                    StatementErr::Unknown => println!("unknown statement found!"),
                    StatementErr::SyntaxErr => println!("invalid statement found!"),
                },
            }
        }
    }
}
