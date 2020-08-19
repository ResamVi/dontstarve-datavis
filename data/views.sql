-- Counts players and servers
CREATE OR REPLACE VIEW count AS
SELECT 
    (SELECT COUNT(*) FROM server) AS server_count, 
    (SELECT COUNT(*) FROM player) AS player_count;

-- Count by Server Origin
CREATE OR REPLACE VIEW count_server AS
SELECT origin, COUNT(origin)
FROM server 
GROUP BY origin 
ORDER BY COUNT(origin) DESC;

-- Count by Player Origin
CREATE OR REPLACE VIEW count_player AS
SELECT origin, COUNT(origin)
FROM player 
GROUP BY origin 
ORDER BY COUNT(origin) DESC;


-- Count by Platform
CREATE OR REPLACE VIEW count_platform AS
SELECT platform, COUNT(platform)
FROM server
GROUP BY platform
ORDER BY COUNT(platform) DESC;

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

-- Count pairs of (Character, Origin) 
CREATE OR REPLACE VIEW count_character_by_origin AS
SELECT character, origin, COUNT(character)
FROM player
GROUP BY character, origin
ORDER BY COUNT(character) DESC;

-- Calculate % of character usage per origin
CREATE OR REPLACE VIEW percentage_character_by_origin AS
SELECT character, c.origin, ROUND((c.count/p.count::DECIMAL)*100, 2) AS percent
FROM count_character_by_origin c
INNER JOIN count_player p
ON c.origin = p.origin;