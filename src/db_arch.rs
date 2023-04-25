use std::collections::HashMap;
// use std::fs::{read_to_string, write};

#[derive(Debug, PartialEq, Clone)]
pub struct Row {
    pub id: i32,
    pub email: String,
    pub username: String,
}

#[derive(Debug, Default, Clone, PartialEq)]
pub struct Table {
    pub num_rows: i32,
    pub rows: Vec<Option<Row>>,
    pub tname: String,
}

impl Table {
    pub fn create_table(tname: String) -> Self {
        Table {
            num_rows: 0,
            rows: vec![],
            tname,
        }
    }
}

#[derive(Debug, Clone, Default, PartialEq)]
pub struct Database {
    pub tables: Vec<Option<Table>>,
    pub index: HashMap<String, i32>,
    pub num_tables: i32,
    pub dname: String,
}

impl Database {
    pub fn create_database(dname: String) -> Self {
        Database {
            tables: vec![],
            index: HashMap::new(),
            num_tables: 0,
            dname,
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

    pub fn push_db(&mut self, table: &Table) {
        self.tables.push(Some(table.to_owned()));
        self.index
            .insert(table.tname.clone(), self.num_tables.clone());
        self.num_tables += 1;
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

    pub fn push_per_db(&mut self, db: &Database) {
        self.dbs.push(Some(db.to_owned()));
        self.index.insert(db.dname.clone(), self.num_dbs.clone());
        self.num_dbs += 1;
    }
}
