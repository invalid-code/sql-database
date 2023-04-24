use std::collections::HashMap;
use std::fs::{read_to_string, write};

#[derive(Debug, PartialEq, Clone)]
pub enum RowType {
    Insert,
    Select,
    Create,
}

#[derive(Debug, Default, Clone, PartialEq)]
pub struct Table {
    pub num_rows: i32,
    pub rows: Vec<Option<RowType>>,
    pub name: String,
}

impl Table {
    pub fn create_table(name: String) -> Self {
        Table {
            num_rows: 0,
            rows: vec![],
            name,
        }
    }
}

#[derive(Debug, Clone, Default, PartialEq)]
pub struct Database {
    pub tables: Vec<Option<Table>>,
    pub index: HashMap<String, i32>,
    pub num_tables: i32,
}

impl Database {
    pub fn create_database() -> Self {
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

pub struct PersistantDatabase {
    pub dbs: Vec<Option<Database>>,
    pub index: HashMap<String, i32>,
    pub num_dbs: i32,
}

impl PersistantDatabase {
    pub fn create_persistant_database() -> Self {
        PersistantDatabase {
            dbs: vec![],
            index: HashMap::new(),
            num_dbs: 0,
        }
    }

    pub fn get_db(&self, dname: &String) -> Option<Database> {
        match self.index.get(dname) {
            Some(i) => self.dbs[i.to_owned() as usize].clone(),
            None => None,
        }
    }
}
