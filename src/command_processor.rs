use crate::db_arch::*;

pub enum StatementErr {
    Prepare(PrepareErr),
    Execute(ExecuteErr),
}

#[derive(Debug)]
pub enum PrepareErr {
    Unknown,
    SyntaxErr,
}

#[derive(Debug)]
pub enum ExecuteErr {
    DatabaseDoesNotExist,
    TableDoesNotExist,
    NoOpenDatabase,
}

#[derive(Debug)]
pub enum StatementType {
    Insert {
        db: String,
        table: String,
        id: i32,
        email: String,
        username: String,
    },
    Select {
        db: String,
        table: String,
    },
    Create {
        dtype: String,
        db: String,
        table: Option<String>,
    },
}

impl Default for StatementType {
    fn default() -> Self {
        Self::Create {
            dtype: "".to_string(),
            db: "".to_string(),
            table: None,
        }
    }
}

impl StatementType {
    pub fn execute_create(
        dstruct: String,
        dstructn: String,
        dname: Option<&str>,
        per_db: Option<&mut PersistantDatabase>,
    ) -> Result<(), ExecuteErr> {
        match per_db {
            Some(per_db) => {
                println!("{}", dstruct);
                if dstruct == "db" {
                    let db = Database::create_database(dstructn.clone());
                    per_db.push_db(&db);
                    // println!("{:?}", per_db);
                }
                if dstruct == "table" {
                    let table = Table::create_table(dstructn.clone());
                    if let Some(dname) = dname {
                        return per_db.push_table(dname, table);
                    }
                }
            }
            None => return Err(ExecuteErr::NoOpenDatabase),
        }
        Ok(())
    }

    pub fn execute_insert(
        id: i32,
        email: String,
        username: String,
        dname: &str,
        tname: &str,
        per_db: Option<&mut PersistantDatabase>,
    ) -> Result<(), ExecuteErr> {
        match per_db {
            Some(per_db) => {
                let row = Row {
                    id,
                    email,
                    username,
                };
                return per_db.push_row(dname, tname, row);
            }
            None => return Err(ExecuteErr::NoOpenDatabase),
        }
    }

    pub fn execute_select(
        dname: &str,
        tname: &str,
        per_db: Option<&mut PersistantDatabase>,
    ) -> Result<(), ExecuteErr> {
        match per_db {
            Some(per_db) => match per_db.get_db(dname) {
                Some(mut db) => match db.get_table(tname) {
                    Some(table) => {
                        for row in &table.rows {
                            println!("{} {} {}", row.id, row.email, row.username);
                        }
                    }
                    None => return Err(ExecuteErr::TableDoesNotExist),
                },
                None => return Err(ExecuteErr::DatabaseDoesNotExist),
            },
            None => return Err(ExecuteErr::NoOpenDatabase),
        }
        Ok(())
    }
    pub fn parse_statement(cmd: &str) -> Result<Self, PrepareErr> {
        let args = cmd.split(" ").collect::<Vec<&str>>();
        if args[0] == "insert" {
            if args[1..].len() < 5 {
                return Err(PrepareErr::SyntaxErr);
            }
            return Ok(Self::Insert {
                id: std::str::FromStr::from_str(args[0]).unwrap(),
                db: args[1].to_owned(),
                table: args[2].to_owned(),
                email: args[3].to_owned(),
                username: args[4].to_owned(),
            });
        }

        if args[0] == "select" {
            if args[1..].len() < 2 {
                return Err(PrepareErr::SyntaxErr);
            }
            return Ok(Self::Select {
                db: std::str::FromStr::from_str(args[0]).unwrap(),
                table: args[1].to_owned(),
            });
        }

        if args[0] == "create" {
            if args[1..].len() == 3 {
                return Ok(Self::Create {
                    dtype: args[0].to_owned(),
                    db: args[1].to_owned(),
                    table: Some(args[2].to_owned()),
                });
            }
            if args[1..].len() == 2 {
                return Ok(Self::Create {
                    dtype: args[1].to_owned(),
                    db: args[2].to_owned(),
                    table: None,
                });
            }
            return Err(PrepareErr::SyntaxErr);
        }

        Err(PrepareErr::Unknown)
    }

