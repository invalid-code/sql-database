use crate::db_arch::*;

#[derive(Debug)]
pub enum StatementErr {
    Unknown,
    SyntaxErr,
}

pub enum ExecuteErr {
    DatabaseDoesNotExist,
    TableDoesNotExist,
}

// pub enum ExecuteResult {
//     Success,
//     Err(ExecuteErr),
// }

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
        db: Option<Database>,
        per_db: &mut PersistantDatabase,
    ) -> Result<(), ExecuteErr> {
        if dstruct == "db" {
            per_db.push_per_db(Database::create_database(dstructn));
        }
        if dstruct == "table" {
            match db {
                Some(mut db) => {
                    let table = Table::create_table(dstructn.clone());
                    db.push_db(table)
                }
                None => {
                    return Err(ExecuteErr::DatabaseDoesNotExist);
                }
            }
        }
        Ok(())
    }

    pub fn execute_insert(
        id: i32,
        email: String,
        username: String,
        dname: String,
        tname: String,
        per_db: &PersistantDatabase,
    ) -> Result<(), ExecuteErr> {
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
                None => return Err(ExecuteErr::TableDoesNotExist),
            },
            None => return Err(ExecuteErr::DatabaseDoesNotExist),
        }
        Ok(())
    }

    pub fn execute_select(
        dname: String,
        tname: String,
        per_db: &PersistantDatabase,
    ) -> Result<(), ExecuteErr> {
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
                None => return Err(ExecuteErr::TableDoesNotExist),
            },
            None => return Err(ExecuteErr::DatabaseDoesNotExist),
        }
        Ok(())
    }
}

impl StatementType {
    pub fn parse_statement(cmd: &String) -> Result<StatementType, StatementErr> {
        if &cmd[..6] == "insert" {
            let args = cmd[7..].split(" ").collect::<Vec<&str>>();
            if args.len() < 5 {
                return Err(StatementErr::SyntaxErr);
            }
            return Ok(StatementType::Insert(
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
                return Ok(StatementType::Create(
                    args[0].to_owned(),
                    args[1].to_owned(),
                    Some(args[2].to_owned()),
                ));
            }
            if args.len() == 2 {
                return Ok(StatementType::Create(
                    args[0].to_owned(),
                    args[1].to_owned(),
                    None,
                ));
            }
            return Err(StatementErr::SyntaxErr);
        }

        Err(StatementErr::Unknown)
    }
}
