export const TestRunner = {
    passed: 0,
    failed: 0,
    results: [],

    assert(condition, message) {
        if (condition) {
            this.passed++;
            this.results.push({ status: 'PASS', message });
        } else {
            this.failed++;
            this.results.push({ status: 'FAIL', message });
        }
    },

    assertEquals(actual, expected, message) {
        this.assert(actual === expected, message + ` (expected: ${expected}, got: ${actual})`);
    },

    printResults() {
        console.log('');
        this.results.forEach(r => {
            const symbol = r.status === 'PASS' ? '\x1b[32m✓\x1b[0m' : '\x1b[31m✗\x1b[0m';
            console.log(`${symbol} ${r.message}`);
        });
        console.log('');
        console.log(`Results: ${this.passed} passed, ${this.failed} failed`);
        return this.failed === 0;
    }
};
