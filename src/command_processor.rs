use crate::db_arch::*;

pub enum MetaCommandResult {
    Success,
    Unknown,
}

#[derive(Debug, PartialEq, Clone)]
pub enum PrepareResult {
    Success(Table),
    Unknown(String),
    SyntaxError(String),
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
    pub row: Option<Row>,
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
            return StatementType::prepare_insert(cmd, statement, None);
        }

        if &cmd[..6] == "select" {
            statement.stype = Some(StatementType::Select);
            return PrepareResult::Success(db.tables[0].clone().unwrap());
        }

        if &cmd[..5] == "create" {
            statement.stype = Some(StatementType::Create);
            return PrepareResult::Success(db.tables[0].clone().unwrap());
        }

        PrepareResult::Unknown(cmd.to_owned())
    }
    pub fn execute_statement(statement: &Statement, table: &mut Table) -> Option<ExecuteResult> {
        // let db
        if let Some(stype) = statement.stype {
            match statement.stype {
                StatementType::Insert => Some(StatementType::execute_insert(statement, table)),
                StatementType::Select => Some(StatementType::execute_select(table)),
                StatementType::Create => Some(StatementType::execute_create(statement, db)),
                // _ => None,
            }
        }
        // None
    }
}

impl StatementType {
    fn prepare_insert(
        cmd: &String,
        statement: &mut Statement,
        table: Option<Table>,
    ) -> PrepareResult {
        match table {
            Some(table) => {
                statement.stype = Some(StatementType::Insert);
                let args: Vec<&str> = cmd[7..].split(" ").collect();
                if args.len() < 2 {
                    return PrepareResult::SyntaxError(format!("{:?}", args));
                }

                statement.row = Some(Row {
                    rtype: RowType::Insert(table.num_rows, args[0].to_owned(), args[1].to_owned()),
                });
                return PrepareResult::Success(table);
            }
            None => PrepareResult::NoExistingTable,
        }
    }

    fn execute_insert(statement: &Statement, table: &mut Table) -> ExecuteResult {
        if let Some(row) = &statement.row {
            table.rows.push(Some(row.to_owned()));
        }

        table.num_rows += 1;
        ExecuteResult::Success
    }

    fn execute_select(table: &Table) -> ExecuteResult {
        for row in &table.rows {
            if let Some(row) = row {
                match &row.rtype {
                    RowType::Insert(id, email, username) => {
                        print!("{} {} {}", id, email, username);
                    }
                    _ => (),
                }
            }
        }
        ExecuteResult::Success
    }

    fn execute_create(statement: &Statement, db: &mut Database) -> ExecuteResult {
        if let Some(row) = &statement.row {
            match &row.rtype {
                RowType::Create(_, tname) => {
                    let table = Table::create_table(tname.to_owned());
                    db.tables.push(Some(table));
                }
                _ => (),
            }
        }
        ExecuteResult::Success
    }
}
