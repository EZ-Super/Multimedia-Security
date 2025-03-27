use image::{ImageBuffer, Rgba};

use rand::Rng;

use crate::modules::get_image::BaseImage;
use crate::modules::point::Point;

#[derive(Clone)]
pub struct Watermark{
    watermark:BaseImage,
    width:u32,
    height:u32,
}

pub struct RandomNumber{
    scale:f32,
}

impl RandomNumber{
    pub fn new(scale:f32)->RandomNumber{
        RandomNumber{
            scale,
        }
    }
}

pub struct HostImage{
    pub base_image:BaseImage,
}

impl Watermark{
    pub fn new(base_image:&BaseImage)->Watermark{
        let clone_base_image = base_image.clone();
        let width = clone_base_image.width;
        let height = clone_base_image.height;
        Watermark{
            watermark:clone_base_image,
            width,
            height,
        }
    }
    pub fn post_host_image(&self){
        self.watermark.dyamic_image.save("result/host_image.png").expect("Failed to save image");
    }
}


impl HostImage{
    pub fn new(base_image:BaseImage)->HostImage{
        HostImage{
            base_image,
        }
    }
    pub fn post_host_image(&self){
        self.base_image.dyamic_image.save("result/host_image.png").expect("Failed to save image");
    }
    pub fn embed_image(&self,watermark: Watermark,watermark_x_number:u32,watermark_y_number:u32,embed_bit:u8)->Result<String,String>{
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

                let watermark_pixel = match watermark.watermark.pixel.get(&Point{x:x_pixel,y:y_pixel}){
                    Some(pixel) => pixel,
                    None => return Err("Watermark pixel not found".to_string())
                };

                let host_pixel = match host_image.pixel.get(&Point{x,y}){
                    Some(pixel) => pixel,
                    None => return Err("Host pixel not found".to_string())
                };
                let mut embed_pixel = 0;
                if watermark_pixel.r == 255 && watermark_pixel.g == 255 && watermark_pixel.b == 255{
                    embed_pixel = 1 ; //00000001
                }
                embed_pixel =  embed_pixel <<embed_bit;

                let result_r = set_bit(host_pixel.r, embed_bit, embed_pixel);
                let result_g = set_bit(host_pixel.g, embed_bit, embed_pixel);
                let result_b = set_bit(host_pixel.b, embed_bit, embed_pixel);

                image.put_pixel(x, y, Rgba([result_r,result_g,result_b,host_pixel.a]));
            }
        }
        let file_name = format!("result/embed_image {} x {} ({}).png",watermark_x_number,watermark_y_number,embed_bit);
        image.save(file_name.clone()).expect("Failed to save image");
        
        Ok(file_name)
    }

    pub fn embed_image_with_random_number(&self,random_number: RandomNumber,embed_bit:u32)->Result<String,String>{
        let host_image = self.base_image.clone();
        let image = host_image.dyamic_image.to_rgba8();

        let new_image_width = if random_number.scale > 1.0 { (host_image.width as f32 * random_number.scale) as u32} else {host_image.width as u32};
        let new_image_height = if random_number.scale > 1.0 {(host_image.height as f32 * random_number.scale) as u32} else {host_image.height as u32};

        let mut new_image = ImageBuffer::from_pixel(new_image_width, new_image_height, Rgba([0,0,0,255]));

        for x in 0..new_image_width{
            for y in 0..new_image_height{
                let point_x = x.clone() as u32;
                let point_y = y.clone() as u32;

                let host_pixel =  match image.get_pixel_checked(point_x, point_y){
                    Some(pixel) => pixel,
                    None => &Rgba([0,0,0,255])
                };

                let mut rng = rand::rng();
                let mut embed_pixel = rng.random_range(0..2);
                let [r,g,b,a] = host_pixel.0;
                embed_pixel =  embed_pixel <<embed_bit;

                let result_r = set_bit(r, embed_bit as u8, embed_pixel as u8);
                let result_g = set_bit(g, embed_bit as u8, embed_pixel as u8);
                let result_b = set_bit(b, embed_bit as u8, embed_pixel as u8);

                new_image.put_pixel(point_x, point_y, Rgba([result_r,result_g,result_b,a]));
            }
        }
        let file_name = format!("result/embed_image with random number {}x{} ({}).png",new_image_height,new_image_width,random_number.scale);
        new_image.save(file_name.clone()).expect("Failed to save image");
        
        Ok(file_name)

    }
}

fn set_bit(n: u8, pos: u8,embed_pixel:u8) -> u8 {
    if embed_pixel == 0{
        return n & !(1 << pos)
    }else{
        return n | (1 << pos)
    }
}