use std::fs::OpenOptions;
use std::fs::File;
use std::io::Write;

pub struct Note{
    pub note_path: String,
    pub file: File,
}

impl Note{
    pub fn new(note_path:String)->Note{
        let file = OpenOptions::new()
            .write(true)
            .create(true)
            .open(&note_path)
            .unwrap();
        Note{
            note_path,
            file,
        }
    }
    pub fn write(&mut self,note:String){
        self.file.write_all(note.as_bytes()).unwrap();
    }
}