use image::Rgba;

use crate::modules::get_image::BaseImage;
use crate::modules::point::RGPPixel;
use crate::modules::point::Point;

use std::collections::HashMap;

pub struct Watermark{
    watermark:BaseImage,
    width:u32,
    height:u32,
    pixel:HashMap<Point,RGPPixel>,
}

pub struct RandomNumber{
    width:u32,
    height:u32,
    embed_numer:Vec<i32>,
}
pub struct HostImage{
    base_image:BaseImage,
    embed_width:u32,
    embed_height:u32,
}

impl Watermark{
    pub fn new(base_image:&BaseImage)->Watermark{
        let clone_base_image = base_image.clone();
        let width = clone_base_image.width;
        let height = clone_base_image.height;
        let pixel = clone_base_image.pixel.clone();
        Watermark{
            watermark:clone_base_image,
            width,
            height,
            pixel:pixel,
        }
    }
}


impl HostImage{
    pub fn new(base_image:BaseImage)->HostImage{
        HostImage{
            base_image,
            embed_width: 0 as u32,
            embed_height: 0 as u32,
        }
    }
    pub fn embed_image(&self,watermark: Watermark,x_point:u32,y_point:u32,watermark_x_number:u32,watermark_y_number:u32,embed_bit:u8)->Result<BaseImage,String>{
        let host_image = self.base_image.clone();
        let mut image = host_image.dyamic_image.to_rgba8();
        let embed_width = watermark.width*watermark_x_number;
        let embed_height = watermark.height * watermark_y_number;

        if host_image.width < embed_width || host_image.height < embed_height{
            return Err("The watermark is too large to embed".to_string())
        }
        for x in 0..embed_width{
            for y in 0..embed_height{


                let x_pixel = x % watermark.width;
                let y_pixel = y % watermark.height;

                let watermark_pixel = match watermark.pixel.get(&Point{x:x_pixel,y:y_pixel}){
                    Some(pixel) => pixel,
                    None => return Err("Watermark pixel not found".to_string())
                };

                let host_pixel = match host_image.pixel.get(&Point{x:x,y:y}){
                    Some(pixel) => pixel,
                    None => return Err("Host pixel not found".to_string())
                };
                let mut embed_pixel = 0;
                if watermark_pixel.r == 255 && watermark_pixel.g == 255 && watermark_pixel.b == 255{
                    embed_pixel = 1 ;
                }
                embed_pixel =  embed_pixel <<embed_bit;
                let result_r = if embed_pixel == 0{
                    host_pixel.r&embed_pixel
                }else {
                    host_pixel.r|embed_pixel
                };
                let result_g = if embed_pixel == 0{
                    host_pixel.g&embed_pixel
                }else {
                    host_pixel.g|embed_pixel
                };
                let result_b = if embed_pixel == 0{
                    host_pixel.b&embed_pixel
                }else {
                    host_pixel.b|embed_pixel
                };

                image.put_pixel(x, y, Rgba([result_r,result_g,result_b,host_pixel.a]));

            }
        }

        image.save("embed_image.jpg").expect("Failed to save image");




        Ok(BaseImage::new("".to_string(), image::ColorType::L16))
    }
}

