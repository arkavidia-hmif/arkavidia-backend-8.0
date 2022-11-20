CREATE TYPE participant_career_interest AS ENUM (
    'software-engineering',
    'product-management',
    'ui-designer',
    'ux-designer',
    'ux-researcher',
    'it-consultant',
    'game-developer',
    'cyber-security',
    'business-analyst',
    'business-intelligence',
    'data-scientist',
    'data-analyst'
);

CREATE TYPE team_category AS ENUM (
    'competitive-programming',
    'datavidia',
    'uxvidia',
    'arkalogica'
);

CREATE TYPE membership_role AS ENUM (
    'leader',
    'member-1',
    'member-2'
);

CREATE TYPE photo_status AS ENUM (
    'waiting-for-verification',
    'verified',
    'declined'
);

CREATE TYPE photo_type AS ENUM (
    'pribadi',
    'kartu-pelajar',
    'bukti-mahasiswa-aktif'
);

CREATE TYPE submission_stage AS ENUM (
    'first-stage',
    'second-stage',
    'final-stage'
);