    pub fn execute_statement(
        command: &str,
        per_db: Option<&mut PersistantDatabase>,
    ) -> Result<(), StatementErr> {
        match StatementType::parse_statement(command) {
            Ok(statement) => match statement {
                StatementType::Insert {
                    db: dname,
                    table: tname,
                    id,
                    email,
                    username,
                } => {
                    match StatementType::execute_insert(
                        id,
                        email,
                        username,
                        dname.as_str(),
                        tname.as_str(),
                        per_db,
                    ) {
                        Ok(_) => (),
                        Err(err) => return Err(StatementErr::Execute(err)),
                    }
                }
                StatementType::Select {
                    db: dname,
                    table: tname,
                } => match StatementType::execute_select(dname.as_str(), tname.as_str(), per_db) {
                    Ok(_) => (),
                    Err(err) => return Err(StatementErr::Execute(err)),
                },
                StatementType::Create {
                    dtype: dstruct,
                    db: dstructn,
                    table: dname,
                } => match dname {
                    Some(dname) => {
                        match StatementType::execute_create(dstruct, dstructn, Some(&dname), per_db)
                        {
                            Ok(_) => (),
                            Err(err) => return Err(StatementErr::Execute(err)),
                        }
                    }
                    None => match StatementType::execute_create(dstruct, dstructn, None, per_db) {
                        Ok(_) => (),
                        Err(err) => return Err(StatementErr::Execute(err)),
                    },
                },
            },
            Err(err) => return Err(StatementErr::Prepare(err)),
        }
        Ok(())
    }
}

#[derive(Debug)]
pub enum MetaCommandErr {
    Unknown,
    NotMetaCommand,
}

#[derive(Debug, Default)]
pub enum MetaCommandType {
    #[default]
    Exit,
    Open(String),
}

impl MetaCommandType {
    pub fn parse_meta_command(cmd: &String) -> Result<Self, MetaCommandErr> {
        let meta_args = cmd.split(" ").collect::<Vec<&str>>();
        let meta_cmdt = meta_args[0];
        if &meta_cmdt[..1] == "." {
            if &meta_cmdt[1..5] == "exit" {
                return Ok(Self::Exit);
            }
            if &meta_cmdt[1..5] == "open" {
                return Ok(Self::Open(meta_args[1].trim().to_string()));
            }
            Err(MetaCommandErr::Unknown)
        } else {
            Err(MetaCommandErr::NotMetaCommand)
        }
    }

    pub fn execute_meta_command(
        cmd: &String,
        per_db: &mut Option<PersistantDatabase>,
        per_db_name: &mut Option<String>,
    ) -> Result<(), MetaCommandErr> {
        match Self::parse_meta_command(cmd) {
            Ok(command) => match command {
                MetaCommandType::Exit => {
                    println!("Goodbye!");
                    std::process::exit(0);
                }
                MetaCommandType::Open(dname) => match PersistantDatabase::open_db(dname.as_str()) {
                    Ok(existing_db) => {
                        *per_db = Some(existing_db);
                        *per_db_name = Some(dname);
                    }
                    Err(err) => match err {
                        PersistantDatabaseErr::UnknownDbErr => (),
                        PersistantDatabaseErr::ReadingErr => {
                            let new_per_db = PersistantDatabase::create_persistant_database();
                            PersistantDatabase::save_db(&dname, &new_per_db);
                            *per_db = Some(new_per_db);
                            *per_db_name = Some(dname);
                        }
                    },
                },
            },
            Err(err) => return Err(err),
        }
        Ok(())
    }
}
