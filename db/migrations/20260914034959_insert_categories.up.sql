-- Seed data: 10 parent categories, each with 10 child categories (110 rows total)
-- Domain: e-commerce marketplace categories

-- ============================================================
-- PARENT CATEGORIES
-- ============================================================
INSERT INTO categories (id, name, parent_id) VALUES
('cat-electronics',    'Electronics',            NULL),
('cat-fashion',        'Fashion & Apparel',      NULL),
('cat-home',           'Home & Living',          NULL),
('cat-beauty',         'Beauty & Personal Care', NULL),
('cat-sports',         'Sports & Outdoors',      NULL),
('cat-toys',           'Toys, Kids & Baby',      NULL),
('cat-automotive',     'Automotive',             NULL),
('cat-books',          'Books & Stationery',     NULL),
('cat-grocery',        'Grocery & Gourmet',      NULL),
('cat-health',         'Health & Wellness',      NULL);

-- ============================================================
-- CHILD CATEGORIES
-- ============================================================

-- Electronics
INSERT INTO categories (id, name, parent_id) VALUES
('cat-electronics-smartphones', 'Smartphones',        'cat-electronics'),
('cat-electronics-laptops',     'Laptops',            'cat-electronics'),
('cat-electronics-tablets',     'Tablets',            'cat-electronics'),
('cat-electronics-cameras',     'Cameras',            'cat-electronics'),
('cat-electronics-headphones',  'Headphones & Audio', 'cat-electronics'),
('cat-electronics-tv',          'Televisions',        'cat-electronics'),
('cat-electronics-console',     'Gaming Consoles',    'cat-electronics'),
('cat-electronics-wearables',   'Smartwatches',       'cat-electronics'),
('cat-electronics-printers',    'Printers & Scanners','cat-electronics'),
('cat-electronics-storage',     'Storage & Drives',   'cat-electronics');

-- Fashion & Apparel
INSERT INTO categories (id, name, parent_id) VALUES
('cat-fashion-mens',       'Men''s Clothing',      'cat-fashion'),
('cat-fashion-womens',     'Women''s Clothing',    'cat-fashion'),
('cat-fashion-shoes',      'Shoes',                'cat-fashion'),
('cat-fashion-bags',       'Bags & Wallets',       'cat-fashion'),
('cat-fashion-watches',    'Watches',              'cat-fashion'),
('cat-fashion-jewelry',    'Jewelry',              'cat-fashion'),
('cat-fashion-sunglasses', 'Eyewear & Sunglasses', 'cat-fashion'),
('cat-fashion-hats',       'Hats & Caps',          'cat-fashion'),
('cat-fashion-belts',      'Belts',                'cat-fashion'),
('cat-fashion-kids',       'Kids'' Fashion',       'cat-fashion');

-- Home & Living
INSERT INTO categories (id, name, parent_id) VALUES
('cat-home-furniture',   'Furniture',           'cat-home'),
('cat-home-kitchen',     'Kitchen & Dining',    'cat-home'),
('cat-home-bedding',     'Bedding',             'cat-home'),
('cat-home-lighting',    'Lighting',            'cat-home'),
('cat-home-decor',       'Home Decor',          'cat-home'),
('cat-home-storage',     'Storage & Organization','cat-home'),
('cat-home-bath',        'Bath',                'cat-home'),
('cat-home-appliances',  'Home Appliances',     'cat-home'),
('cat-home-cleaning',    'Cleaning Supplies',   'cat-home'),
('cat-home-garden',      'Garden & Outdoor',    'cat-home');

-- Beauty & Personal Care
INSERT INTO categories (id, name, parent_id) VALUES
('cat-beauty-skincare',  'Skincare',            'cat-beauty'),
('cat-beauty-makeup',    'Makeup',              'cat-beauty'),
('cat-beauty-haircare',  'Hair Care',           'cat-beauty'),
('cat-beauty-fragrance', 'Fragrances',          'cat-beauty'),
('cat-beauty-nails',     'Nail Care',           'cat-beauty'),
('cat-beauty-mens',      'Men''s Grooming',     'cat-beauty'),
('cat-beauty-bath',      'Bath & Body',         'cat-beauty'),
('cat-beauty-tools',     'Beauty Tools',        'cat-beauty'),
('cat-beauty-oral',      'Oral Care',           'cat-beauty'),
('cat-beauty-suncare',   'Sun Care',            'cat-beauty');

