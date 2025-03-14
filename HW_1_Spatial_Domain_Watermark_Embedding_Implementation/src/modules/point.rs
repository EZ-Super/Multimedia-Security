use std::hash::{Hash, Hasher};


#[derive(Debug,Eq,PartialEq)]
pub struct Point{
    pub x: i32,
    pub y: i32,
}

impl Hash for Point{
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.x.hash(state);
        self.y.hash(state);
    }
}