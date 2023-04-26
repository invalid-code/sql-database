use crate::db_arch::*;

#[derive(Debug)]
pub enum StatementErr {
    Unknown,
    SyntaxErr,
}

#[derive(Debug)]
pub enum ExecuteErr {
    DatabaseDoesNotExist,
    TableDoesNotExist,
}

#[derive(Debug)]
pub enum StatementType {
    Insert(i32, String, String, String, String),
    Select(String, String),
    Create(String, String, Option<String>),
}

impl StatementType {
    pub fn execute_create(
        dstruct: String,
        dstructn: String,
        dname: Option<&str>,
        per_db: &mut PersistantDatabase,
    ) -> Result<(), ExecuteErr> {
        if dstruct == "db" {
            let db = Database::create_database(dstructn.clone());
            per_db.push_db(&db);
        }
        if dstruct == "table" {
            let table = Table::create_table(dstructn.clone());
            match dname {
                Some(dname) => match per_db.push_table(dname, table) {
                    Ok(_) => (),
                    Err(err) => match err {
                        ExecuteErr::DatabaseDoesNotExist => println!("database doesn't exist"),
                        ExecuteErr::TableDoesNotExist => (),
                    },
                },
                None => (),
            }
        }
        Ok(())
    }

    pub fn execute_insert(
        id: i32,
        email: String,
        username: String,
        dname: &str,
        tname: &str,
        per_db: &mut PersistantDatabase,
    ) -> Result<(), ExecuteErr> {
        let row = Row {
            id,
            email,
            username,
        };
        match per_db.push_row(dname, tname, row) {
            Ok(_) => (),
            Err(err) => return Err(err),
        }
        Ok(())
    }

    pub fn execute_select(
        dname: &str,
        tname: &str,
        per_db: &mut PersistantDatabase,
    ) -> Result<(), ExecuteErr> {
        match per_db.get_db(dname) {
            Some(mut db) => match db.get_table(tname) {
                Some(table) => {
                    for row in &table.rows {
                        println!("{} {} {}", row.id, row.email, row.username);
                    }
                }
                None => return Err(ExecuteErr::TableDoesNotExist),
            },
            None => return Err(ExecuteErr::DatabaseDoesNotExist),
        }
        Ok(())
    }
}

impl StatementType {
    pub fn parse_statement(cmd: &str) -> Result<Self, StatementErr> {
        if &cmd[..6] == "insert" {
            let args = cmd[7..].split(" ").collect::<Vec<&str>>();
            if args.len() < 5 {
                return Err(StatementErr::SyntaxErr);
            }
            return Ok(Self::Insert(
                std::str::FromStr::from_str(args[0]).unwrap(),
                args[1].to_owned(),
                args[2].to_owned(),
                args[3].to_owned(),
                args[4].to_owned(),
            ));
        }

        if &cmd[..6] == "select" {
            let args = cmd[7..].split(" ").collect::<Vec<&str>>();
            if args.len() < 2 {
                return Err(StatementErr::SyntaxErr);
            }
            return Ok(Self::Select(
                std::str::FromStr::from_str(args[0]).unwrap(),
                args[1].to_owned(),
            ));
        }

        if &cmd[..6] == "create" {
            let args = cmd[6..].split(" ").collect::<Vec<&str>>();
            if args.len() == 3 {
                return Ok(Self::Create(
                    args[0].to_owned(),
                    args[1].to_owned(),
                    Some(args[2].to_owned()),
                ));
            }
            if args.len() == 2 {
                return Ok(Self::Create(args[0].to_owned(), args[1].to_owned(), None));
            }
            return Err(StatementErr::SyntaxErr);
        }

        Err(StatementErr::Unknown)
    }

    pub fn execute_statement(command: &str, per_db: &mut PersistantDatabase) {
        match StatementType::parse_statement(command) {
            Ok(statement) => match statement {
                StatementType::Insert(id, email, username, dname, tname) => {
                    match StatementType::execute_insert(
                        id,
                        email,
                        username,
                        dname.as_str(),
                        tname.as_str(),
                        per_db,
                    ) {
                        Ok(_) => (),
                        Err(err) => match err {
                            ExecuteErr::DatabaseDoesNotExist => {
                                println!("database does not exist!");
                            }
                            ExecuteErr::TableDoesNotExist => {
                                println!("table does not exist!");
                            }
                        },
                    }
                }
                StatementType::Select(dname, tname) => {
                    match StatementType::execute_select(dname.as_str(), tname.as_str(), per_db) {
                        Ok(_) => (),
                        Err(err) => match err {
                            ExecuteErr::DatabaseDoesNotExist => {
                                println!("database does not exist!");
                            }
                            ExecuteErr::TableDoesNotExist => {
                                println!("table does not exist!");
                            }
                        },
                    }
                }
                StatementType::Create(dstruct, dstructn, dname) => match dname {
                    Some(dname) => {
                        match StatementType::execute_create(dstruct, dstructn, Some(&dname), per_db)
                        {
                            Ok(_) => (),
                            Err(err) => match err {
                                ExecuteErr::DatabaseDoesNotExist => {
                                    println!("database does not exist!");
                                }
                                ExecuteErr::TableDoesNotExist => {
                                    println!("table does not exist!");
                                }
                            },
                        }
                    }
                    None => match StatementType::execute_create(dstruct, dstructn, None, per_db) {
                        Ok(_) => (),
                        Err(err) => match err {
                            ExecuteErr::DatabaseDoesNotExist => {
                                println!("database does not exist!");
                            }
                            ExecuteErr::TableDoesNotExist => {
                                println!("table does not exist!");
                            }
                        },
                    },
                },
            },
            Err(err) => match err {
                StatementErr::Unknown => {
                    println!("unknown statement found!");
                }
                StatementErr::SyntaxErr => {
                    println!("invalid statement found!");
                }
            },
        }
    }
}

pub enum MetaCommandErr {
    Unknown,
}

pub enum MetaCommandType {
    Exit,
    Open(String),
}

impl MetaCommandType {
    pub fn parse_meta_command(cmd: &String) -> Result<Self, MetaCommandErr> {
        if &cmd[..1] == "." {
            if &cmd[2..5] == "exit" {
                return Ok(Self::Exit);
            }
        }
        Err(MetaCommandErr::Unknown)
    }

    pub fn execute_meta_command(cmd: &String) {
        match Self::parse_meta_command(cmd) {
            Ok(command) => match command {
                MetaCommandType::Exit => {
                    println!("Goodbye!");
                    std::process::exit(0);
                }
                MetaCommandType::Open(_) => todo!(),
            },
            Err(_err) => (),
            // match err {
            //     MetaCommandErr::Unknown => println!("unknown meta command found!"),
            // },
        }
    }
}