-- Sports & Outdoors
INSERT INTO categories (id, name, parent_id) VALUES
('cat-sports-fitness',   'Fitness Equipment',   'cat-sports'),
('cat-sports-cycling',   'Cycling',             'cat-sports'),
('cat-sports-running',   'Running',             'cat-sports'),
('cat-sports-camping',   'Camping & Hiking',    'cat-sports'),
('cat-sports-fishing',   'Fishing',             'cat-sports'),
('cat-sports-football',  'Football',            'cat-sports'),
('cat-sports-basketball','Basketball',          'cat-sports'),
('cat-sports-swimming',  'Swimming',            'cat-sports'),
('cat-sports-yoga',      'Yoga & Pilates',      'cat-sports'),
('cat-sports-golf',      'Golf',                'cat-sports');

-- Toys, Kids & Baby
INSERT INTO categories (id, name, parent_id) VALUES
('cat-toys-action',      'Action Figures',      'cat-toys'),
('cat-toys-dolls',       'Dolls',               'cat-toys'),
('cat-toys-blocks',      'Building Blocks',     'cat-toys'),
('cat-toys-board',       'Board Games',         'cat-toys'),
('cat-toys-puzzles',     'Puzzles',             'cat-toys'),
('cat-toys-rc',          'Remote Control Toys', 'cat-toys'),
('cat-toys-educational', 'Educational Toys',    'cat-toys'),
('cat-toys-babygear',    'Baby Gear',           'cat-toys'),
('cat-toys-diapering',   'Diapering',           'cat-toys'),
('cat-toys-feeding',     'Baby Feeding',        'cat-toys');

-- Automotive
INSERT INTO categories (id, name, parent_id) VALUES
('cat-auto-tires',       'Tires & Wheels',      'cat-automotive'),
('cat-auto-oil',         'Oils & Fluids',       'cat-automotive'),
('cat-auto-battery',     'Batteries',           'cat-automotive'),
('cat-auto-interior',    'Interior Accessories','cat-automotive'),
('cat-auto-exterior',    'Exterior Accessories','cat-automotive'),
('cat-auto-electronics', 'Car Electronics',     'cat-automotive'),
('cat-auto-tools',       'Tools & Equipment',   'cat-automotive'),
('cat-auto-cleaning',    'Car Care & Cleaning', 'cat-automotive'),
('cat-auto-lighting',    'Car Lighting',        'cat-automotive'),
('cat-auto-motorcycle',  'Motorcycle Parts',    'cat-automotive');

-- Books & Stationery
INSERT INTO categories (id, name, parent_id) VALUES
('cat-books-fiction',    'Fiction',             'cat-books'),
('cat-books-nonfiction', 'Non-Fiction',         'cat-books'),
('cat-books-children',   'Children''s Books',   'cat-books'),
('cat-books-comics',     'Comics & Manga',      'cat-books'),
('cat-books-textbooks',  'Textbooks',           'cat-books'),
('cat-books-religion',   'Religion & Spirituality','cat-books'),
('cat-books-cooking',    'Cookbooks',           'cat-books'),
('cat-books-notebooks',  'Notebooks',           'cat-books'),
('cat-books-pens',       'Pens & Writing',      'cat-books'),
('cat-books-art',        'Art Supplies',        'cat-books');

-- Grocery & Gourmet
INSERT INTO categories (id, name, parent_id) VALUES
('cat-grocery-beverages', 'Beverages',          'cat-grocery'),
('cat-grocery-snacks',    'Snacks',             'cat-grocery'),
('cat-grocery-coffee',    'Coffee & Tea',       'cat-grocery'),
('cat-grocery-dairy',     'Dairy & Eggs',       'cat-grocery'),
('cat-grocery-bakery',    'Bakery',             'cat-grocery'),
('cat-grocery-pantry',    'Pantry Staples',     'cat-grocery'),
('cat-grocery-frozen',    'Frozen Food',        'cat-grocery'),
('cat-grocery-organic',   'Organic Food',       'cat-grocery'),
('cat-grocery-condiments','Sauces & Condiments','cat-grocery'),
('cat-grocery-candy',     'Candy & Chocolate',  'cat-grocery');

-- Health & Wellness
INSERT INTO categories (id, name, parent_id) VALUES
('cat-health-vitamins',   'Vitamins & Supplements','cat-health'),
('cat-health-firstaid',   'First Aid',          'cat-health'),
('cat-health-medical',    'Medical Devices',    'cat-health'),
('cat-health-nutrition',  'Sports Nutrition',   'cat-health'),
('cat-health-personal',   'Personal Hygiene',   'cat-health'),
('cat-health-eyecare',    'Eye Care',           'cat-health'),
('cat-health-mobility',   'Mobility Aids',      'cat-health'),
('cat-health-monitors',   'Health Monitors',    'cat-health'),
('cat-health-masks',      'Masks & Protection', 'cat-health'),
('cat-health-wellness',   'Wellness & Relaxation','cat-health');
