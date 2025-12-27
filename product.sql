-- Creates the products table
CREATE TABLE products (
    -- BIGSERIAL คือ BIGINT ที่เพิ่มค่าอัตโนมัติ (auto-increment) เหมาะสำหรับ Primary Key
    id BIGSERIAL PRIMARY KEY,

    -- SKU (Stock Keeping Unit) สำหรับอ้างอิงสินค้า, ตั้งค่าให้ไม่ซ้ำกัน (UNIQUE)
    sku VARCHAR(100) UNIQUE NOT NULL,

    name VARCHAR(255) NOT NULL,

    description TEXT,

    -- DECIMAL เหมาะสำหรับเก็บข้อมูลทางการเงิน เพราะจะไม่มีปัญหาเรื่องทศนิยมคลาดเคลื่อนเหมือน FLOAT
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),

    -- จำนวนสินค้าในคลัง, ค่าเริ่มต้นคือ 0 และต้องไม่น้อยกว่า 0
    stock_quantity INT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),

    -- Foreign key ไปยังตาราง categories (ถ้ามี)
    -- เราจะสร้าง FK constraint หลังจากมีตาราง categories แล้ว
    category_id BIGINT,

    image_url VARCHAR(255),

    -- TIMESTAMPTZ จะเก็บเวลาพร้อมโซนเวลา (Timezone) ซึ่งเป็น Best Practice
    -- DEFAULT NOW() จะใส่เวลาปัจจุบันให้ทันทีที่สร้างแถวใหม่
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- เพิ่ม Index เพื่อช่วยให้การค้นหาด้วยชื่อหรือ sku เร็วขึ้น
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_sku ON products(sku);

-- (Optional but Recommended) สร้าง Trigger เพื่ออัปเดต updated_at อัตโนมัติ
-- 1. สร้าง Function
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. ผูก Function นี้เข้ากับตาราง products
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
