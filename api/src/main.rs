#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate rocket;
#[macro_use] extern crate rocket_contrib;

use rocket_contrib::json::Json;
use rocket_contrib::databases::postgres;

use chrono::{Utc, Duration};

#[database("shortterm")]
struct ShortTerm(postgres::Connection);

#[database("longterm")]
struct LongTerm(postgres::Connection);

// -- Data model --

type Item = (String, i64); // e.g. ["Wigfrid", 15] or ["Russia", 1243]

type IsoItem = (String, String, f64); // e.g. ["Greece", "GR", 42.5], ["Poland", "PL", 25.5]

type Percentage = (String, f64); // e.g. ["Greece", 42.5], ["Poland", 25.5]

#[derive(serde::Serialize)]
struct Series<T> {
    name: String,
    data: Vec<T>,
}

// -- Handlers --
#[get("/series/continents")]
fn series_continents(conn: LongTerm) -> Json<Vec<Series<Item>>> {

    let query_string = "
        SELECT *
        FROM series_continent
        WHERE date BETWEEN NOW() - INTERVAL '5 DAYS' AND NOW()
        ORDER BY date DESC";

    const CONTINENTS: [&str; 6] = ["Asia", "Europe", "North America", "South America", "Africa", "Oceania"];

    // Prepare result object
    let mut result: Vec<Series<Item>> = vec![];
    for continent in CONTINENTS.iter() {
        result.push(Series{name: continent.to_string(), data: Vec::new()})
    }

    for row in &conn.0.query(query_string, &[]).unwrap() {
        
        let date: chrono::NaiveDateTime = row.get(1);

        for (i, _) in CONTINENTS.iter().enumerate() {
            let count: i32 = row.get(2 + i);
            result[i].data.push((date.format("%Y-%m-%dT%H:%M:%S%.f").to_string(), count as i64));
        }
    }

    Json(result)
}

#[get("/series/characters")]
fn series_characters(conn: LongTerm) -> Json<Vec<Series<Item>>> {
    
    let query_string = "
        SELECT *
        FROM series_character_count
        WHERE date BETWEEN NOW() - INTERVAL '5 DAYS' AND NOW()
        ORDER BY date DESC";

    const CHARACTERS: [&str; 17] = [
        "Wilson", "Willow", "Wolfgang",
        "Wendy", "WX-78", "Wickerbottom",
        "Woodie", "Wes", "Maxwell",
        "Wigfrid", "Webber", "Warly",
        "Wormwood", "Winona", "Wortox", "Wurt",
        "Walter"];

    // Prepare result object
    let mut result: Vec<Series<Item>> = vec![];
    for character in CHARACTERS.iter() {
        result.push(Series{name: character.to_string(), data: Vec::new()})
    }

    for row in &conn.0.query(query_string, &[]).unwrap() {
        let date: chrono::NaiveDateTime = row.get(1);

        for (i, _) in CHARACTERS.iter().enumerate() {
            let count: i32 = row.get(2 + i);
            result[i].data.push((date.format("%Y-%m-%dT%H:%M:%S%.f").to_string(), count as i64));
        }
    }
    
    Json(result)
}

#[get("/series/preferences/<character>")]
fn series_preferences_by_characters(conn: LongTerm, character: String) -> Json<Vec<Series<Percentage>>> {
    
    let query_string = "
        SELECT * FROM
        (
            SELECT *, ROW_NUMBER() OVER(ORDER BY date DESC) AS rnk
            FROM series_character_ranking
            WHERE character = $1
        ) AS t
        WHERE rnk % 3 = 0
        LIMIT 500"; // take every 3rd entry
    
    let mut result: Vec<Series<Percentage>> = vec![];
    for row in &conn.0.query(query_string, &[&character]).unwrap() {
        
        let date: chrono::NaiveDateTime = row.get(1);
        
        let mut entry = Series{name: date.format("%Y-%m-%dT%H:%M:%S%.f").to_string(), data: Vec::new()};
        
        // id | date |character | first | first_percent | second | second_percent | ... | fifth
        for i in [3, 5, 7, 9, 11].iter() {
            entry.data.push((row.get(i), row.get(i+1)))
        }
        
        result.push(entry);
    }
    
    Json(result)
}

