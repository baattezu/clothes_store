ALTER TABLE clothes_info
    ADD CONSTRAINT check_updated_at_after_created_at CHECK (updated_at >= created_at),
    ADD CONSTRAINT check_module_duration_range CHECK (cloth_cost > 0 AND cloth_cost <= 200000),
    ADD CONSTRAINT check_cloth_size CHECK (cloth_size IN ('s', 'l', 'xl', 'xxl'));
