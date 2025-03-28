use image::Rgba;
use log::info;


use crate::modules::get_image::BaseImage;
use crate::modules::point::Point;

pub struct PSNR{
    pub psnr: f64,
    pub mse: f64,
    pub max: f64,
    pub host_image:BaseImage,
    pub embed_image:BaseImage,
}

impl PSNR{
    pub fn new(embed_image:BaseImage,host_image:BaseImage)->PSNR{
        PSNR{
            psnr:0.0,
            mse:0.0,
            max:255.0,
            host_image,
            embed_image,
        }
    }
    pub fn calculate_psnr(&mut self)->&mut Self{
        let host_pixel = self.host_image.clone();
        let water_pixel = self.embed_image.clone();

        let hegiht = if host_pixel.height > water_pixel.height{host_pixel.height} else {water_pixel.height};
        let width = if host_pixel.width > water_pixel.width{host_pixel.width} else {water_pixel.width};
        //println!("height:{} width:{}",hegiht,width);
        for x in 0..hegiht{
            for y in 0..width{
                let host_pixel = host_pixel.pixel.get(&Point { x, y});
                let embedded_pixel = water_pixel.pixel.get(&Point { x, y});
            
                let [hr,hg,hb,_] = match host_pixel{
                    Some(pixel) => pixel.get_rgb(),
                    None => [0,0,0,255],
                };
                let [er,eg,eb,_] = match embedded_pixel{
                    Some(pixel) => pixel.get_rgb(),
                    None => [0,0,0,255],
                };

                let diff_r = (hr as i32 - er as i32).pow(2);
                let diff_g = (hg as i32 - eg as i32).pow(2);
                let diff_b = (hb as i32 - eb as i32).pow(2);
                //info!("diff_r:{} diff_g:{} diff_b:{}",diff_r,diff_g,diff_b);
                self.mse += diff_r as f64;
                self.mse += diff_g as f64;
                self.mse += diff_b as f64; 
            }
        }      
        let height_width = hegiht * width;
        self.mse = self.mse / (height_width as f64);
        self.psnr = 10.0 * (self.max.powi(2) / self.mse).log10();

        self
    }

    pub fn get_psnr(&self)->f64{
        self.psnr
    }
}