use std::collections::HashMap;

#[derive(Debug, PartialEq, Clone)]
pub enum RowType {
    Insert(i32, String, String),
    Select(i32, String, String),
    Create(String, String),
}

#[derive(Debug, Default, Clone, PartialEq)]
pub struct Table {
    pub num_rows: i32,
    pub rows: Vec<Option<RowType>>,
    pub name: String,
}

#[derive(Debug, Clone, Default, PartialEq)]
pub struct Database {
    pub tables: Vec<Option<Table>>,
    pub index: HashMap<String, i32>,
    pub num_tables: i32,
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
        Database {
            tables: vec![],
            index: HashMap::new(),
            num_tables: 0,
        }
    }
}
