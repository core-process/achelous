const TEST_DIR = '/tmp/testservice';

const fs = require('fs');
const express = require('express');
const bodyParser = require('body-parser');

const app = express();

app.use(bodyParser.json({
    limit: '50mb',
}));

app.post('/upload', (req, res) => {
    fs.writeFileSync(
        TEST_DIR + '/' + Date.now() + '_upload.json',
        JSON.stringify(req.body)
    );
    res.send('OK');
});

app.post('/report', (req, res) => {
    fs.writeFileSync(
        TEST_DIR + '/' + Date.now() + '_report.json',
        JSON.stringify(req.body)
    );
    res.send('OK');
});

if (!fs.existsSync(TEST_DIR)) {
    fs.mkdirSync(TEST_DIR);
}

app.listen(3000, () => console.log('Testservice listening on port 3000!'));
