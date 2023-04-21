const TABLE_MAX_ROWS: i32 = 5000;

pub enum MetaCommandResult {
    Success,
    Unknown,
}

#[derive(Debug)]
pub enum PrepareResult {
    Success,
    Unknown,
    SyntaxError,
    InputTooLong,
}

#[derive(Debug, Clone, Default, PartialEq)]
pub enum ExecuteResult {
    #[default]
    Success,
    TableFull,
}

#[derive(Debug)]
pub enum StatementType {
    Insert,
    Select,
}

#[derive(Debug, Default)]
pub struct Statement {
    pub stype: Option<StatementType>,
    pub row: Option<Row>,
}

#[derive(Debug, Clone, Default)]
pub struct Row {
    pub id: i32,
    pub username: String,
    pub email: String,
}

#[derive(Debug, Default)]
pub struct Table {
    pub num_rows: i32,
    pub rows: Vec<Option<Row>>,
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
    pub fn prepare_statement(cmd: &String, statement: &mut Statement) -> PrepareResult {
        if &cmd[..6] == "insert" {
            return StatementType::prepare_insert(cmd, statement);
        }

        if &cmd[..6] == "select" {
            statement.stype = Some(StatementType::Select);
            return PrepareResult::Success;
        }
        PrepareResult::Unknown
    }
    pub fn execute_statement(statement: &Statement, table: &mut Table) -> Option<ExecuteResult> {
        match statement.stype {
            Some(StatementType::Insert) => Some(StatementType::execute_insert(statement, table)),
            Some(StatementType::Select) => Some(StatementType::execute_select(table)),
            _ => None,
        }
    }
}

impl StatementType {
    fn prepare_insert(cmd: &String, statement: &mut Statement) -> PrepareResult {
        statement.stype = Some(StatementType::Insert);
        let args: Vec<&str> = cmd[7..].split(" ").collect();
        if args.len() < 3 {
            return PrepareResult::SyntaxError;
        }

        if args[1].len() > 255 || args[2].len() > 30 {
            return PrepareResult::InputTooLong;
        }

        statement.row = Some(Row {
            id: std::str::FromStr::from_str(args[0]).unwrap(),
            username: args[2].to_owned(),
            email: args[1].to_owned(),
        });
        PrepareResult::Success
    }

    fn execute_insert(statement: &Statement, table: &mut Table) -> ExecuteResult {
        if table.num_rows >= TABLE_MAX_ROWS {
            return ExecuteResult::TableFull;
        }

        if let Some(row) = &statement.row {
            table.rows.push(Some(row.to_owned()));
        }

        table.num_rows += 1;
        ExecuteResult::Success
    }

    fn execute_select(table: &Table) -> ExecuteResult {
        for row in &table.rows {
            if let Some(row) = row {
                print!("{} {} {}", row.id, row.email, row.username);
            }
        }
        ExecuteResult::Success
    }
}

impl Table {
    pub fn create_table() -> Table {
        Table {
            num_rows: 0,
            rows: vec![],
        }
    }
}
