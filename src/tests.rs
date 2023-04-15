#[cfg(test)]
mod tests {
    use crate::command_processor::*;
    #[test]
    fn test_insert() {
        let mut table = Table::create_table();
        let row = Row {
            id: 1,
            username: "admin@web.net".to_owned(),
            email: "admin".to_owned(),
        };
        let statement = Statement {
            stype: Some(StatementType::Insert),
            row: Some(row),
        };

        if let Some(res) = Statement::execute_statement(&statement, &mut table) {
            assert_eq!(res, ExecuteResult::Success);
        }
    }
}
