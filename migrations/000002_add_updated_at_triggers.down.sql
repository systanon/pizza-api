DROP TRIGGER IF EXISTS trg_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS trg_variants_updated_at ON variants;
DROP TRIGGER IF EXISTS trg_products_updated_at ON products;
DROP TRIGGER IF EXISTS trg_product_variant_prices_updated_at ON product_variant_prices;
DROP TRIGGER IF EXISTS trg_addons_updated_at ON addons;

DROP FUNCTION IF EXISTS set_updated_at();