CREATE TABLE public.activity_log (
	trace_id varchar(50) NOT NULL,
	endpoint varchar(50) NULL,
	"path" varchar(50) NULL,
	description text NULL,
	created_at varchar(50) NULL,
	request_payload text NULL,
    response_payload text NULL
);
CREATE INDEX "ACTIVITY_LOG_INSERT_DATE_IDX" ON public.activity_log USING btree (path, created_at);