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
    let mut host_image = match get_image::BaseImage::new(path.clone(),ColorType::L8){
        Ok(host_image) => host_image,
        Err(error) => {
            panic!("{}",error);
        }
    };
    let _ = host_image.get_pixel();
    let new_image = embed_watermark::HostImage::new(host_image.clone());


/*
    let random_number_embeding_50 = embed_watermark::RandomNumber::new(0.5);
    let random_number_embeding_50_file = new_image.embed_image_with_random_number(random_number_embeding_50, 0);



    let random_number_embeding_200 = embed_watermark::RandomNumber::new(2.0);
    let random_number_embeding_200_file = new_image.embed_image_with_random_number(random_number_embeding_200, 0);

 */


    let nfu_path = env::var("nfuwm").expect("nfu read error");
    let mut nfu = match get_image::BaseImage::new(nfu_path,ColorType::L8){
        Ok(nfu) => nfu,
        Err(error) => {
            panic!("{}",error);
        }
    };
    let _ = nfu.get_pixel();
    let watermark = Watermark::new(&nfu);
    watermark.post_host_image();
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
