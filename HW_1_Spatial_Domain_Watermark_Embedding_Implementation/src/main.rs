use image::ColorType;
use modules::embed_watermark::Watermark;
use std::env;



mod modules;
use modules::get_image;
use modules::embed_watermark;
use modules::note;

fn main() {
    dotenvy::dotenv().expect("Failed to read .env file");

    let path = env::var("elaine").expect("elaine read error");
    let mut image = get_image::BaseImage::new(path,ColorType::L8);
    let _ = image.get_pixel();
//    let _ = image.show_pixel();
    println!("Hello, world!");
 

 

 
    let nfu_path = env::var("nfuwm").expect("nfu read error");

    let mut nfu = get_image::BaseImage::new(nfu_path,ColorType::L8);
    let _ = nfu.get_pixel();
//   let _ = nfu.show_pixel();
    println!("Hello, world!");


    let mut new_image = embed_watermark::HostImage::new(image);
    let watermark = Watermark::new(&nfu);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 0);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 0);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 1);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 2);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 3);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 4);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 5);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 6);
    let _ = new_image.embed_image(watermark.clone(), 1, 1, 7);

    let _ = new_image.embed_image(watermark.clone(),  2, 1, 0);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 1);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 2);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 3);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 4);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 5);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 6);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 7);



}
