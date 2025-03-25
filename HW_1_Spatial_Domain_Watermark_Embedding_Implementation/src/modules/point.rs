use std::hash::{Hash, Hasher};


#[derive(Debug,Eq,PartialEq,Clone)]
pub struct Point{
    pub x: u32,
    pub y: u32,
}

impl Hash for Point{
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.x.hash(state);
        self.y.hash(state);
    }

}

#[derive(Debug,Clone)]
pub struct RGBPixel {
    pub r: u8,
    pub g: u8,
    pub b: u8,
    pub a: u8,
}

impl RGBPixel {
    pub fn new(r: u8, g: u8, b: u8, a: u8) -> RGBPixel {
        RGBPixel {
            r,
            g,
            b,
            a,
        }
    }
    pub fn get_rgb(&self)->[u8;4]{
        [self.r,self.g,self.b,self.a]
    }
}