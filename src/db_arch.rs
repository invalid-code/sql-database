use std::collections::HashMap;
use std::fs::{read_to_string, write};

#[derive(Debug, PartialEq, Clone)]
pub enum RowType {
    Insert(i32, String, String, String),
    Select(i32, String, String, String),
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

struct PersistantDatabase {
    dbs: Vec<Database>,
}

impl PersistantDatabase {
    fn read_file(path: &String) -> String {
        loop {
            match read_to_string(path.clone()) {
                Ok(contents) => {
                    let db_file = contents;
                    return db_file;
                }
                Err(_) => {
                    write(&path, "").expect("couldn create database");
                }
            }
        }
    }

    fn write_file(path: &String, contents: String) {
        write(path, contents).expect("couldn create database");
    }

    fn read_database(path: &String) -> Self {
        let file = read_file(path);
    }
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

    pub fn get_table(self: &Self, tname: String) -> Option<Table> {
        if let Some(tindex) = self.index.get(&tname) {
            if let Some(table) = &self.tables[tindex.to_owned() as usize] {
                return Some(table.to_owned());
            }
        }
        None
    }

    pub fn contains_tables(self: &Self) -> bool {
        if self.tables.len() < 1 {
            false
        } else {
            true
        }
    }
}
