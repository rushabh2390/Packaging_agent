-- Enable UUID extension for evaluation IDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. FEFCO / ECMA Standard Catalog
CREATE TABLE IF NOT EXISTS box_styles (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,            -- e.g., 'FEFCO 0201'
    name VARCHAR(120) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    folding_allowance_mm FLOAT DEFAULT 3.0,
    joint_flap_mm FLOAT DEFAULT 30.0,
    is_active BOOLEAN DEFAULT TRUE
);

-- 2. Cardboard Board / Material Specifications
CREATE TABLE IF NOT EXISTS board_materials (
    id SERIAL PRIMARY KEY,
    material_code VARCHAR(30) UNIQUE NOT NULL,   -- e.g., 'B-FLUTE-STANDARD'
    flute_type VARCHAR(10) NOT NULL,
    thickness_mm FLOAT NOT NULL,
    ect_rating VARCHAR(20),
    cost_per_sq_m NUMERIC(10, 4) NOT NULL,
    max_load_kg FLOAT NOT NULL
);

-- 3. Industry Clearance & Protection Rules
CREATE TABLE IF NOT EXISTS category_protection_rules (
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(60) UNIQUE NOT NULL,   -- e.g., 'Glassware & Ceramics'
    min_clearance_mm FLOAT NOT NULL,
    recommended_padding VARCHAR(100) NOT NULL,
    fragility_factor INT CHECK (fragility_factor BETWEEN 1 AND 5)
);

-- 4. User Evaluation Sessions
CREATE TABLE IF NOT EXISTS packaging_evaluations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIMEZONE DEFAULT CURRENT_TIMESTAMP,
    product_name VARCHAR(100) NOT NULL,
    dimensions_lwh_mm JSONB NOT NULL,
    weight_kg FLOAT NOT NULL,
    category VARCHAR(60) NOT NULL,
    selected_fefco_code VARCHAR(20) REFERENCES box_styles(code),
    calculated_dim_weight FLOAT NOT NULL,
    agent_reasoning TEXT
);