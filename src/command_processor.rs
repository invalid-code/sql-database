use crate::db_arch::*;

pub enum StatementErr {
    Unknown,
    SyntaxErr,
}

pub enum StatementResult {
    Err(StatementErr),
    Success,
}

pub enum ExecuteErr {
    DatabaseDoesNotExist,
    TableDoesNotExist,
}

pub enum ExecuteResult {
    Success,
    Err(ExecuteErr),
}

pub enum StatementType {
    Insert(i32, String, String, String, String),
    Select(String, String),
    Create(String, String, Option<String>),
}

impl StatementType {
    pub fn execute_create(
        dstruct: String,
        dstructn: String,
        db: Option<Database>,
        per_db: &mut PersistantDatabase,
    ) -> ExecuteResult {
        if dstruct == "db" {
            per_db.dbs.push(Some(Database::create_database()));
            per_db
                .index
                .insert(dstructn.clone(), per_db.num_dbs.clone());
            per_db.num_dbs += 1;
        }
        if dstruct == "table" {
            let table = Table::create_table(dstruct);
            match db {
                Some(mut db) => {
                    db.tables.push(Some(table));
                    db.index.insert(dstructn.clone(), db.num_tables.clone());
                    db.num_tables += 1;
                }
                None => {
                    return ExecuteResult::Err(ExecuteErr::DatabaseDoesNotExist);
                }
            }
        }
        ExecuteResult::Success
    }

    pub fn execute_insert(
        id: i32,
        email: String,
        username: String,
        dname: String,
        tname: String,
        per_db: &PersistantDatabase,
    ) -> ExecuteResult {
        match per_db.get_db(&dname) {
            Some(db) => match db.get_table(tname) {
                Some(mut table) => {
                    let row = Row {
                        id,
                        email,
                        username,
                    };
                    table.rows.push(Some(row));
                }
                None => return ExecuteResult::Err(ExecuteErr::TableDoesNotExist),
            },
            None => return ExecuteResult::Err(ExecuteErr::DatabaseDoesNotExist),
        }
        ExecuteResult::Success
    }

    pub fn execute_select(
        dname: String,
        tname: String,
        per_db: &PersistantDatabase,
    ) -> ExecuteResult {
        match per_db.get_db(&dname) {
            Some(db) => match db.get_table(tname) {
                Some(table) => {
                    for row in &table.rows {
                        match row {
                            Some(row) => println!("{} {} {}", row.id, row.email, row.username),
                            None => (),
                        }
                    }
                }
                None => return ExecuteResult::Err(ExecuteErr::TableDoesNotExist),
            },
            None => return ExecuteResult::Err(ExecuteErr::DatabaseDoesNotExist),
        }
        ExecuteResult::Success
    }
}

pub enum Statement {
    Result(StatementResult),
    Statement(StatementType),
}

impl Statement {
    pub fn parse_statement(cmd: &String) -> Self {
        if &cmd[..6] == "insert" {
            let args = cmd[7..].split(" ").collect::<Vec<&str>>();
            if args.len() < 5 {
                return Self::Result(StatementResult::Err(StatementErr::SyntaxErr));
            }
            return Self::Statement(StatementType::Insert(
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
                return Self::Result(StatementResult::Err(StatementErr::SyntaxErr));
            }
            return Self::Statement(StatementType::Select(
                std::str::FromStr::from_str(args[0]).unwrap(),
                args[1].to_owned(),
            ));
        }

        if &cmd[..5] == "create" {
            let args = cmd[6..].split(" ").collect::<Vec<&str>>();
            if args.len() < 3 {
                return Self::Statement(StatementType::Create(
                    args[0].to_owned(),
                    args[1].to_owned(),
                    Some(args[2].to_owned()),
                ));
            }
            if args.len() < 2 {
                return Self::Statement(StatementType::Create(
                    args[0].to_owned(),
                    args[1].to_owned(),
                    None,
                ));
            }
            return Self::Result(StatementResult::Err(StatementErr::SyntaxErr));
        }

        Self::Result(StatementResult::Err(StatementErr::Unknown))
    }
}
