DO $$
BEGIN
	CREATE TABLE IF NOT EXISTS activity_log (
		trace_id varchar(100) NOT NULL,
		endpoint varchar(100) NULL,
		"path" varchar(100) NULL,
		description text NULL,
		created_at varchar(100) NULL,
		request_payload text NULL,
    	response_payload text NULL
	);
	CREATE INDEX "ACTIVITY_LOG_INSERT_DATE_IDX" ON public.activity_log USING btree (path, created_at);
END $$;