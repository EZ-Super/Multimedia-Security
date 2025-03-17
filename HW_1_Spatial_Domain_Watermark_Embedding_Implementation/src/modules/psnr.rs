use crate::modules::note;

struct PSNR{
    pub psnr: f64,
    pub mse: f64,
    pub max: f64,
    pub note : note::Note,
}

impl PSNR{
    pub fn new(note:note::Note)->PSNR{
        PSNR{
            psnr:0.0,
            mse:0.0,
            max:255.0,
            note,
        }
    }
    pub fn calculate_psnr(&mut self,host_image:&image::DynamicImage,watermark_image:&image::DynamicImage){
        let host_image = host_image.to_rgba8();
        let watermark_image = watermark_image.to_rgba8();
        let host_image = host_image.to_vec();
        let watermark_image = watermark_image.to_vec();
        let mut mse = 0.0;
        for i in 0..host_image.len(){
            let host_pixel = host_image[i] as f64;
            let watermark_pixel = watermark_image[i] as f64;
            mse += (host_pixel - watermark_pixel).powi(2);
        }
        mse = mse / host_image.len() as f64;
        self.mse = mse;
        self.psnr = 10.0 * (self.max.powi(2) / mse).log10();
        self.note.write(format!("PSNR:{}",self.psnr));
    }
}