use crate::db_arch::*;

pub enum MetaCommandResult {
    Success,
    Unknown,
}

#[derive(Debug, PartialEq, Clone)]
pub enum PrepareResult {
    Success,
    Unknown,
    SyntaxError,
    NoExistingTable,
}

#[derive(Debug, Clone, Default, PartialEq)]
pub enum ExecuteResult {
    #[default]
    Success,
}

#[derive(Debug, PartialEq, Clone)]
pub enum StatementType {
    Insert,
    Select,
    Create,
}

#[derive(Debug, Default, PartialEq, Clone)]
pub struct Statement {
    pub stype: Option<StatementType>,
    pub row: Option<RowType>,
}

pub fn execute_meta_command(cmd: &String) -> Option<MetaCommandResult> {
    if cmd == ".exit\n" {
        println!("Goodbye!");
        std::process::exit(0);
    } else {
        Some(MetaCommandResult::Unknown)
    }
}

impl Statement {
    pub fn prepare_statement(
        cmd: &String,
        statement: &mut Statement,
        db: &Database,
    ) -> PrepareResult {
        if &cmd[..6] == "insert" {
            return StatementType::prepare_insert(cmd, statement, db);
        }

        if &cmd[..6] == "select" {
            statement.stype = Some(StatementType::Select);
            return PrepareResult::Success;
        }

        if &cmd[..5] == "create" {
            statement.stype = Some(StatementType::Create);
            return PrepareResult::Success;
        }

        PrepareResult::Unknown
    }

    pub fn execute_statement(statement: &Statement, db: &mut Database) -> Option<ExecuteResult> {
        match statement.stype {
            Some(StatementType::Insert) => Some(StatementType::execute_insert(statement, db)),
            Some(StatementType::Select) => Some(StatementType::execute_select(statement, db)),
            Some(StatementType::Create) => Some(StatementType::execute_create(statement, db)),
            _ => None,
        }
    }
}

impl StatementType {
    fn prepare_insert(cmd: &String, statement: &mut Statement, db: &Database) -> PrepareResult {
        statement.stype = Some(StatementType::Insert);
        let args: Vec<&str> = cmd[7..].split(" ").collect();
        if args.len() < 3 {
            return PrepareResult::SyntaxError;
        }
        if let Some(table) = db.get_table(args[2].to_owned()) {
            statement.row = Some(RowType::Insert(
                table.num_rows,
                args[0].to_owned(),
                args[1].to_owned(),
                table.name,
            ));
        }
        PrepareResult::Success
    }

    fn execute_insert(statement: &Statement, db: &mut Database) -> ExecuteResult {
        if let Some(row) = &statement.row {
            if let RowType::Insert(_, _, _, tname) = row {
                if let Some(mut table) = db.get_table(tname.to_owned()) {
                    if let Some(row) = &statement.row {
                        table.rows.push(Some(row.to_owned()));
                    }
                    table.num_rows += 1;
                }
            }
        }
        ExecuteResult::Success
    }

    fn execute_select(statement: &Statement, db: &mut Database) -> ExecuteResult {
        if let Some(row) = &statement.row {
            if let RowType::Select(_, _, _, tname) = row {
                if let Some(table) = db.get_table(tname.to_owned()) {
                    for row in table.rows {
                        if let Some(row) = row {
                            if let RowType::Select(id, email, username, _) = row {
                                print!("{} {} {}", id, email, username);
                            }
                        }
                    }
                }
            }
        }
        ExecuteResult::Success
    }

    fn execute_create(statement: &Statement, db: &mut Database) -> ExecuteResult {
        if let Some(row) = &statement.row {
            if let RowType::Create(_, name) = row {
                let table = Table::create_table(name.to_owned());
                db.tables.push(Some(table));
                db.index.insert(name.to_owned(), db.num_tables + 1);
                db.num_tables += 1;
            }
        }
        ExecuteResult::Success
    }
}
