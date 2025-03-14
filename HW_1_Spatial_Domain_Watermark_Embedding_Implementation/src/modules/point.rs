use std::hash::{Hash, Hasher};


#[derive(Debug,Eq,PartialEq)]
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


pub struct RGPPixel {
    r: u8,
    g: u8,
    b: u8,
    a: u8,
}

impl RGPPixel {
    pub fn new(r: u8, g: u8, b: u8, a: u8) -> RGPPixel {
        RGPPixel {
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