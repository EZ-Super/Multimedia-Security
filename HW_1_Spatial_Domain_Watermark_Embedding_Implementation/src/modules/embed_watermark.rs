use crate::modules::get_image::BaseImage;
use crate::modules::point::RGPPixel;
use crate::modules::point::Point;

use std::collections::HashMap;

struct Watermark{
    width:u32,
    height:u32,
    pixel:HashMap<Point,RGPPixel>,
}

struct RandomNumber{
    width:u32,
    height:u32,
    embed_numer:Vec<i32>,
}
struct EmbedImage{
    base_image:BaseImage,
}

impl Watermark{
    fn new(base_image:BaseImage)->Watermark{
        let width = base_image.width;
        let height = base_image.height;
        let pixel = base_image.pixel;
        Watermark{
            width,
            height,
            pixel:pixel,
        }
    }

}


impl EmbedImage{

}