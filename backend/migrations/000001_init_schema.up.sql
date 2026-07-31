-- Initial schema for PW Toolkit

CREATE TABLE IF NOT EXISTS classes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(50) UNIQUE NOT NULL,
    name_ru VARCHAR(100) NOT NULL,
    name_en VARCHAR(100),
    base_stats JSON NOT NULL,
    stat_growth_per_level JSON NOT NULL,
    allowed_equipment_types JSON,
    icon_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS items (
    id INT PRIMARY KEY AUTO_INCREMENT,
    game_id INT UNIQUE NOT NULL,
    name_ru VARCHAR(200) NOT NULL,
    name_en VARCHAR(200),
    type VARCHAR(50) NOT NULL,
    subtype VARCHAR(50),
    level_requirement INT DEFAULT 1,
    class_restriction JSON,
    base_stats JSON NOT NULL,
    gem_slots INT DEFAULT 0,
    set_id INT,
    icon_url VARCHAR(500),
    rarity VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_type (type),
    INDEX idx_level (level_requirement)
);

CREATE TABLE IF NOT EXISTS gems (
    id INT PRIMARY KEY AUTO_INCREMENT,
    game_id INT UNIQUE NOT NULL,
    name_ru VARCHAR(200) NOT NULL,
    name_en VARCHAR(200),
    level INT NOT NULL,
    bonuses JSON NOT NULL,
    icon_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS builds (
    id VARCHAR(10) PRIMARY KEY,
    internal_id BIGINT AUTO_INCREMENT UNIQUE,
    name VARCHAR(200),
    class_id INT NOT NULL,
    level INT NOT NULL,
    equipment JSON NOT NULL,
    cards JSON,
    books JSON,
    genie_id INT,
    pangu_souls JSON,
    star_disks JSON,
    titles JSON,
    calculated_stats JSON,
    view_count INT DEFAULT 0,
    last_viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (class_id) REFERENCES classes(id),
    INDEX idx_class (class_id),
    INDEX idx_last_viewed (last_viewed_at)
);

-- Seed data: 10 PW classes (version Заговор судьбы)
INSERT INTO classes (code, name_ru, name_en, base_stats, stat_growth_per_level, allowed_equipment_types) VALUES
('warrior', 'Воин', 'Warrior',
 '{"hp": 500, "mp": 100, "attack": 20, "defense": 15}',
 '{"hp": 50, "mp": 5, "attack": 3, "defense": 2}',
 '["sword", "axe", "spear", "heavy_armor"]'),
('wizard', 'Маг', 'Wizard',
 '{"hp": 300, "mp": 400, "attack": 35, "defense": 8}',
 '{"hp": 20, "mp": 30, "attack": 5, "defense": 1}',
 '["staff", "magic_weapon", "light_armor", "cloth_armor"]'),
('cleric', 'Жрец', 'Cleric',
 '{"hp": 350, "mp": 350, "attack": 25, "defense": 10}',
 '{"hp": 25, "mp": 25, "attack": 3, "defense": 2}',
 '["hammer", "magic_weapon", "light_armor", "cloth_armor"]'),
('archer', 'Лучник', 'Archer',
 '{"hp": 400, "mp": 200, "attack": 30, "defense": 12}',
 '{"hp": 30, "mp": 10, "attack": 4, "defense": 1}',
 '["bow", "crossbow", "light_armor"]'),
('beast_slayer', 'Оборотень', 'Beast Slayer',
 '{"hp": 600, "mp": 80, "attack": 25, "defense": 20}',
 '{"hp": 60, "mp": 3, "attack": 3, "defense": 3}',
 '["axe", "spear", "heavy_armor"]'),
('venomancer', 'Друид', 'Venomancer',
 '{"hp": 380, "mp": 300, "attack": 20, "defense": 10}',
 '{"hp": 25, "mp": 20, "attack": 2, "defense": 2}',
 '["magic_weapon", "light_armor", "cloth_armor"]'),
('assassin', 'Убийца', 'Assassin',
 '{"hp": 420, "mp": 150, "attack": 40, "defense": 10}',
 '{"hp": 30, "mp": 8, "attack": 5, "defense": 1}',
 '["daggers", "light_armor"]'),
('psychic', 'Шаман', 'Psychic',
 '{"hp": 350, "mp": 380, "attack": 30, "defense": 9}',
 '{"hp": 22, "mp": 28, "attack": 4, "defense": 1}',
 '["orb", "magic_weapon", "cloth_armor"]'),
('seeker', 'Страж', 'Seeker',
 '{"hp": 400, "mp": 250, "attack": 32, "defense": 13}',
 '{"hp": 28, "mp": 15, "attack": 4, "defense": 2}',
 '["sword", "magic_weapon", "light_armor"]'),
('mystic', 'Мистик', 'Mystic',
 '{"hp": 360, "mp": 320, "attack": 28, "defense": 11}',
 '{"hp": 24, "mp": 22, "attack": 3, "defense": 2}',
 '["magic_weapon", "light_armor", "cloth_armor"]');
