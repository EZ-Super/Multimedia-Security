use image::ColorType;
use std::env;

use jpeg_decoder::Decoder;
use std::fs::File;
use std::io::BufReader;


mod modules;
use modules::get_image;

fn main() {
    dotenvy::dotenv().expect("Failed to read .env file");
/*
    let path = env::var("elaine").expect("elaine read error");
    let mut image = get_image::BaseImage::new(path,ColorType::L8);
    let _ = image.get_pixel();
    let _ = image.show_pixel();
    println!("Hello, world!");
 */

 
    let nfu_path = env::var("nfuwm").expect("nfu read error");

    let mut nfu = get_image::BaseImage::new(nfu_path,ColorType::L8);
    let _ = nfu.get_pixel();
    let _ = nfu.show_pixel();
    println!("Hello, world!");
/*
    let file = File::open(env::var("nfuwm").expect("nfuwm read error")).unwrap();
    let mut decoder = Decoder::new(BufReader::new(file));
    let pixel = decoder.decode().expect("error decoding image");
    let metadata = decoder.info().unwrap();

    let width = metadata.width as usize;
    let height = metadata.height as usize;

    for x in 0..width {
        for y in 0..height {
            let index = (y*width + x) * 3;
            let r = pixel[index];
            let g = pixel[index + 1];
            let b = pixel[index + 2];
            println!("x:{} y:{} r:{} g:{} b:{}", x, y, r, g, b);
        }
    }
     */

}
