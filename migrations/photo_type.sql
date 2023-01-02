DO $$ BEGIN
    CREATE TYPE photo_type AS ENUM (
        'pribadi',
        'kartu-pelajar',
        'bukti-mahasiswa-aktif'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$