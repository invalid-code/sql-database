#[derive(Debug, Clone, Default)]
pub struct Row {
    // pub id: i32,
    // pub username: String,
    // pub email: String,
    pub rtype: RowType,
}

#[derive(Debug, Clone)]
pub enum RowType {
    Insert(i32, String, String),
    Create(String, String),
}

impl Default for RowType {
    fn default() -> Self {
        Self::Create(String::from("table"), String::from("table"))
    }
}

#[derive(Debug, Default, Clone, PartialEq)]
pub struct Table {
    pub num_rows: i32,
    pub rows: Vec<Option<Row>>,
    pub name: String,
}

pub struct Database {
    pub tables: Vec<Option<Table>>,
}

impl Table {
    pub fn create_table(name: String) -> Table {
        Table {
            num_rows: 0,
            rows: vec![],
            name,
        }
    }
}

impl Database {
    pub fn create_database() -> Database {
        Database { tables: vec![] }
    }
}
