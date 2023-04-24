use crate::db_arch::*;

// pub enum MetaCommandResult {
//     Success,
//     Unknown,
// }

// #[derive(Debug, PartialEq, Clone)]
// pub enum PrepareResult {
//     Success,
//     Unknown,
//     SyntaxError,
//     NoExistingTable,
//     NoExistingDatabase,
// }

// #[derive(Debug, Clone, Default, PartialEq)]
// pub enum ExecuteResult {
//     #[default]
//     Success,
// }
pub enum StatementErr {
    Unknown,
    SyntaxErr,
}

pub enum StatementType {
    Insert(i32, String, String),
    Select(i32, String, String),
    Create(String, String),
    Err(StatementErr),
}

// #[derive(Debug, Default, PartialEq, Clone)]
// pub struct Statement {
//     pub stype: Option<StatementType>,
//     pub row: Option<RowType>,
// }

// pub fn execute_meta_command(cmd: &String) -> Option<MetaCommandResult> {
//     if cmd == ".exit\n" {
//         println!("Goodbye!");
//         std::process::exit(0);
//     } else {
//         Some(MetaCommandResult::Unknown)
//     }
// }

// impl Statement {
//     pub fn prepare_statement(
//         cmd: &String,
//         statement: &mut Statement,
//         db: Option<&Database>,
//     ) -> PrepareResult {
//         match db {
//             Some(db) => {
//                 if !db.contains_tables() {
//                     return PrepareResult::NoExistingTable;
//                 }

//                 if &cmd[..6] == "insert" {
//                     return StatementType::prepare_insert(cmd, statement, db);
//                 }

//                 if &cmd[..6] == "select" {
//                     statement.stype = Some(StatementType::Select);
//                     return PrepareResult::Success;
//                 }
//             }
//             None => return PrepareResult::NoExistingDatabase,
//         }
//         if &cmd[..5] == "create" {
//             statement.stype = Some(StatementType::Create);
//             return PrepareResult::Success;
//         }

//         PrepareResult::Unknown
//     }

//     pub fn execute_statement(statement: &Statement, db: &mut Database) -> Option<ExecuteResult> {
//         match statement.stype {
//             Some(StatementType::Insert) => Some(StatementType::execute_insert(statement, db)),
//             Some(StatementType::Select) => Some(StatementType::execute_select(statement, db)),
//             Some(StatementType::Create) => Some(StatementType::execute_create(statement, db)),
//             _ => None,
//         }
//     }
// }

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
            if args.len() < 2 {
                return Self::Err(StatementErr::SyntaxErr);
            }
            return Self::Create(args[0].to_owned(), args[1].to_owned());
        }

        Self::Err(StatementErr::Unknown)
    }

    // fn parse_args(cmd: &String, stype: &String) {}

    // fn prepare_insert(cmd: &String, statement: &mut Statement, db: &Database) -> PrepareResult {
    //     statement.stype = Some(StatementType::Insert);
    //     let args: Vec<&str> = cmd[7..].split(" ").collect();
    //     if args.len() < 3 {
    //         return PrepareResult::SyntaxError;
    //     }
    //     let table = db.get_table(args[2].to_owned()).unwrap();
    //     statement.row = Some(RowType::Insert(
    //         table.num_rows,
    //         args[0].to_owned(),
    //         args[1].to_owned(),
    //         table.name,
    //     ));
    //     PrepareResult::Success
    // }

    // fn execute_insert(statement: &Statement, db: &mut Database) -> ExecuteResult {
    //     let row = statement.row.as_ref().unwrap();
    //     if let RowType::Insert(_, _, _, tname) = row {
    //         let mut table = db.get_table(tname.to_owned()).unwrap();
    //         table.rows.push(Some(row.to_owned()));
    //     }
    //     ExecuteResult::Success
    // }

    // fn execute_select(statement: &Statement, db: &mut Database) -> ExecuteResult {
    //     let row = statement.row.as_ref().unwrap();
    //     if let RowType::Select(_, _, _, tname) = row {
    //         let table = db.get_table(tname.to_owned()).unwrap();
    //         for row in table.rows {
    //             if let Some(row) = row {
    //                 if let RowType::Select(id, email, username, _) = row {
    //                     print!("{} {} {}", id, email, username);
    //                 }
    //             }
    //         }
    //     }
    //     ExecuteResult::Success
    // }

    // fn execute_create(statement: &Statement, db: &mut Database) -> ExecuteResult {
    //     if let Some(row) = &statement.row {
    //         if let RowType::Create(_, name) = row {
    //             // if data_structure[]
    //             let table = Table::create_table(name.to_owned());
    //             db.tables.push(Some(table));
    //             db.index.insert(name.to_owned(), db.num_tables + 1);
    //             db.num_tables += 1;
    //         }
    //     }
    //     ExecuteResult::Success
    // }
}
