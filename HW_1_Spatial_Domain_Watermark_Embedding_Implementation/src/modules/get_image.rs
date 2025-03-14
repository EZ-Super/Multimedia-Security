use std::collections::HashMap;

use image::{open, GenericImageView};

use crate::modules::point::Point;


pub struct BaseImage{
    path : String,
    width : u32,
    height : u32,
    pixel: HashMap<Point,u8>,
    color_type:image::ColorType,
    dyamic_image: image::DynamicImage,
}

impl BaseImage{
    pub fn new(path:String,color_type:image::ColorType)->BaseImage{
        let image = match open(&path){
            Ok(image) => image,
            Err(error) => panic!("Error loading image: {}", error)
        };
        let (width, height) = image.dimensions();
        BaseImage{
            path,
            width,
            height,
            pixel:HashMap::new(),
            color_type,
            dyamic_image:image,
        }

    }
    pub fn get_pixel(&mut self){
        for x in 0..self.width{
            for y in 0..self.height{
                let [r,g,b,a] = self.dyamic_image.get_pixel(x,y).0;
                println!("r:{} g:{} b:{} a:{}",r,g,b,a);
            }
        }
    }
}
