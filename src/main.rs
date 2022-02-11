use clap::Parser;
use std::process;
use std::io::{stdin, Write};
use std::path::Path;
use std::fs::File;

#[derive(Parser, Debug)]
#[clap(version = "1.0.0")]
#[clap(author = "Jason Hilder. <jhilder95@gmail.com>")]
#[clap(about = "A dynamic wallpaper daemon.", long_about = None)]
struct Args {
    #[clap(short, long)]
    name: String,

    #[clap(short, long, default_value_t = 1)]
    count: u8,

    #[clap(short, long)]
    test: Option<String>,
}

const DEFAULT_PATH: &str = "./config-walter";

fn main() {
    // Function to check if config file exsists. If not prompt user to use -i or --install
    init_process();

    let args: Args = Args::parse();

    for _ in 0..args.count {
        println!("Hello {}!", args.name)
    }
}

fn init_process() {
    let config_path = Path::new(DEFAULT_PATH).exists();

    if config_path {
        // Read contents of config_path into a Config Struct
        // Config Struct needs to stay alive for the entire app process `static
        println!("Continue");
    } else {
        let mut user_input = String::new();
        println!("Config file does not exist, do you want to run the setup?(yes/no)");

        match stdin().read_line(&mut user_input) {
            Ok(_) => match user_input.as_str().trim() {
                "yes" => create_config_file(),
                "no" => process::exit(0x0100),
                _ => println!("Invalid answer please type yes or no."),
            },
            Err(error) => println!("error: {}", error),
        }
    }
}

#[derive(Debug)]
struct InputError;

enum ChangeType {
    Weather,
    Time
}

impl ChangeType {
    fn from_str(s: &str) -> Result<ChangeType, InputError> {
        match s {
            "weather" | "Weather" => Ok(ChangeType::Weather),
            "time" | "Time" => Ok(ChangeType::Time),
            _ => Err(InputError)
        }
    }
}

struct Config {
    change_type: ChangeType,
    image_path: Path
}

impl Config {
    // TODO: Create constructor asociated function
    fn new() {
        unimplemented!();
    }
}

// TODO:
// Prompt user questions, with answers build
// a Config_File struct, use the struct to
// write to file at path json it for easy reading/writing
fn create_config_file() {
    // TODO: Get os agnostic path to pictures/walter/
    let default_image_path = Path::new("~/Pictures/walter");
    let mut config_file = File::create("./config_file").unwrap();
    let mut user_choice = String::new();
    let mut user_path = String::new();

    println!("Wallpaper to change dependant on?(weather/time)");
    stdin().read_line(&mut user_choice).expect("Invalid input");
    //if let Ok?

    println!("Default program wallpaper path is set to {:?}.", default_image_path);
    println!("Press enter to keep default path else enter valid path:");
    stdin().read_line(&mut user_path).expect("Invalid input");

    if user_path.trim().len() > 0 {
        // With user_path input try vaild it and create a new Path
    }

    config_file.write_all(user_choice.as_bytes()).unwrap()
}
