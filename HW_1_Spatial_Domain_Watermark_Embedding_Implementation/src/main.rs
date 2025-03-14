use image::ColorType;
use std::env;

mod modules;
use modules::get_image;

fn main() {
    dotenvy::dotenv().expect("Failed to read .env file");
    let path = env::var("elaine").expect("elaine read error");
    let mut image = get_image::BaseImage::new(path,ColorType::L8);
    let _ = image.get_pixel();

    println!("Hello, world!");
}
