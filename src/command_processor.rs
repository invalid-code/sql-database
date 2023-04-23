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
    Success(Table),
}

#[derive(Debug)]
pub enum StatementType {
    Insert,
    Select,
    Create,
}

#[derive(Debug, Default)]
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
            return StatementType::prepare_insert(cmd, statement, db.tables[0]);
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
            Some(StatementType::Select) => Some(StatementType::execute_select(db)),
            Some(StatementType::Create) => Some(StatementType::execute_create(statement, db)),
            _ => None,
        }
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
                    return PrepareResult::SyntaxError;
                }

                statement.row = Some(Row {
                    rtype: RowType::Insert(table.num_rows, args[0].to_owned(), args[1].to_owned()),
                });
                return PrepareResult::Success;
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

    fn execute_select(db: &Database) -> ExecuteResult {
        for row in &table.rows {
            if let Some(row) = row {
                print!("{} {} {}", row.id, row.email, row.username);
            }
        }
        ExecuteResult::Success
    }

    fn execute_create(statement: &Statement, db: &Database) -> ExecuteResult {
        let table = Table::create_table();
        db.tables.push(Some(table));
        ExecuteResult::Success
    }
}
