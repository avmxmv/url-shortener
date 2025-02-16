-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_links_original_url_unique ON links(original_url);
CREATE UNIQUE INDEX idx_links_short_url_unique ON links(short_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_links_original_url_unique;
DROP INDEX idx_links_short_url_unique;
-- +goose StatementEnd