#[get("/player/<attribute>/<name>")]
fn player(conn: LongTerm, attribute: String, name: String) -> Json<Vec<Percentage>> {

    if attribute != "country" && attribute != "character" {
        return Json(vec![]);
    }

    let query_string = format!("
            SELECT
            {},
            ROUND(
                EXTRACT(epoch from SUM(duration))::DECIMAL / (
                    SELECT EXTRACT(epoch from SUM(duration))
                    FROM player
                    WHERE name = $1
                )::DECIMAL * 100, 2
            )::float AS percentage
        FROM player
        WHERE name = $1
        GROUP BY {}
        ORDER BY percentage DESC", attribute, attribute);

        let mut percentages: Vec<Percentage> = vec![];
        for row in &conn.0.query(query_string.as_str(), &[&name]).unwrap() {
            match attribute.as_str() {
                "character" => percentages.push( (rename_char(row.get(0)), row.get(1)) ),
                "country" => percentages.push((row.get(0), row.get(1))),
                _ => panic!("should not happen")
                
            } 
        }
    
        Json(percentages)
}

#[get("/characters?<modded>")]
fn characters(conn: ShortTerm, modded: bool) -> Json<Vec<Item>> {

    let query_string = if modded {
        "SELECT character, count
        FROM count_character
        LIMIT 30"
    } else {
        "SELECT character, count
        FROM count_character
        WHERE character
        IN ('wendy', 'wathgrithr', 'wilson', 'woodie', 'wolfgang', 'wickerbottom', 'wx78', 'walter', 'webber', 'winona', 'waxwell', 'wortox', 'wormwood', 'wurt', 'wes', 'willow', 'warly')"
    };

    let mut characters: Vec<Item> = vec![];
    for row in &conn.0.query(query_string, &[]).unwrap() {
        characters.push((rename_char(row.get(0)), row.get(1)));
    }

    Json(characters)
}

#[get("/characters/<country>")]
fn characters_by_country(conn: ShortTerm, country: String) -> Json<Vec<Item>> {

    let query_string = "
        SELECT character, count
        FROM count_character_by_country
        WHERE country = $1
        AND character
        IN ('wendy', 'wathgrithr', 'wilson', 'woodie', 'wolfgang', 'wickerbottom', 'wx78', 'walter', 'webber', 'winona', 'waxwell', 'wortox', 'wormwood', 'wurt', 'wes', 'willow', 'warly')";    

    let mut characters: Vec<Item> = vec![];
    for row in &conn.0.query(query_string, &[&title_case(&country)]).unwrap() {
        characters.push((rename_char(row.get(0)), row.get(1)));
    }

    Json(characters)
}

#[get("/characters/percentage/<character>")]
fn country_percentage_by_character(conn:ShortTerm, character: String) -> Json<Vec<IsoItem>> {
    let query_string = "
        SELECT character, country, iso, percent
        FROM percentage_character_by_country
        WHERE character = $1
        AND total_count > 30
        ORDER BY percent DESC
        LIMIT 5";

    let mut countries: Vec<IsoItem> = vec![];
    for row in &conn.0.query(query_string, &[&character]).unwrap() {        
        countries.push((row.get(1), row.get(2), row.get(3)));
    }

    Json(countries)
}

#[get("/characters/country/<country>")]
fn country_percentage_by_country(conn:ShortTerm, country: String) -> Json<Vec<Percentage>> {
    let query_string = "
        SELECT character, percent
        FROM percentage_character_by_country
        WHERE country = $1
        AND character
        IN ('wendy', 'wathgrithr', 'wilson', 'woodie', 'wolfgang', 'wickerbottom', 'wx78', 'walter', 'webber', 'winona', 'waxwell', 'wortox', 'wormwood', 'wurt', 'wes', 'willow', 'warly')";

        let mut countries: Vec<Percentage> = vec![];
        for row in &conn.0.query(query_string, &[&title_case(&country)]).unwrap() {
            countries.push((rename_char(row.get(0)), row.get(1)));
        }

        Json(countries)
}

#[get("/meta/age")]
fn age(conn: ShortTerm) -> Json<i64> {

    let rows = conn.0.query("SELECT date FROM last_update", &[]).unwrap();
    
    let last_update = rows.get(0).get(0);
    let now         = Utc::now().naive_utc();
    let age         = Duration::num_minutes(&now.signed_duration_since(last_update));

    Json(age)
}

#[get("/meta/countries")]
fn countries(conn: ShortTerm) -> Json<Vec<String>> {

    let mut countries: Vec<String> = vec![];
    for row in &conn.0.query("SELECT country FROM count_player", &[]).unwrap() {
        countries.push(row.get(0));
    }

    Json(countries)
}

#[get("/meta/<ent_type>")]
fn volume(conn: ShortTerm, ent_type: String) -> Json<i64> {

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
fn count(conn: ShortTerm, ent_type: String) -> Json<Vec<Item>> {

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
        .mount("/", routes![
            series_characters, series_preferences_by_characters, series_continents,
            country_percentage_by_character, country_percentage_by_country,
            characters_by_country, player, characters, countries, volume, count, age])
        .attach(ShortTerm::fairing())
        .attach(LongTerm::fairing())
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