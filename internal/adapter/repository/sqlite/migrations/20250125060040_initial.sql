-- +goose Up
CREATE TABLE IF NOT EXISTS proxy(
                                    id INTEGER PRIMARY KEY,
                                    proxy TEXT NOT NULL,
                                    port INTEGER NOT NULL,
                                    country TEXT,
                                    city TEXT,
                                    ISP TEXT,
                                    timezone INTEGER DEFAULT -1,
                                    alive INTEGER CHECK(alive IN (0, 1, 2))
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_ip_port ON proxy(proxy, port);
CREATE INDEX IF NOT EXISTS idx_proxy ON proxy(proxy);
CREATE INDEX IF NOT EXISTS idx_country ON proxy(country);
CREATE INDEX IF NOT EXISTS idx_city ON proxy(city);
CREATE INDEX IF NOT EXISTS idx_ISP ON proxy(ISP);
CREATE INDEX IF NOT EXISTS idx_alive ON proxy(alive);

-- +goose Down
DROP TABLE proxy;
