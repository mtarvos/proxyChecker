-- +goose Up
-- +goose StatementBegin
ALTER TABLE proxy ADD COLUMN out_ip TEXT DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_out_ip ON proxy(out_ip);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_out_ip;

CREATE TABLE IF NOT EXISTS proxy_temp(
    id INTEGER PRIMARY KEY,
    proxy TEXT NOT NULL,
    port INTEGER NOT NULL,
    country TEXT,
    city TEXT,
    ISP TEXT,
    timezone INTEGER,
    alive INTEGER CHECK(alive IN (0, 1, 2))
);

insert into proxy_temp (id, proxy, port, timezone, country, city, ISP, alive)
SELECT id, proxy, port, timezone, country, city, ISP, alive
FROM proxy;

DROP TABLE proxy;

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_ip_port ON proxy_temp(proxy, port);
CREATE INDEX IF NOT EXISTS idx_proxy ON proxy_temp(proxy);
CREATE INDEX IF NOT EXISTS idx_country ON proxy_temp(country);
CREATE INDEX IF NOT EXISTS idx_city ON proxy_temp(city);
CREATE INDEX IF NOT EXISTS idx_ISP ON proxy_temp(ISP);
CREATE INDEX IF NOT EXISTS idx_alive ON proxy_temp(alive);

ALTER TABLE proxy_temp RENAME TO proxy;
-- +goose StatementEnd
