CREATE TABLE IF NOT EXISTS url_maps (
    id int comment"编号 雪花算法生成唯一编号" PRIMARY KEY,
    long_url TEXT comment"初始url" NOT NULL,
    short_url TEXT comment"生成的短url" NOT NULL ,
    is_custom bool comment"是否是自定义的短url" DEFAULT FALSE,
    expired_at timestamp comment"过期时间 时间戳存储" NOT NULL ,
    created_at timestamp comment"创建时间 时间戳存储" NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_short_code ON urls(short_url);
CREATE INDEX idx_expired_at ON urls(expired_at);