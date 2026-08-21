const fs = require('fs');
fs.mkdirSync('dist', { recursive: true });
fs.writeFileSync('dist/index.html', '<!doctype html><title>Park Visitor Sync</title><main>Park visitor synchronization console</main>');
