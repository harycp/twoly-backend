ALTER TABLE memories
ADD COLUMN converted_from_date_plan_id UUID REFERENCES date_plans(id) ON DELETE SET NULL;
