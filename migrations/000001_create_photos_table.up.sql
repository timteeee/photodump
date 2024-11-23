create table if not exists photos (
	id uuid primary key,
	filename varchar not null,
	size int not null,
	uploaded_at timestamptz not null
);

create index if not exists photo_upload_time on photos (uploaded_at);
