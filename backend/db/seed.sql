-- 1. FEFCO & ECMA Structural Catalog
INSERT INTO box_styles (code, name, category, description, folding_allowance_mm, joint_flap_mm) VALUES
('FEFCO 0201', 'Regular Slotted Container (RSC)', 'Shipping', 'Standard shipping box with 4 meeting flaps top and bottom.', 3.0, 35.0),
('FEFCO 0427', 'Roll-End Tuck Top Mailer', 'E-commerce', 'Self-locking mailer box with double-wall sides. No tape required.', 1.5, 15.0),
('FEFCO 0300', 'Telescope Box (Lid & Tray)', 'Display', 'Two-piece box where the lid slides over the bottom tray.', 2.0, 20.0),
('FEFCO 0211', 'Snap Lock Bottom Box', 'Retail', 'Slotted box with self-locking bottom flaps.', 3.0, 30.0),
('FEFCO 0411', 'Wrap Around Blank', 'Automated', 'Flat folder sheet folded tightly around flat items.', 1.5, 15.0)
ON CONFLICT (code) DO NOTHING;

-- 2. Cardboard Flute Profiles & Material Costs
INSERT INTO board_materials (material_code, flute_type, thickness_mm, ect_rating, cost_per_sq_m, max_load_kg) VALUES
('E-FLUTE-FINE', 'E', 1.5, 'ECT 23', 0.3500, 5.0),
('B-FLUTE-STANDARD', 'B', 3.0, 'ECT 32 (200#)', 0.4800, 15.0),
('C-FLUTE-SHIPPING', 'C', 4.0, 'ECT 32', 0.5200, 20.0),
('BC-DOUBLE-WALL', 'BC', 6.5, 'ECT 48 (275#)', 0.9200, 40.0)
ON CONFLICT (material_code) DO NOTHING;

-- 3. Category Protection & Clearance Allowances
INSERT INTO category_protection_rules (category_name, min_clearance_mm, recommended_padding, fragility_factor) VALUES
('Glassware & Ceramics', 50.0, 'Bubble Wrap / Molded Pulp Corners', 5),
('Consumer Electronics', 35.0, 'EPE Foam Inserts / Anti-Static Wrap', 4),
('Bottles & Liquids', 40.0, 'Corrugated Divider Grids', 4),
('Books & Documents', 10.0, 'Tight Wrap / Kraft Paper Fill', 2),
('Apparel & Textiles', 5.0, 'Polybag / Tissue Paper Wrap', 1),
('Heavy Machinery / Hardware', 25.0, 'Double-wall Cushion Pad Base', 3)
ON CONFLICT (category_name) DO NOTHING;