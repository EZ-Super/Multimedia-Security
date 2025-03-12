use image::{open, GenericImageView};
use std::env;

fn main() {
    dotenv::dotenv().ok();
    let path = env::var("elaine").expect("elaine read error");
    match open(path){
        Ok(image) => {
            if image.color() == image::ColorType::Rgb8 {
                let (width, height) = image.dimensions();
                for x in 0..width{
                    for y in 0..height{
                        let [r,g,b,a] = image.get_pixel(x, y).0;
                        println!("Pixel at ({}, {}): R: {}, G: {}, B: {}, A: {}", x, y, r, g, b, a);
                    }
                }
            } else {
                println!("Unsupported image color type {:?}", image.color());
            }
        },
        Err(error) => {
            println!("Error loading image: {}", error);
        }
    }



    println!("Hello, world!");
}
