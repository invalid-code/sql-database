use cucumber::World;

#[derive(World, Debug, Default)]
struct DataWorld {}

fn main() {
    futures::executor::block_on(DataWorld::run("tests/features/data_persistance.feature"));
}
