const ROWS_PER_PAGE: i32 = 1000;
const TABLE_MAX_PAGES: i32 = 5;

pub enum MetaCommandResult {
    Success,
    Unknown,
}

pub enum PrepareResult {
    Success,
    Unknown,
    SyntaxError,
}

#[derive(Debug, PartialEq)]
pub enum ExecuteResult {
    Success,
    TableFull,
}

pub enum StatementType {
    Insert,
    Select,
}

#[derive()]
pub struct Statement {
    pub stype: Option<StatementType>,
    pub row: Option<Row>,
}

#[derive(Debug, Clone)]
pub struct Row {
    pub id: i32,
    pub username: String,
    pub email: String,
}

pub struct Table {
    pub num_rows: i32,
    pub pages: [Page; TABLE_MAX_PAGES as usize],
}

#[derive(Debug)]
pub struct Page {
    pub rows: [Option<Row>; ROWS_PER_PAGE as usize],
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
            statement.stype = Some(StatementType::Insert);
            let args: Vec<&str> = cmd[7..].split(" ").collect();
            if args.len() < 3 {
                return PrepareResult::SyntaxError;
            }
            statement.row = Some(Row {
                id: std::str::FromStr::from_str(args[0]).unwrap(),
                username: args[2].to_owned(),
                email: args[1].to_owned(),
            });
            return PrepareResult::Success;
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
    fn execute_insert(statement: &Statement, table: &mut Table) -> ExecuteResult {
        if table.num_rows >= ROWS_PER_PAGE * TABLE_MAX_PAGES {
            return ExecuteResult::TableFull;
        }

        if let Some(row) = &statement.row {
            let coords = Row::row_slot(row.id);
            table.pages[coords.1].rows[coords.0] = Some(row.to_owned());
        }

        table.num_rows += 1;
        ExecuteResult::Success
    }

    fn execute_select(table: &Table) -> ExecuteResult {
        for page in &table.pages {
            for row in &page.rows {
                match row {
                    Some(y) => print!("{} {} {}", y.id, y.email, y.username),
                    None => (),
                }
            }
        }
        ExecuteResult::Success
    }
}

impl Table {
    pub fn create_table() -> Table {
        let pages: [Page; 5] = (0..TABLE_MAX_PAGES)
            .into_iter()
            .map(|_| Page {
                rows: (0..ROWS_PER_PAGE)
                    .into_iter()
                    .map(|_| None)
                    .collect::<Vec<Option<Row>>>()
                    .try_into()
                    .expect("failed to convert"),
            })
            .collect::<Vec<Page>>()
            .try_into()
            .expect("failed to convert");

        Table { num_rows: 0, pages }
    }
}

impl Row {
    fn row_slot(row_num: i32) -> (usize, usize) {
        let page_num = (row_num / ROWS_PER_PAGE) as usize;
        let row_num = (row_num % TABLE_MAX_PAGES) as usize;
        (row_num, page_num)
    }
}
