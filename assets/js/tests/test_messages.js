// Run with: node assets/js/tests/test_messages.js
import {TestRunner} from "./test_runner.js";
import * as messages from "../entities/messages.js";

// Every message must survive the wire round trip: instance → JSON → fromObject → same instance.
// Catches a field that fromObject forgets to copy or copies under a different name.
for (const [name, MessageClass] of Object.entries(messages)) {
    const original = new MessageClass();
    for (const field of Object.keys(original)) {
        original[field] = field + '-value';
    }

    const decoded = new MessageClass().fromObject(JSON.parse(JSON.stringify(original)));

    TestRunner.assertEquals(
        JSON.stringify(decoded),
        JSON.stringify(original),
        `${name} survives JSON round trip through fromObject`,
    );
}

if (!TestRunner.printResults()) {
    process.exit(1);
}
