-- 1. FEFCO / ECMA Standard Catalog
CREATE TABLE IF NOT EXISTS box_styles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT UNIQUE NOT NULL,                  -- e.g., 'FEFCO 0201'
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    description TEXT NOT NULL,
    folding_allowance_mm REAL DEFAULT 3.0,
    joint_flap_mm REAL DEFAULT 30.0,
    is_active INTEGER DEFAULT 1
);

-- 2. Cardboard Board / Material Specifications
CREATE TABLE IF NOT EXISTS board_materials (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    material_code TEXT UNIQUE NOT NULL,         -- e.g., 'B-FLUTE-STANDARD'
    flute_type TEXT NOT NULL,
    thickness_mm REAL NOT NULL,
    ect_rating TEXT,
    cost_per_sq_m REAL NOT NULL,
    max_load_kg REAL NOT NULL
);

-- 3. Industry Clearance & Protection Rules
CREATE TABLE IF NOT EXISTS category_protection_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_name TEXT UNIQUE NOT NULL,         -- e.g., 'Glassware & Ceramics'
    min_clearance_mm REAL NOT NULL,
    recommended_padding TEXT NOT NULL,
    fragility_factor INTEGER CHECK (fragility_factor BETWEEN 1 AND 5)
);

-- 4. User Evaluation Sessions
CREATE TABLE IF NOT EXISTS packaging_evaluations (
    id TEXT PRIMARY KEY,                        -- Generated in Go application using google/uuid
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    product_name TEXT NOT NULL,
    dimensions_lwh_mm TEXT NOT NULL,            -- Stored as JSON string
    weight_kg REAL NOT NULL,
    category TEXT NOT NULL,
    selected_fefco_code TEXT REFERENCES box_styles(code),
    calculated_dim_weight REAL NOT NULL,
    agent_reasoning TEXT
);