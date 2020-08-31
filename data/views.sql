-- Counts players and servers
CREATE OR REPLACE VIEW count AS
SELECT 
    (SELECT COUNT(*) FROM server) AS server_count, 
    (SELECT COUNT(*) FROM player) AS player_count;

-- Count by Server Country
CREATE OR REPLACE VIEW count_server AS
SELECT country, COUNT(country)
FROM server 
GROUP BY country 
ORDER BY COUNT(country) DESC;

-- Count by Player country
CREATE OR REPLACE VIEW count_player AS
SELECT country, COUNT(country)
FROM player 
GROUP BY country 
ORDER BY COUNT(country) DESC;

-- Count by Platform
CREATE OR REPLACE VIEW count_platform AS
SELECT platform, COUNT(platform)
FROM server
GROUP BY platform
ORDER BY COUNT(platform) DESC;

-- Last Update
CREATE OR REPLACE VIEW last_update AS
SELECT date 
FROM server 
ORDER BY date DESC 
LIMIT 1;

-- Count by Intent
CREATE OR REPLACE VIEW count_intent AS
SELECT intent, COUNT(intent)
FROM server
GROUP BY intent
HAVING COUNT(intent) > 10
ORDER BY COUNT(intent) DESC;

-- Vanilla vs Modded
CREATE OR REPLACE VIEW count_vanilla AS
SELECT mods, COUNT(*)
FROM server
GROUP BY mods;

-- Count by season
CREATE OR REPLACE VIEW count_season AS
SELECT season, COUNT(season)
FROM server
WHERE season <> ''
GROUP BY season
ORDER BY COUNT(season) DESC;

-- Count characters
CREATE OR REPLACE VIEW count_character AS
SELECT character, COUNT(character)
FROM player
GROUP BY character
ORDER BY COUNT(character) DESC;

-- Count pairs of (Character, Country) 
CREATE OR REPLACE VIEW count_character_by_country AS
SELECT character, country, iso, COUNT(character)
FROM player
GROUP BY character, country, iso
ORDER BY COUNT(character) DESC;

-- Calculate % of character usage per country (need to cast to float because numeric type in rust not supported (fuck sake man...))
CREATE OR REPLACE VIEW percentage_character_by_country AS
SELECT character, c.country, c.iso, ROUND((c.count/p.count::DECIMAL)*100, 2)::float AS percent, c.count AS char_count, p.count AS total_count
FROM count_character_by_country c
INNER JOIN count_player p
ON c.country = p.country;