const { test } = require('node:test');
const assert = require('node:assert');

test('hello world', () => {
    assert.strictEqual('hello world', 'hello world');
});
