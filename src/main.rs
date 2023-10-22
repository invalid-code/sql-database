use std::{
    collections::HashMap,
    error::Error,
    fmt::Display,
    fs::{self, File},
    io::{stdin, stdout, Write},
    process,
};

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
struct Row {
    id: i32,
    name: String,
    employed: bool,
}

impl Row {
    fn new(id: i32, name: String, employed: bool) -> Self {
        Self { id, name, employed }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct Table {
    rows: Vec<Row>,
}

impl Table {
    fn new() -> Self {
        Self { rows: Vec::new() }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct Database {
    tables: Vec<Table>,
    index: HashMap<String, i32>,
}

impl Database {
    fn new() -> Self {
        Self {
            tables: Vec::new(),
            index: HashMap::new(),
        }
    }

    fn get_table(&self, table_name: &str) -> &Table {
        &self.tables[self.index[table_name] as usize]
    }

    fn get_mut_table(&mut self, table_name: &str) -> &mut Table {
        &mut self.tables[self.index[table_name] as usize]
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct PersistantDatabase {
    databases: Vec<Database>,
    index: HashMap<String, i32>,
    database_name: String,
    current_database: String,
}

impl PersistantDatabase {
    fn new(database_name: String) -> Self {
        Self {
            databases: Vec::new(),
            index: HashMap::new(),
            current_database: String::new(),
            database_name,
        }
    }

    fn get_current_database(&self) -> &Database {
        &self.databases[self.index[&self.current_database] as usize]
    }

    fn get_mut_current_database(&mut self) -> &mut Database {
        &mut self.databases[*self.index.get(&self.current_database).unwrap() as usize]
    }

    fn save_db(&self, database_name: Option<String>) {
        match database_name {
            Some(database_name) => {
                let mut file = File::create(database_name).unwrap();
                file.write(serde_json::to_string(self).unwrap().as_bytes())
                    .unwrap();
            }
            None => {
                let mut file = File::create(self.database_name.clone()).unwrap();
                file.write(serde_json::to_string(self).unwrap().as_bytes())
                    .unwrap();
            }
        }
    }
}

#[derive(Debug)]
enum ParsingError {
    UnknownStatement,
    MissingKeyword(String),
}

impl Display for ParsingError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::UnknownStatement => write!(f, "unknown command"),
            Self::MissingKeyword(keyword) => write!(f, "missing keyword: {}", keyword),
        }
    }
}

impl Error for ParsingError {}

enum Statement {
    Insert(String, Vec<String>, Row),
    Select(Vec<String>, String),
    Create(String, String),
    Use(String, String),
}

impl Statement {
    fn parse(buf: &str) -> Result<Self, ParsingError> {
        if &buf[..6] == "insert" {
            let into_index = match buf.find("into") {
                Some(keyword) => keyword,
                None => return Err(ParsingError::MissingKeyword("into".to_owned())),
            };
            let value_index = match buf.find("value") {
                Some(keyword) => keyword,
                None => return Err(ParsingError::MissingKeyword("value".to_owned())),
            };
            let opening_parenthesis_index = match buf[..value_index].find("(") {
                Some(keyword) => keyword,
                None => return Err(ParsingError::MissingKeyword("(".to_owned())),
            };
            let table_name = buf[into_index + 5..opening_parenthesis_index - 1].to_owned();
            let row = {
                let args = buf[value_index + 7..buf.len() - 1]
                    .split(", ")
                    .collect::<Vec<&str>>();
                Row::new(
                    args[0].parse::<i32>().unwrap(),
                    args[1].parse::<String>().unwrap(),
                    args[2].parse::<bool>().unwrap(),
                )
            };
            let columns = buf[opening_parenthesis_index + 1..value_index - 2]
                .split(", ")
                .map(|column| column.to_owned())
                .collect::<Vec<String>>();
            Ok(Self::Insert(table_name, columns, row))
        } else if &buf[..6] == "select" {
            let from_keyword_index = match buf.find("from") {
                Some(keyword) => keyword,
                None => return Err(ParsingError::MissingKeyword("from".to_owned())),
            };
            let columns = buf[8..from_keyword_index - 1]
                .split(", ")
                .map(|column| column.to_owned())
                .collect::<Vec<String>>();
            let table_name = buf[from_keyword_index + 5..].to_owned();
            Ok(Self::Select(columns, table_name))
        } else if &buf[..6] == "create" {
            let args = buf[7..]
                .split(" ")
                .map(|arg| arg.to_owned())
                .collect::<Vec<String>>();
            Ok(Self::Create(args[0].to_owned(), args[1].to_owned()))
        } else if &buf[..3] == "use" {
            let args = buf[4..]
                .split(" ")
                .map(|arg| arg.to_owned())
                .collect::<Vec<String>>();
            Ok(Self::Use(args[0].clone(), args[1].clone()))
        } else {
            Err(ParsingError::UnknownStatement)
        }
    }

    fn execute(&self, per_db: &mut Option<PersistantDatabase>) -> Result<(), NoOpenDatabaseErr> {
        match per_db {
            Some(per_db) => match self {
                Self::Insert(table_name, _columns, value) => {
                    let current_database = per_db.get_mut_current_database();
                    let table = current_database.get_mut_table(table_name);
                    table.rows.push(value.to_owned());
                }
                Self::Select(_columns, table_name) => {
                    let current_database = per_db.get_current_database();
                    let table = current_database.get_table(table_name);
                    for row in table.rows.clone().into_iter() {
                        println!("{:?}", row);
                    }
                }
                Self::Create(struct_type, struct_name) => {
                    if struct_type == "database" {
                        per_db
                            .index
                            .insert(struct_name.to_owned(), per_db.databases.len() as i32);
                        per_db.databases.push(Database::new());
                    } else if struct_type == "table" {
                        let current_database = per_db.get_mut_current_database();
                        current_database
                            .index
                            .insert(struct_name.to_owned(), current_database.tables.len() as i32);
                        current_database.tables.push(Table::new());
                    }
                }
                Self::Use(struct_type, struct_name) => {
                    if struct_type == "database" {
                        match per_db.index.get(struct_name) {
                            Some(_) => per_db.current_database = struct_name.to_owned(),
                            None => return Err(NoOpenDatabaseErr),
                        }
                    } else if struct_type == "table" {
                    }
                }
            },
            None => (),
        }
        Ok(())
    }
}

#[derive(Debug)]
struct NoOpenDatabaseErr;

impl Display for NoOpenDatabaseErr {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "no open database")
    }
}

impl Error for NoOpenDatabaseErr {}

#[derive(Debug)]
struct UnknownCommandError;

impl Display for UnknownCommandError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "unknown command")
    }
}

