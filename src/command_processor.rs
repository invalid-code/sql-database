use crate::db_arch::*;

pub enum StatementErr {
    Unknown,
    SyntaxErr,
}

pub enum StatementType {
    Insert(i32, String, String),
    Select(i32, String, String),
    Create(String, String, Option<String>),
    Err(StatementErr),
}

impl StatementType {
    pub fn parse_statement(cmd: &String) -> Self {
        if &cmd[..6] == "insert" {
            let args = cmd[7..].split(" ").collect::<Vec<&str>>();
            if args.len() < 3 {
                return Self::Err(StatementErr::SyntaxErr);
            }
            return Self::Insert(
                std::str::FromStr::from_str(args[0]).unwrap(),
                args[1].to_owned(),
                args[2].to_owned(),
            );
        }

        if &cmd[..6] == "select" {
            let args = cmd[7..].split(" ").collect::<Vec<&str>>();
            if args.len() < 3 {
                return Self::Err(StatementErr::SyntaxErr);
            }
            return Self::Select(
                std::str::FromStr::from_str(args[0]).unwrap(),
                args[1].to_owned(),
                args[2].to_owned(),
            );
        }

        if &cmd[..5] == "create" {
            let args = cmd[6..].split(" ").collect::<Vec<&str>>();
            if args.len() < 3 {
                return Self::Create(
                    args[0].to_owned(),
                    args[1].to_owned(),
                    Some(args[2].to_owned()),
                );
            }
            if args.len() < 2 {
                return Self::Create(args[0].to_owned(), args[1].to_owned(), None);
            }
            return Self::Err(StatementErr::SyntaxErr);
        }

        Self::Err(StatementErr::Unknown)
    }

    pub fn execute_create(
        dstruct: String,
        dstructn: String,
        db: Option<Database>,
        per_db: &mut PersistantDatabase,
    ) {
        if dstruct == "db" {
            per_db.dbs.push(Some(Database::create_database()));
            per_db.index.insert(dstructn, per_db.num_dbs.clone());
            per_db.num_dbs += 1;
        }
        if dstruct == "table" {
            let table = Table::create_table(dstruct);
        }
    }

    fn execute_insert(
        &self,
        id: i32,
        email: String,
        username: String,
        table: Option<Table>,
        db: Option<Database>,
    ) {
    }

    fn execute_select(
        &self,
        id: i32,
        email: String,
        username: String,
        table: Option<Table>,
        db: Option<Database>,
    ) {
    }
}
