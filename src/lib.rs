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
        // let path = String::from("database.db");

        let mut per_db = PersistantDatabase::create_persistant_database();

        loop {
            let mut command = String::new();

            read_input(&mut command);

            MetaCommandType::execute_meta_command(&command);

            StatementType::execute_statement(&command, &mut per_db);
        }
    }
}
