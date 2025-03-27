use flexi_logger::detailed_format;
use flexi_logger::FileSpec;
use image::ColorType;
use modules::embed_watermark::Watermark;
use std::env;
use log::{info,warn,error};
use flexi_logger::Logger;



mod modules;
use modules::get_image;
use modules::embed_watermark;
use modules::psnr;


fn main() {
    dotenvy::dotenv().expect("Failed to read .env file");

    let logger = Logger::try_with_str("info")
        .unwrap()
        .log_to_file(
            FileSpec::default()
                .directory("log")
                .basename("result")
                .suffix("log")
        )
        .write_mode(flexi_logger::WriteMode::BufferAndFlush)
        .duplicate_to_stderr(flexi_logger::Duplicate::All)
        .format_for_files(detailed_format)
        .start()
        .unwrap();
    info!("Start watermark embedding");

    
    let path = env::var("elaine").expect("elaine read error");
    let mut host_image = match get_image::BaseImage::new(path.clone(),ColorType::L8){
        Ok(host_image) => host_image,
        Err(error) => {
            panic!("{}",error);
        }
    };
    let _ = host_image.get_pixel();
    let new_image = embed_watermark::HostImage::new(host_image.clone());

    let random_number_embeding_50 = embed_watermark::RandomNumber::new(0.5);
    let random_number_embeding_50_file_name = new_image.embed_image_with_random_number(random_number_embeding_50, 0);
    let mut random_50_file = get_image::BaseImage::new(random_number_embeding_50_file_name.unwrap().clone(),ColorType::L8).unwrap();
    let _ = random_50_file.get_pixel();
    let mut random_number_embeding = psnr::PSNR::new(random_50_file,host_image.clone());
    let random_50_psnr = format!("random number PSNR 0.5: {}",random_number_embeding.calculate_psnr().get_psnr());
    println!("{}",random_50_psnr);
    info!("{}",random_50_psnr);


    let random_number_embeding_200 = embed_watermark::RandomNumber::new(2.0);
    let random_number_embeding_200_file_name = new_image.embed_image_with_random_number(random_number_embeding_200, 0);
    let mut random_200_file = get_image::BaseImage::new(random_number_embeding_200_file_name.unwrap().clone(),ColorType::L8).unwrap();
    let _ = random_200_file.get_pixel();
    let mut random_number_embeding = psnr::PSNR::new(random_200_file,host_image.clone());
    let random_200_psnr = format!("random number PSNR 2.0: {}",random_number_embeding.calculate_psnr().get_psnr());
    println!("{}",random_200_psnr);
    info!("{}",random_200_psnr);
 
    let nfu_path = env::var("nfuwm").expect("nfu read error");
    let mut nfu = match get_image::BaseImage::new(nfu_path,ColorType::L8){
        Ok(nfu) => nfu,
        Err(error) => {
            panic!("{}",error);
        }
    };
    let _ = nfu.get_pixel();
    let watermark = Watermark::new(&nfu);



    let _ = new_image.embed_image(watermark.clone(),  1, 1, 0);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 1);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 2);

    let embedded_1x1_3_file_name = new_image.embed_image(watermark.clone(),  1, 1, 3);
    let mut embedded_1x1_3_file = get_image::BaseImage::new(embedded_1x1_3_file_name.unwrap().clone(),ColorType::L8).unwrap();
    let _ = embedded_1x1_3_file.get_pixel();
    let mut psnr_1x1_3 = psnr::PSNR::new(embedded_1x1_3_file,host_image.clone());
    let embedded_1x1_3_psnr = format!("1x1_3 PSNR: {}",psnr_1x1_3.calculate_psnr().get_psnr());
    println!("{}",embedded_1x1_3_psnr);
    info!("{}",embedded_1x1_3_psnr);

    let _ = new_image.embed_image(watermark.clone(),  1, 1, 4);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 5);
    let _ = new_image.embed_image(watermark.clone(),  1, 1, 6);

    let embedded_1x1_7_file_name = new_image.embed_image(watermark.clone(), 1, 1, 7);
    let mut embedded_1x1_7_file = get_image::BaseImage::new(embedded_1x1_7_file_name.unwrap().clone(),ColorType::L8).unwrap();
    let _ = embedded_1x1_7_file.get_pixel();
    let mut psnr_1x1_7 = psnr::PSNR::new(embedded_1x1_7_file,host_image.clone());
    let embedded_1x1_7_psnr = format!("1x1_7 PSNR: {}",psnr_1x1_7.calculate_psnr().get_psnr());
    println!("{}",embedded_1x1_7_psnr);
    info!("{}",embedded_1x1_7_psnr);


    let _ = new_image.embed_image(watermark.clone(),  2, 1, 0);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 1);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 2);
    let embedded_2x1_3_file_name = new_image.embed_image(watermark.clone(),  2, 1, 3);
    let mut embedded_2x1_3_file = get_image::BaseImage::new(embedded_2x1_3_file_name.unwrap().clone(),ColorType::L8).unwrap();
    let _ = embedded_2x1_3_file.get_pixel();
    let mut psnr_2x2_3 = psnr::PSNR::new(embedded_2x1_3_file,host_image.clone());
    let embedded_2x1_3_psnr = format!("2x1_3 PSNR: {}",psnr_2x2_3.calculate_psnr().get_psnr());
    println!("{}",embedded_2x1_3_psnr);
    info!("{}",embedded_2x1_3_psnr);

    let _ = new_image.embed_image(watermark.clone(),  2, 1, 4);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 5);
    let _ = new_image.embed_image(watermark.clone(),  2, 1, 6);
    let embedded_2x1_7_file_name = new_image.embed_image(watermark.clone(),  2, 1, 7);
    let mut embedded_2x1_7_file = get_image::BaseImage::new(embedded_2x1_7_file_name.unwrap().clone(),ColorType::L8).unwrap();
    let _ = embedded_2x1_7_file.get_pixel();
    let mut psnr_2x2_7 = psnr::PSNR::new(embedded_2x1_7_file,host_image.clone());
    let embedded_2x1_7_psnr = format!("2x1_7 PSNR: {}",psnr_2x2_7.calculate_psnr().get_psnr());
    println!("{}",embedded_2x1_7_psnr);
    info!("{}",embedded_2x1_7_psnr);
    
    logger.flush();


}
