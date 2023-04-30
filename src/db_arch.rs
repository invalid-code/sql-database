use super::command_processor::ExecuteErr;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::{read_to_string, write};

#[derive(Debug, PartialEq, Clone, Serialize, Deserialize)]
pub struct Row {
    pub id: i32,
    pub email: String,
    pub username: String,
}

#[derive(Debug, Default, Clone, PartialEq, Serialize, Deserialize)]
pub struct Table {
    pub num_rows: i32,
    pub rows: Vec<Row>,
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

#[derive(Debug, Clone, Default, PartialEq, Serialize, Deserialize)]
pub struct Database {
    pub tables: Vec<Table>,
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

    pub fn get_table(&mut self, tname: &str) -> Option<Table> {
        match self.index.get(tname) {
            Some(tindex) => Some(self.tables[tindex.to_owned() as usize].clone()),
            None => None,
        }
    }

    pub fn contains_tables(&self) -> bool {
        if self.tables.len() < 1 {
            false
        } else {
            true
        }
    }

    pub fn push_table(&mut self, table: Table) {
        self.tables.push(table.to_owned());
        self.index
            .insert(table.tname.clone(), self.num_tables.clone());
        self.num_tables += 1;
    }
}

#[derive(Debug, Default, Serialize, Deserialize)]
pub struct PersistantDatabase {
    pub dbs: Vec<Database>,
    pub index: HashMap<String, i32>,
    pub num_dbs: i32,
}

pub enum PersistantDatabaseErr {
    UnknownDbErr,
    ReadingErr,
}

impl PersistantDatabase {
    pub fn create_persistant_database() -> Self {
        PersistantDatabase {
            dbs: vec![],
            index: HashMap::new(),
            num_dbs: 0,
        }
    }

    /// returns a copy of the database if it exists
    pub fn get_db(&mut self, dname: &str) -> Option<Database> {
        match self.index.get(dname) {
            Some(dindex) => Some(self.dbs[dindex.to_owned() as usize].clone()),
            None => None,
        }
    }

    /// push a table to the db
    pub fn push_table(&mut self, dname: &str, table: Table) -> Result<(), ExecuteErr> {
        match self.index.get(dname) {
            Some(dindex) => {
                self.dbs[dindex.to_owned() as usize].push_table(table);
            }
            None => return Err(ExecuteErr::DatabaseDoesNotExist),
        }
        Ok(())
    }

    pub fn push_row(&mut self, dname: &str, tname: &str, row: Row) -> Result<(), ExecuteErr> {
        match self.index.get(dname) {
            Some(dindex) => {
                let db = &mut self.dbs[dindex.to_owned() as usize];
                match db.index.get(tname) {
                    Some(tindex) => {
                        db.tables[tindex.to_owned() as usize].rows.push(row);
                        db.tables[tindex.to_owned() as usize].num_rows += 1;
                    }
                    None => return Err(ExecuteErr::TableDoesNotExist),
                }
            }
            None => return Err(ExecuteErr::DatabaseDoesNotExist),
        }
        Ok(())
    }

    pub fn push_db(&mut self, db: &Database) {
        self.dbs.push(db.to_owned());
        self.index.insert(db.dname.clone(), self.num_dbs.clone());
        self.num_dbs += 1;
    }

    pub fn open_db(file: &str) -> Result<Self, PersistantDatabaseErr> {
        match read_file(file) {
            Ok(db) => match serde_json::from_str(&db) {
                Ok(db) => Ok(db),
                Err(_) => Err(PersistantDatabaseErr::UnknownDbErr),
            },
            Err(_) => Err(PersistantDatabaseErr::ReadingErr),
        }
    }

    pub fn save_db(file: &str, db: &PersistantDatabase) -> Result<(), std::io::Error> {
        write(file, serde_json::to_string(db).unwrap())
    }
}

pub fn write_file(file: &str, db: &str) -> Result<(), std::io::Error> {
    write(file, db)
}

pub fn read_file(file: &str) -> Result<String, std::io::Error> {
    read_to_string(file)
}
