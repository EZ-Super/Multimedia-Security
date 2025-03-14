use std::collections::HashMap;

use image::{open, GenericImageView};

use crate::modules::point::Point;

use super::point;

#[allow(dead_code)]
pub struct BaseImage{
    path : String,
    pub width : u32,
    pub height : u32,
    pub pixel: HashMap<Point,point::RGPPixel>,
    color_type:image::ColorType,
    dyamic_image: image::DynamicImage,
}

impl BaseImage{
    pub fn new(path:String,color_type:image::ColorType)->BaseImage{
        println!("path:{}",path);
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
                //println!("r:{} g:{} b:{} a:{}",r,g,b,a);

                let point = Point{x,y};
                let pixel = point::RGPPixel::new(r,g,b,a);
                self.pixel.insert(point,pixel );
            }
        }
    }
    #[allow(unused)]
    pub fn show_pixel(&self){
        for x in 0..self.width{
            for y in 0..self.height{
                let point = Point{x,y};
                let pixel = self.pixel.get(&point).unwrap();

                let [r,g,b,a] = pixel.get_rgb();
                println!("r:{} g:{} b:{} a:{}",r,g,b,a);
            }
        }
    }
}