impl Error for UnknownCommandError {}

enum Command {
    Exit,
    Open(String),
}

impl Command {
    fn parse(buf: &str) -> Result<Self, UnknownCommandError> {
        if &buf[..4] == "exit" {
            Ok(Self::Exit)
        } else if &buf[..4] == "open" {
            let database_name = buf[5..].to_owned();
            Ok(Self::Open(database_name))
        } else {
            Err(UnknownCommandError)
        }
    }

    fn execute(&self, per_db: &mut Option<PersistantDatabase>) {
        match self {
            Self::Exit => {
                if let Some(per_db) = per_db {
                    per_db.save_db(None);
                }
                println!("Goodbye!");
                process::exit(0);
            }
            Self::Open(database_name) => {
                match fs::read_to_string(database_name) {
                    Ok(file) => *per_db = Some(serde_json::from_str(&file).unwrap()),
                    Err(_) => {
                        *per_db = Some(PersistantDatabase::new(database_name.to_owned()));
                        if let Some(per_db) = per_db {
                            per_db.save_db(Some(database_name.to_owned()));
                        }
                    }
                };
            }
        }
    }
}

enum Input {
    Command,
    Statement,
    None,
}

impl Input {
    fn parse_input(buf: &str) -> Self {
        if &buf[..1] == "." {
            Self::Command
        } else if &buf[..1] == "\n" {
            Self::None
        } else {
            Self::Statement
        }
    }
}

fn read_input(buf: &mut String) {
    print!("> ");
    stdout().flush().unwrap();
    stdin().read_line(buf).unwrap();
}

fn main() {
    let mut per_db: Option<PersistantDatabase> = None;
    loop {
        let mut buf = String::new();
        read_input(&mut buf);
        buf = buf.trim().to_owned();

        let mut command: Option<Command> = None;
        let mut statement: Option<Statement> = None;

        match Input::parse_input(&buf) {
            Input::Command => match Command::parse(&buf[1..]) {
                Ok(parsed_command) => command = Some(parsed_command),
                Err(err) => println!("{}", err),
            },
            Input::Statement => match Statement::parse(&buf) {
                Ok(parsed_statement) => statement = Some(parsed_statement),
                Err(err) => println!("{}", err),
            },
            Input::None => (),
        }

        if let Some(command) = command {
            command.execute(&mut per_db);
        }

        if let Some(statement) = statement {
            match statement.execute(&mut per_db) {
                Ok(_) => (),
                Err(err) => println!("{}", err),
            };
        }
    }
}
