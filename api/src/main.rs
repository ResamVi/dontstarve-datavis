#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate rocket;
#[macro_use] extern crate rocket_contrib;

use rocket_contrib::json::Json;
use rocket_contrib::databases::postgres;

use chrono::{Utc, Duration};

#[database("db")]
struct Db(postgres::Connection);

// -- Data model --

type Item = (String, i64); // e.g. ["Wigfrid", 15] or ["Russia", 1243]

// -- Handlers --

#[get("/characters?<modded>")]
fn characters(conn: Db, modded: bool) -> Json<Vec<Item>> {

    let query_string = if modded {
        "SELECT character, count \
        FROM count_character \
        LIMIT 20"
    } else {
        "SELECT character, count \
        FROM count_character \
        WHERE character \
        IN ('wendy', 'wathgrithr', 'wilson', 'woodie', 'wolfgang', 'wickerbottom', 'wx78', 'walter', 'webber', 'winona', 'waxwell', 'wortox', 'wormwood', 'wurt', 'wes', 'willow', 'warly')"
    };

    let mut characters: Vec<Item> = vec![];
    for row in &conn.0.query(query_string, &[]).unwrap() {
        characters.push((rename_char(row.get(0)), row.get(1)));
    }

    Json(characters)
}

#[get("/characters/<country>")]
fn characters_by_country(conn: Db, country: String) -> Json<Vec<Item>> {

    let query_string = "
        SELECT character, count \
        FROM count_character_by_country \
        WHERE country = $1 \
        AND character \
        IN ('wendy', 'wathgrithr', 'wilson', 'woodie', 'wolfgang', 'wickerbottom', 'wx78', 'walter', 'webber', 'winona', 'waxwell', 'wortox', 'wormwood', 'wurt', 'wes', 'willow', 'warly');";

    let mut characters: Vec<Item> = vec![];
    for row in &conn.0.query(query_string, &[&title_case(&country)]).unwrap() {
        characters.push((rename_char(row.get(0)), row.get(1)));
    }

    Json(characters)
}

#[get("/meta/age")]
fn age(conn: Db) -> Json<i64> {

    let rows = conn.0.query("SELECT date FROM last_update", &[]).unwrap();
    
    let last_update = rows.get(0).get(0);
    let now         = Utc::now().naive_utc();
    let age         = Duration::num_minutes(&now.signed_duration_since(last_update));

    Json(age)
}

#[get("/meta/countries")]
fn countries(conn: Db) -> Json<Vec<String>> {

    let mut countries: Vec<String> = vec![];
    for row in &conn.0.query("SELECT country FROM count_player", &[]).unwrap() {
        countries.push(row.get(0));
    }

    Json(countries)
}

#[get("/meta/<ent_type>")]
fn volume(conn: Db, ent_type: String) -> Json<i64> {

    let query_string = match ent_type.as_str() {
        "players"   => "SELECT player_count FROM count",
        "servers"   => "SELECT server_count FROM count",
        _           => panic!("ent_type not found"),
    };
    
    let rows = conn.0.query(query_string, &[]).unwrap();
    let volume = rows.get(0).get(0);

    Json(volume)
}

#[get("/count/<ent_type>")]
fn count(conn: Db, ent_type: String) -> Json<Vec<Item>> {

    let query_string = match ent_type.as_str() {
        "allplayers"    => "SELECT country, count FROM count_player",
        "players"       => "SELECT country, count FROM count_player LIMIT 20",
        "servers"       => "SELECT country, count FROM count_server LIMIT 20",
        "platforms"     => "SELECT platform, count FROM count_platform LIMIT 4",
        "intent"        => "SELECT intent, count FROM count_intent LIMIT 4",
        "modded"        => "SELECT mods, count FROM count_vanilla LIMIT 2",
        "season"        => "SELECT season, count FROM count_season LIMIT 4",
        _               => panic!("ent_type not found"),
    };

    let mut counts: Vec<Item> = vec![];
    for row in &conn.0.query(query_string, &[]).unwrap() {
        counts.push((row.get(0), row.get(1)));
    }

    Json(counts)
}

fn main() {
    let cors = rocket_cors::CorsOptions::default().to_cors().unwrap();

    rocket::ignite()
        .mount("/", routes![characters, characters_by_country, age, countries, count, volume])
        .attach(Db::fairing())
        .attach(cors)
        .launch();
}

// -- Don't Starve specific domain logic --

fn rename_char(name: String) -> String {
    match name.as_str() {
        ""              => "<Selecting>".to_string(),
        "wathgrithr"    => "Wigfrid".to_string(),
        "waxwell"       => "Maxwell".to_string(),
        "monkey_king"   => "Wilbur".to_string(),
        "wx78"          => "WX-78".to_string(),
        _               => capitalize(name),
    }
}

fn capitalize(word: String) -> String {
    word.chars().take(1).flat_map(char::to_uppercase).chain(word.chars().skip(1)).collect::<String>()
}

fn title_case(words: &str) -> String {
    words.split(" ").map(|word| capitalize(word.to_string())).collect::<Vec<String>>().join(" ")
}